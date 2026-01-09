package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	if args.Error(0) == nil {
		org.ID = uuid.New() // Simulate DB generating ID
	}
	return args.Error(0)
}

func (m *MockOrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

type MockInvitationRepository struct {
	mock.Mock
}

func (m *MockInvitationRepository) Create(ctx context.Context, invitation *domain.Invitation) error {
	args := m.Called(ctx, invitation)
	return args.Error(0)
}

func (m *MockInvitationRepository) GetByToken(ctx context.Context, token string) (*domain.Invitation, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invitation), args.Error(1)
}

func (m *MockInvitationRepository) Update(ctx context.Context, invitation *domain.Invitation) error {
	args := m.Called(ctx, invitation)
	return args.Error(0)
}

// Tests

func TestIdentityService_SignUp(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockInvRepo := new(MockInvitationRepository)
	service := NewIdentityService(mockUserRepo, mockOrgRepo, mockInvRepo)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", ctx, "admin@test.com").Return(nil, errors.New("not found")).Once()
		mockOrgRepo.On("Create", ctx, mock.AnythingOfType("*domain.Organization")).Return(nil).Once()
		mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		user, err := service.SignUp(ctx, "Test Org", "admin@test.com", "password123")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "admin@test.com", user.Email)
		assert.Equal(t, domain.UserRoleAdmin, user.Role)
		assert.Equal(t, "Test Org", user.Organization.Name)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		existingUser := &domain.User{Email: "admin@test.com"}
		mockUserRepo.On("GetByEmail", ctx, "admin@test.com").Return(existingUser, nil).Once()

		user, err := service.SignUp(ctx, "Test Org", "admin@test.com", "password123")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user already exists", err.Error())
	})
}

func TestIdentityService_InviteUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockInvRepo := new(MockInvitationRepository)
	service := NewIdentityService(mockUserRepo, mockOrgRepo, mockInvRepo)

	ctx := context.Background()
	invitorID := uuid.New()
	orgID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		invitor := &domain.User{
			ID:             invitorID,
			OrganizationID: orgID,
			Role:           domain.UserRoleAdmin,
		}

		mockUserRepo.On("GetByID", ctx, invitorID).Return(invitor, nil).Once()
		// Mock GetByEmail for each email (assuming they don't exist)
		mockUserRepo.On("GetByEmail", ctx, "invitee@test.com").Return(nil, errors.New("not found")).Once()
		mockInvRepo.On("Create", ctx, mock.AnythingOfType("*domain.Invitation")).Return(nil).Once()

		invitations, err := service.InviteUsers(ctx, invitorID, []string{"invitee@test.com"}, domain.UserRoleUser)
		assert.NoError(t, err)
		assert.Len(t, invitations, 1)
		assert.Equal(t, "invitee@test.com", invitations[0].Email)
		assert.Equal(t, domain.InvitationStatusPending, invitations[0].Status)
	})

	t.Run("InsufficientPermissions", func(t *testing.T) {
		invitor := &domain.User{
			ID:             invitorID,
			OrganizationID: orgID,
			Role:           domain.UserRoleUser, // Normal user cannot invite
		}

		mockUserRepo.On("GetByID", ctx, invitorID).Return(invitor, nil).Once()

		invitations, err := service.InviteUsers(ctx, invitorID, []string{"invitee@test.com"}, domain.UserRoleUser)
		assert.Error(t, err)
		assert.Nil(t, invitations)
		assert.Equal(t, "insufficient permissions to invite users", err.Error())
	})

	t.Run("InvitorNotFound", func(t *testing.T) {
		mockUserRepo.On("GetByID", ctx, invitorID).Return(nil, errors.New("not found")).Once()

		invitations, err := service.InviteUsers(ctx, invitorID, []string{"invitee@test.com"}, domain.UserRoleUser)
		assert.Error(t, err)
		assert.Nil(t, invitations)
	})
}

func TestIdentityService_AcceptInvitation(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOrgRepo := new(MockOrganizationRepository)
	mockInvRepo := new(MockInvitationRepository)
	service := NewIdentityService(mockUserRepo, mockOrgRepo, mockInvRepo)

	ctx := context.Background()
	token := "valid-token"

	t.Run("Success", func(t *testing.T) {
		invitation := &domain.Invitation{
			Token:     token,
			Status:    domain.InvitationStatusPending,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Email:     "newuser@test.com",
			Role:      domain.UserRoleUser,
		}

		mockInvRepo.On("GetByToken", ctx, token).Return(invitation, nil).Once()
		mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		mockInvRepo.On("Update", ctx, mock.MatchedBy(func(inv *domain.Invitation) bool {
			return inv.Status == domain.InvitationStatusAccepted
		})).Return(nil).Once()

		user, err := service.AcceptInvitation(ctx, token, "password123", "John", "Doe")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "newuser@test.com", user.Email)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		invitation := &domain.Invitation{
			Token:     token,
			Status:    domain.InvitationStatusPending,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		}

		mockInvRepo.On("GetByToken", ctx, token).Return(invitation, nil).Once()
		mockInvRepo.On("Update", ctx, mock.MatchedBy(func(inv *domain.Invitation) bool {
			return inv.Status == domain.InvitationStatusExpired
		})).Return(nil).Once()

		user, err := service.AcceptInvitation(ctx, token, "password123", "John", "Doe")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invitation expired", err.Error())
	})

	t.Run("InvalidToken", func(t *testing.T) {
		mockInvRepo.On("GetByToken", ctx, "invalid").Return(nil, errors.New("not found")).Once()

		user, err := service.AcceptInvitation(ctx, "invalid", "password123", "John", "Doe")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
