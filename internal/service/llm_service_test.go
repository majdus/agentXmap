package service

import (
	"agentXmap/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLLMRepository is a mock implementation of domain.LLMRepository
type MockLLMRepository struct {
	mock.Mock
}

func (m *MockLLMRepository) ListProviders(ctx context.Context) ([]domain.LLMProvider, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LLMProvider), args.Error(1)
}

func (m *MockLLMRepository) GetModel(ctx context.Context, id uuid.UUID) (*domain.LLMModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LLMModel), args.Error(1)
}

func (m *MockLLMRepository) ListModels(ctx context.Context, providerID uuid.UUID) ([]domain.LLMModel, error) {
	args := m.Called(ctx, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LLMModel), args.Error(1)
}

func (m *MockLLMRepository) ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockLLMRepository) GetCertifications(ctx context.Context, modelID uuid.UUID) ([]domain.Certification, error) {
	args := m.Called(ctx, modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Certification), args.Error(1)
}

func TestLLMService_ListProviders(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		expectedProviders := []domain.LLMProvider{
			{Name: "OpenAI"},
			{Name: "Anthropic"},
		}
		mockRepo.On("ListProviders", ctx).Return(expectedProviders, nil)

		providers, err := service.ListProviders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedProviders), len(providers))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		mockRepo.On("ListProviders", ctx).Return(nil, errors.New("db error"))

		_, err := service.ListProviders(ctx)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestLLMService_ListModels(t *testing.T) {
	ctx := context.Background()
	providerID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		expectedModels := []domain.LLMModel{
			{ApiModelName: "gpt-4"},
			{ApiModelName: "gpt-3.5-turbo"},
		}
		mockRepo.On("ListModels", ctx, providerID).Return(expectedModels, nil)

		models, err := service.ListModels(ctx, providerID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedModels), len(models))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		mockRepo.On("ListModels", ctx, providerID).Return(nil, errors.New("db error"))

		_, err := service.ListModels(ctx, providerID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestLLMService_GetModel(t *testing.T) {
	ctx := context.Background()
	modelID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		expectedModel := &domain.LLMModel{ID: modelID, ApiModelName: "gpt-4"}
		mockRepo.On("GetModel", ctx, modelID).Return(expectedModel, nil)

		model, err := service.GetModel(ctx, modelID)
		assert.NoError(t, err)
		assert.Equal(t, "gpt-4", model.ApiModelName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		mockRepo.On("GetModel", ctx, modelID).Return(nil, errors.New("model not found"))

		_, err := service.GetModel(ctx, modelID)
		assert.Error(t, err)
		assert.Equal(t, "model not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		mockRepo.On("GetModel", ctx, modelID).Return(nil, errors.New("db error"))

		_, err := service.GetModel(ctx, modelID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestLLMService_ListAgentsUsingModel(t *testing.T) {
	ctx := context.Background()
	modelID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		expectedAgents := []domain.Agent{
			{Name: "Agent 007"},
		}
		mockRepo.On("ListAgentsUsingModel", ctx, modelID).Return(expectedAgents, nil)

		agents, err := service.ListAgentsUsingModel(ctx, modelID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedAgents), len(agents))
		mockRepo.AssertExpectations(t)
	})
}

func TestLLMService_ListModelCertifications(t *testing.T) {
	ctx := context.Background()
	modelID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockLLMRepository)
		service := NewLLMService(mockRepo)

		expectedCerts := []domain.Certification{
			{Name: "EU AI Act Compliant"},
		}
		mockRepo.On("GetCertifications", ctx, modelID).Return(expectedCerts, nil)

		certs, err := service.ListModelCertifications(ctx, modelID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedCerts), len(certs))
		mockRepo.AssertExpectations(t)
	})
}
