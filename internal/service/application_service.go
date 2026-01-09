package service

import (
	"agentXmap/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ApplicationService interface {
	CreateApplication(ctx context.Context, ownerID uuid.UUID, name, description string) (*domain.Application, error)
	GetApplication(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	CreateAPIKey(ctx context.Context, appID uuid.UUID, name string) (string, *domain.ApplicationKey, error)
}

type DefaultApplicationService struct {
	appRepo domain.ApplicationRepository
}

func NewApplicationService(appRepo domain.ApplicationRepository) *DefaultApplicationService {
	return &DefaultApplicationService{appRepo: appRepo}
}

func (s *DefaultApplicationService) CreateApplication(ctx context.Context, ownerID uuid.UUID, name, description string) (*domain.Application, error) {
	if name == "" {
		return nil, errors.New("application name is required")
	}

	app := &domain.Application{
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		IsActive:    true,
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *DefaultApplicationService) GetApplication(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("application not found")
	}
	return app, nil
}

func (s *DefaultApplicationService) CreateAPIKey(ctx context.Context, appID uuid.UUID, name string) (string, *domain.ApplicationKey, error) {
	// Generate a secure random key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	rawKey := "sk-live-" + hex.EncodeToString(keyBytes)

	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("failed to hash key: %w", err)
	}

	// Create the key record
	key := &domain.ApplicationKey{
		ApplicationID: appID,
		KeyHash:       string(hashedKey),
		KeyPrefix:     rawKey[:16] + "...", // Store prefix for display
		Name:          name,
	}
	// Correct prefix storage (sk-live-...)
	// Actually typical pattern is to store "sk-live-" + first few chars.
	// rawKey is "sk-live-" (8 chars) + 64 hex chars.
	// Let's store "sk-live-" + first 4 chars of hex as prefix.
	// Total prefix length: 8+4 = 12 chars.
	key.KeyPrefix = rawKey[:12]

	if err := s.appRepo.CreateKey(ctx, key); err != nil {
		return "", nil, err
	}

	return rawKey, key, nil
}

// ValidateKey (future helper) checks if a provided raw API key matches the stored hash
func (s *DefaultApplicationService) ValidateKey(rawKey, storedHash string) bool {
	if !strings.HasPrefix(rawKey, "sk-live-") {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(rawKey))
	return err == nil
}
