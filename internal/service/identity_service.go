package service

import (
	"agentXmap/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// IdentityService defines the interface for user identity management.
type IdentityService interface {
	SignUp(ctx context.Context, orgName, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	InviteUsers(ctx context.Context, invitorID uuid.UUID, emails []string, role domain.UserRole) ([]*domain.Invitation, error)
	AcceptInvitation(ctx context.Context, token, password, firstName, lastName string) (*domain.User, error)
}

type DefaultIdentityService struct {
	userRepo       domain.UserRepository
	orgRepo        domain.OrganizationRepository
	invitationRepo domain.InvitationRepository
	// In a real app we would have a PasswordHasher and EmailService interface here
}

// NewIdentityService creates a new instance of DefaultIdentityService.
func NewIdentityService(
	userRepo domain.UserRepository,
	orgRepo domain.OrganizationRepository,
	invitationRepo domain.InvitationRepository,
) *DefaultIdentityService {
	return &DefaultIdentityService{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		invitationRepo: invitationRepo,
	}
}

// Helper for password hashing
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Helper for random token generation
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *DefaultIdentityService) SignUp(ctx context.Context, orgName, email, password string) (*domain.User, error) {
	// TODO: Add proper validation (email format, password strength)
	// TODO: Check if email already exists? (Repo will throw error, but better to check)
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Slug generation (simplistic for now)
	slug := Slugify(orgName)

	// We should run this in a transaction.
	// Gorm doesn't easily expose "RunTransaction" on Repos unless we share the DB instance or pass it around.
	// For now we do it sequentially. If Org creation succeeds but User fails, we have an orphan Org.
	// Ideally Repos should accept a DB/Tx interface.

	org := &domain.Organization{
		Name: orgName,
		Slug: slug,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	user := &domain.User{
		OrganizationID: org.ID,
		Email:          email,
		PasswordHash:   hashedPassword,
		Role:           domain.UserRoleAdmin, // First user is Admin
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		// Rollback Org?
		// s.orgRepo.Delete(ctx, org.ID)
		return nil, err
	}

	user.Organization = *org
	return user, nil
}

func (s *DefaultIdentityService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *DefaultIdentityService) InviteUsers(ctx context.Context, invitorID uuid.UUID, emails []string, role domain.UserRole) ([]*domain.Invitation, error) {
	invitor, err := s.userRepo.GetByID(ctx, invitorID)
	if err != nil {
		return nil, errors.New("invitor not found")
	}

	if invitor.Role != domain.UserRoleAdmin && invitor.Role != domain.UserRoleManager {
		return nil, errors.New("insufficient permissions to invite users")
	}

	var invitations []*domain.Invitation

	for _, email := range emails {
		// Skip existing users? or allow re-invite?
		// For now simple check:
		if _, err := s.userRepo.GetByEmail(ctx, email); err == nil {
			// User exists, skip or error? Let's skip and continue.
			// Or maybe return error for partial failure?
			continue
		}

		token, err := generateToken()
		if err != nil {
			return nil, err
		}

		invitation := &domain.Invitation{
			OrganizationID: invitor.OrganizationID,
			InvitorID:      invitor.ID,
			Email:          email,
			Token:          token,
			Role:           role,
			Status:         domain.InvitationStatusPending,
			ExpiresAt:      time.Now().Add(48 * time.Hour),
		}

		if err := s.invitationRepo.Create(ctx, invitation); err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (s *DefaultIdentityService) AcceptInvitation(ctx context.Context, token, password, firstName, lastName string) (*domain.User, error) {
	invitation, err := s.invitationRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid invitation token")
	}

	if invitation.Status != domain.InvitationStatusPending {
		return nil, errors.New("invitation is not pending")
	}

	if time.Now().After(invitation.ExpiresAt) {
		invitation.Status = domain.InvitationStatusExpired
		_ = s.invitationRepo.Update(ctx, invitation)
		return nil, errors.New("invitation expired")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		OrganizationID: invitation.OrganizationID,
		Email:          invitation.Email,
		PasswordHash:   hashedPassword,
		Role:           invitation.Role,
		FirstName:      firstName,
		LastName:       lastName,
	}

	// This should also be transactional
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	invitation.Status = domain.InvitationStatusAccepted
	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		// Log error, but user is created
	}

	return user, nil
}

// Simple slugify helper
// Slugify converts a string to a valid URL slug.
func Slugify(s string) string {
	// Lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with dashes
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")

	// Trim dashes from start and end
	s = strings.Trim(s, "-")

	return s
}
