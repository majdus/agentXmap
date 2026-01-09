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
