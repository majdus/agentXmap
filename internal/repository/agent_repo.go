package repository

import (
	"context"
	"errors"

	"agentXmap/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type agentRepository struct {
	db *gorm.DB
}

// NewAgentRepository creates a new postgres repository for Agents.
func NewAgentRepository(db *gorm.DB) domain.AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(ctx context.Context, agent *domain.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *agentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	var agent domain.Agent
	// Preload minimal relations
	if err := r.db.WithContext(ctx).Preload("Versions").Preload("Organization").First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.Agent, error) {
	var agents []domain.Agent
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *agentRepository) GetAssignedLLMs(ctx context.Context, agentID uuid.UUID) ([]domain.AgentLLM, error) {
	var agentLLMs []domain.AgentLLM
	// Load AgentLLM with associated LLMModel details
	if err := r.db.WithContext(ctx).
		Preload("LLMModel").
		Where("agent_id = ?", agentID).
		Find(&agentLLMs).Error; err != nil {
		return nil, err
	}
	return agentLLMs, nil
}

func (r *agentRepository) Update(ctx context.Context, agent *domain.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *agentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Agent{}, "id = ?", id).Error
}

func (r *agentRepository) CreateVersion(ctx context.Context, version *domain.AgentVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *agentRepository) GetResources(ctx context.Context, agentID uuid.UUID) ([]domain.Resource, error) {
	var resources []domain.Resource
	// Join AgentResourceAccess to find resources linked to this agent
	err := r.db.WithContext(ctx).
		Joins("JOIN agent_resource_accesses ON agent_resource_accesses.resource_id = resources.id").
		Where("agent_resource_accesses.agent_id = ?", agentID).
		Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *agentRepository) GetAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	// Join AgentAssignment to find users linked to this agent
	err := r.db.WithContext(ctx).
		Joins("JOIN agent_assignments ON agent_assignments.user_id = users.id").
		Where("agent_assignments.agent_id = ?", agentID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *agentRepository) GetAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	var apps []domain.Application
	// Join ApplicationAgentAccess (and then Application if needed, but ApplicationAgentAccess belongs to Application? No, ApplicationAgentAccess links Application and Agent)
	// struct: ApplicationAgentAccess has ApplicationID and AgentID.
	// We want to list Applications.
	// JOIN application_agent_accesses ON application_agent_accesses.application_id = applications.id
	// WHERE application_agent_accesses.agent_id = ?

	err := r.db.WithContext(ctx).
		Joins("JOIN application_agent_accesses ON application_agent_accesses.application_id = applications.id").
		Where("application_agent_accesses.agent_id = ?", agentID).
		Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}
