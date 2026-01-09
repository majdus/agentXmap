package repository

import (
	"context"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new postgres repository for Organizations.
func NewOrganizationRepository(db *gorm.DB) domain.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization
	if err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	var org domain.Organization
	if err := r.db.WithContext(ctx).First(&org, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &org, nil
}
