package service

import (
	"agentXmap/internal/domain"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockApplicationRepository is a mock implementation of domain.ApplicationRepository
type MockApplicationRepository struct {
	mock.Mock
}

func (m *MockApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *MockApplicationRepository) CreateKey(ctx context.Context, key *domain.ApplicationKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestApplicationService_CreateApplication(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockApplicationRepository)
		service := NewApplicationService(mockRepo)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Application")).Return(nil)

		app, err := service.CreateApplication(ctx, ownerID, "Test App", "Description")
		assert.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, "Test App", app.Name)
		assert.Equal(t, ownerID, app.OwnerID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty Name", func(t *testing.T) {
		mockRepo := new(MockApplicationRepository)
		service := NewApplicationService(mockRepo)

		_, err := service.CreateApplication(ctx, ownerID, "", "Description")
		assert.Error(t, err)
		assert.Equal(t, "application name is required", err.Error())
	})
}

func TestApplicationService_GetApplication(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockApplicationRepository)
		service := NewApplicationService(mockRepo)

		expectedApp := &domain.Application{ID: appID, Name: "Test App"}
		mockRepo.On("GetByID", ctx, appID).Return(expectedApp, nil)

		app, err := service.GetApplication(ctx, appID)
		assert.NoError(t, err)
		assert.Equal(t, "Test App", app.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockApplicationRepository)
		service := NewApplicationService(mockRepo)

		mockRepo.On("GetByID", ctx, appID).Return(nil, errors.New("db error"))

		_, err := service.GetApplication(ctx, appID)
		assert.Error(t, err)
		assert.Equal(t, "application not found", err.Error())
	})
}

func TestApplicationService_CreateAPIKey(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockApplicationRepository)
		service := NewApplicationService(mockRepo)

		mockRepo.On("CreateKey", ctx, mock.AnythingOfType("*domain.ApplicationKey")).Return(nil)

		rawKey, key, err := service.CreateAPIKey(ctx, appID, "Test Key")
		assert.NoError(t, err)
		assert.NotNil(t, key)
		assert.NotEmpty(t, rawKey)
		assert.True(t, strings.HasPrefix(rawKey, "sk-live-"))
		assert.NotEmpty(t, key.KeyHash)
		assert.Equal(t, "Test Key", key.Name)
		assert.Equal(t, appID, key.ApplicationID)
		mockRepo.AssertExpectations(t)
	})
}
