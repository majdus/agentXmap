package repository

import (
	"context"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *resourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) Create(ctx context.Context, res *domain.Resource) error {
	return r.db.WithContext(ctx).Create(res).Error
}

func (r *resourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	var res domain.Resource
	if err := r.db.WithContext(ctx).Preload("Type").Preload("Secret").First(&res, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
