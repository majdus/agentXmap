package repository

import (
	"context"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Add Interface to domain/interfaces.go first?
// No, I should assume I need to implement what makes sense or update domain/interfaces.go in parallel.
// But as per plan I should stick to adding implementation unless I forgot the interface definition.
// Checking previous steps, I defined `UserRepository`, `AgentRepository`, `LLMRepository`, `AuditRepository`.
// I missed `ApplicationRepository` and `ResourceRepository` in the initial `interfaces.go`.
// I will create the implementations and then update the interfaces.go file.

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *applicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var app domain.Application
	if err := r.db.WithContext(ctx).Preload("Keys").Preload("AgentAccess").First(&app, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) CreateKey(ctx context.Context, key *domain.ApplicationKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}
