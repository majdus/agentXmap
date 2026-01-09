package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// AgentService defines the interface for agent management.
type AgentService interface {
	CreateAgent(ctx context.Context, orgID, userID uuid.UUID, name string, config json.RawMessage) (*domain.Agent, error)
	GetAgent(ctx context.Context, id uuid.UUID) (*domain.Agent, error)
	ListAgents(ctx context.Context, orgID uuid.UUID) ([]domain.Agent, error)
	ListAgentsByStatus(ctx context.Context, orgID uuid.UUID, status domain.AgentStatus) ([]domain.Agent, error)
	GetActiveMonthlyCost(ctx context.Context, orgID uuid.UUID) (float64, error)
	ListAgentResources(ctx context.Context, agentID uuid.UUID) ([]domain.Resource, error)
	ListAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]domain.User, error)
	ListAssignedAgents(ctx context.Context, userID uuid.UUID) ([]domain.Agent, error)
	GetAgentLLMs(ctx context.Context, agentID uuid.UUID) ([]domain.AgentLLM, error)
	ListAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error)
	ListAgentCertifications(ctx context.Context, agentID uuid.UUID) ([]domain.Certification, error)
	UpdateAgent(ctx context.Context, id, userID uuid.UUID, name string, config json.RawMessage, status domain.AgentStatus) (*domain.Agent, error)
	DeleteAgent(ctx context.Context, id uuid.UUID) error
}

type DefaultAgentService struct {
	agentRepo domain.AgentRepository
}

// NewAgentService creates a new instance of DefaultAgentService.
func NewAgentService(agentRepo domain.AgentRepository) *DefaultAgentService {
	return &DefaultAgentService{
		agentRepo: agentRepo,
	}
}

func (s *DefaultAgentService) CreateAgent(ctx context.Context, orgID, userID uuid.UUID, name string, config json.RawMessage) (*domain.Agent, error) {
	if name == "" {
		return nil, errors.New("agent name is required")
	}

	agent := &domain.Agent{
		OrganizationID: orgID,
		Name:           name,
		Status:         domain.AgentStatusActive,
		Configuration:  config,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := s.agentRepo.Create(ctx, agent); err != nil {
		return nil, err
	}

	// Create initial version
	version := &domain.AgentVersion{
		AgentID:               agent.ID,
		VersionNumber:         1,
		ConfigurationSnapshot: config,
		ReasonForChange:       "Initial creation",
		CreatedBy:             &userID,
	}

	if err := s.agentRepo.CreateVersion(ctx, version); err != nil {
		// Log error but generally agent is created.
		// Ideally transactional.
	}

	return agent, nil
}

func (s *DefaultAgentService) GetAgent(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	agent, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

func (s *DefaultAgentService) ListAgents(ctx context.Context, orgID uuid.UUID) ([]domain.Agent, error) {
	return s.agentRepo.ListByOrg(ctx, orgID)
}

func (s *DefaultAgentService) ListAgentsByStatus(ctx context.Context, orgID uuid.UUID, status domain.AgentStatus) ([]domain.Agent, error) {
	return s.agentRepo.ListByStatus(ctx, orgID, status)
}

func (s *DefaultAgentService) GetActiveMonthlyCost(ctx context.Context, orgID uuid.UUID) (float64, error) {
	agents, err := s.agentRepo.ListByStatus(ctx, orgID, domain.AgentStatusActive)
	if err != nil {
		return 0, err
	}

	var totalCost float64
	for _, agent := range agents {
		// Logic:
		// Monthly -> CostAmount
		// Yearly -> CostAmount / 12
		// OneTime -> 0 (Not recurring monthly)
		// Custom -> CostAmount (Assume custom is entered as normalized, or just take it)
		switch agent.BillingCycle {
		case domain.BillingCycleMonthly:
			totalCost += agent.CostAmount
		case domain.BillingCycleYearly:
			totalCost += agent.CostAmount / 12.0
		case domain.BillingCycleOneTime:
			// One time costs are not monthly recurring
			continue
		case domain.BillingCycleCustom:
			totalCost += agent.CostAmount
		default:
			// Default to monthly if not specified/unknown, or safer: CostAmount
			totalCost += agent.CostAmount
		}
	}
	return totalCost, nil
}

func (s *DefaultAgentService) ListAgentResources(ctx context.Context, agentID uuid.UUID) ([]domain.Resource, error) {
	// First check if agent exists?
	// The repo query will just return empty list if agent doesn't exist or has no resources.
	// But good practice maybe to check agent existence for 404?
	// For now, let's just return what the repo gives.
	return s.agentRepo.GetResources(ctx, agentID)
}

func (s *DefaultAgentService) ListAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]domain.User, error) {
	return s.agentRepo.GetAssignedUsers(ctx, agentID)
}

func (s *DefaultAgentService) ListAssignedAgents(ctx context.Context, userID uuid.UUID) ([]domain.Agent, error) {
	return s.agentRepo.GetAssignedAgents(ctx, userID)
}

func (s *DefaultAgentService) GetAgentLLMs(ctx context.Context, agentID uuid.UUID) ([]domain.AgentLLM, error) {
	return s.agentRepo.GetAssignedLLMs(ctx, agentID)
}

func (s *DefaultAgentService) ListAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	return s.agentRepo.GetAssignedApplications(ctx, agentID)
}

func (s *DefaultAgentService) ListAgentCertifications(ctx context.Context, agentID uuid.UUID) ([]domain.Certification, error) {
	return s.agentRepo.GetCertifications(ctx, agentID)
}

func (s *DefaultAgentService) UpdateAgent(ctx context.Context, id, userID uuid.UUID, name string, config json.RawMessage, status domain.AgentStatus) (*domain.Agent, error) {
	agent, err := s.agentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("agent not found")
	}

	// Check if config changed
	configChanged := false
	if string(agent.Configuration) != string(config) {
		configChanged = true
	}

	agent.Name = name
	agent.Status = status
	agent.Configuration = config
	agent.UpdatedBy = &userID
	agent.UpdatedAt = time.Now()

	if err := s.agentRepo.Update(ctx, agent); err != nil {
		return nil, err
	}

	if configChanged {
		// Calculate new version number
		newVersionNum := 1
		if len(agent.Versions) > 0 {
			// Assuming versions are loaded or ordered, find max.
			// Ideally repo should provide GetLatestVersion or we just count + 1 if we loaded all (which might be heavy).
			// Start simple: existing versions list might be partial if not preload all.
			// Ideally we query DB for max version.
			// For this iteration, let's assume we implement a logic to get items.
			// Or we just increment based on count if preloaded?
			// Let's rely on versions being preloaded in GetByID (as per repo implementation).
			for _, v := range agent.Versions {
				if v.VersionNumber >= newVersionNum {
					newVersionNum = v.VersionNumber + 1
				}
			}
		} else {
			// Fallback if no versions loaded, assume 2 (since 1 was initial)?
			// Or assume 1 + 1.
			newVersionNum = 2
		}

		version := &domain.AgentVersion{
			AgentID:               agent.ID,
			VersionNumber:         newVersionNum,
			ConfigurationSnapshot: config,
			ReasonForChange:       "Configuration updated",
			CreatedBy:             &userID,
		}
		_ = s.agentRepo.CreateVersion(ctx, version)
	}

	return agent, nil
}

func (s *DefaultAgentService) DeleteAgent(ctx context.Context, id uuid.UUID) error {
	// Soft delete is handled by Repository/GORM
	return s.agentRepo.Delete(ctx, id)
}
