package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResourceRepository is a mock implementation of domain.ResourceRepository
type MockResourceRepository struct {
	mock.Mock
}

func (m *MockResourceRepository) Create(ctx context.Context, res *domain.Resource) error {
	args := m.Called(ctx, res)
	return args.Error(0)
}

func (m *MockResourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Resource), args.Error(1)
}

func (m *MockResourceRepository) ListAgentsWithAccess(ctx context.Context, resourceID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, resourceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func TestResourceService_CreateResource(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	config := json.RawMessage(`{"host":"localhost"}`)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resource")).Return(nil)

		res, err := service.CreateResource(ctx, orgID, "postgres-db", "Test DB", config)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Test DB", res.Name)
		assert.Equal(t, "postgres-db", res.TypeID)
		assert.Equal(t, orgID, res.OrganizationID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Empty Name", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		_, err := service.CreateResource(ctx, orgID, "postgres-db", "", config)
		assert.Error(t, err)
		assert.Equal(t, "resource name is required", err.Error())
	})

	t.Run("Validation Error - Empty Type", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		_, err := service.CreateResource(ctx, orgID, "", "Test DB", config)
		assert.Error(t, err)
		assert.Equal(t, "resource type is required", err.Error())
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resource")).Return(errors.New("db error"))

		_, err := service.CreateResource(ctx, orgID, "postgres-db", "Test DB", config)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestResourceService_GetResource(t *testing.T) {
	ctx := context.Background()
	resID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		expectedRes := &domain.Resource{ID: resID, Name: "Test Resource"}
		mockRepo.On("GetByID", ctx, resID).Return(expectedRes, nil)

		res, err := service.GetResource(ctx, resID)
		assert.NoError(t, err)
		assert.Equal(t, "Test Resource", res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockResourceRepository)
		service := NewResourceService(mockRepo)

		mockRepo.On("GetByID", ctx, resID).Return(nil, errors.New("db error"))

		_, err := service.GetResource(ctx, resID)
		assert.Error(t, err)
		assert.Equal(t, "resource not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
