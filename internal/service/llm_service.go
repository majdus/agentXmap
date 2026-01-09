package service

import (
	"agentXmap/internal/domain"
	"context"

	"github.com/google/uuid"
)

type LLMService interface {
	ListProviders(ctx context.Context) ([]domain.LLMProvider, error)
	ListModels(ctx context.Context, providerID uuid.UUID) ([]domain.LLMModel, error)
	GetModel(ctx context.Context, id uuid.UUID) (*domain.LLMModel, error)
	ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]domain.Agent, error)
}

type DefaultLLMService struct {
	llmRepo domain.LLMRepository
}

func NewLLMService(llmRepo domain.LLMRepository) *DefaultLLMService {
	return &DefaultLLMService{llmRepo: llmRepo}
}

func (s *DefaultLLMService) ListProviders(ctx context.Context) ([]domain.LLMProvider, error) {
	return s.llmRepo.ListProviders(ctx)
}

func (s *DefaultLLMService) ListModels(ctx context.Context, providerID uuid.UUID) ([]domain.LLMModel, error) {
	return s.llmRepo.ListModels(ctx, providerID)
}

func (s *DefaultLLMService) GetModel(ctx context.Context, id uuid.UUID) (*domain.LLMModel, error) {
	return s.llmRepo.GetModel(ctx, id)
}

func (s *DefaultLLMService) ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]domain.Agent, error) {
	return s.llmRepo.ListAgentsUsingModel(ctx, modelID)
}
