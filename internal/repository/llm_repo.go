package repository

import (
	"context"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type llmRepository struct {
	db *gorm.DB
}

func NewLLMRepository(db *gorm.DB) domain.LLMRepository {
	return &llmRepository{db: db}
}

func (r *llmRepository) ListProviders(ctx context.Context) ([]domain.LLMProvider, error) {
	var providers []domain.LLMProvider
	// Preload Models for each provider
	if err := r.db.WithContext(ctx).Preload("Models").Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *llmRepository) GetModel(ctx context.Context, id uuid.UUID) (*domain.LLMModel, error) {
	var model domain.LLMModel
	if err := r.db.WithContext(ctx).Preload("Provider").First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *llmRepository) ListModels(ctx context.Context, providerID uuid.UUID) ([]domain.LLMModel, error) {
	var models []domain.LLMModel
	query := r.db.WithContext(ctx)
	if providerID != uuid.Nil {
		query = query.Where("provider_id = ?", providerID)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *llmRepository) ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]domain.Agent, error) {
	var agents []domain.Agent
	if err := r.db.WithContext(ctx).
		Joins("JOIN agent_llms ON agent_llms.agent_id = agents.id").
		Where("agent_llms.llm_model_id = ?", modelID).
		Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *llmRepository) GetCertifications(ctx context.Context, modelID uuid.UUID) ([]domain.Certification, error) {
	var certifications []domain.Certification
	err := r.db.WithContext(ctx).
		Joins("JOIN llm_model_certifications ON llm_model_certifications.certification_id = certifications.id").
		Where("llm_model_certifications.llm_model_id = ?", modelID).
		Find(&certifications).Error
	if err != nil {
		return nil, err
	}
	return certifications, nil
}
