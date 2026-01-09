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

func (r *applicationRepository) GetAssignedAgents(ctx context.Context, appID uuid.UUID) ([]domain.Agent, error) {
	var agents []domain.Agent
	// Join ApplicationAgentAccess to find agents linked to this application
	// Remember ApplicationAgentAccess has ApplicationID and AgentID
	err := r.db.WithContext(ctx).
		Joins("JOIN application_agent_accesses ON application_agent_accesses.agent_id = agents.id").
		Where("application_agent_accesses.application_id = ?", appID).
		Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *applicationRepository) GetCertifications(ctx context.Context, appID uuid.UUID) ([]domain.Certification, error) {
	var certifications []domain.Certification
	err := r.db.WithContext(ctx).
		Joins("JOIN application_certifications ON application_certifications.certification_id = certifications.id").
		Where("application_certifications.application_id = ?", appID).
		Find(&certifications).Error
	if err != nil {
		return nil, err
	}
	return certifications, nil
}
