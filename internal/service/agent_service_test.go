package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAgentRepository
type MockAgentRepository struct {
	mock.Mock
}

func (m *MockAgentRepository) Create(ctx context.Context, agent *domain.Agent) error {
	args := m.Called(ctx, agent)
	// Simulate setting ID
	if agent.ID == uuid.Nil {
		agent.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockAgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.AgentStatus) ([]domain.Agent, error) {
	args := m.Called(ctx, orgID, status)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) Update(ctx context.Context, agent *domain.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockAgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAgentRepository) CreateVersion(ctx context.Context, version *domain.AgentVersion) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockAgentRepository) GetResources(ctx context.Context, agentID uuid.UUID) ([]domain.Resource, error) {
	args := m.Called(ctx, agentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *MockAgentRepository) GetAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]domain.User, error) {
	args := m.Called(ctx, agentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockAgentRepository) GetAssignedLLMs(ctx context.Context, agentID uuid.UUID) ([]domain.AgentLLM, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.AgentLLM), args.Error(1)
}

func (m *MockAgentRepository) GetAssignedAgents(ctx context.Context, userID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) GetAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.Application), args.Error(1)
}

func TestAgentService_CreateAgent(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAgentService(mockRepo)
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		name := "My Agent"
		config := json.RawMessage(`{"model": "gpt-4"}`)

		mockRepo.On("Create", ctx, mock.MatchedBy(func(a *domain.Agent) bool {
			return a.Name == name && a.OrganizationID == orgID
		})).Return(nil)

		mockRepo.On("CreateVersion", ctx, mock.MatchedBy(func(v *domain.AgentVersion) bool {
			return v.VersionNumber == 1
		})).Return(nil)

		agent, err := service.CreateAgent(ctx, orgID, userID, name, config)

		assert.NoError(t, err)
		assert.NotNil(t, agent)
		assert.Equal(t, name, agent.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty Name", func(t *testing.T) {
		_, err := service.CreateAgent(ctx, orgID, userID, "", nil)
		assert.Error(t, err)
	})

	t.Run("Duplicate Name", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)
		name := "Duplicate Agent"
		config := json.RawMessage(`{}`)

		// Simulate duplicate key error from DB
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("duplicate key value violates unique constraint"))

		_, err := service.CreateAgent(ctx, orgID, userID, name, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
	})
}

func TestAgentService_UpdateAgent(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	agentID := uuid.New()

	t.Run("Success - Config Changed", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		// Existing agent
		existingAgent := &domain.Agent{
			ID:             agentID,
			OrganizationID: orgID,
			Name:           "Old Name",
			Configuration:  json.RawMessage(`{"model": "gpt-3.5"}`),
			Versions: []domain.AgentVersion{
				{VersionNumber: 1},
			},
		}

		mockRepo.On("GetByID", ctx, agentID).Return(existingAgent, nil)

		newConfig := json.RawMessage(`{"model": "gpt-4"}`)
		newName := "New Name"

		mockRepo.On("Update", ctx, mock.MatchedBy(func(a *domain.Agent) bool {
			return a.Name == newName && string(a.Configuration) == string(newConfig)
		})).Return(nil)

		mockRepo.On("CreateVersion", ctx, mock.MatchedBy(func(v *domain.AgentVersion) bool {
			return v.VersionNumber == 2 && string(v.ConfigurationSnapshot) == string(newConfig)
		})).Return(nil)

		agent, err := service.UpdateAgent(ctx, agentID, userID, newName, newConfig, domain.AgentStatusActive)

		assert.NoError(t, err)
		assert.Equal(t, newName, agent.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - No Config Change", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		// Existing agent
		config := json.RawMessage(`{"model": "gpt-4"}`)
		existingAgent := &domain.Agent{
			ID:             agentID,
			OrganizationID: orgID,
			Name:           "Old Name",
			Configuration:  config,
			Versions: []domain.AgentVersion{
				{VersionNumber: 1},
			},
		}

		mockRepo.On("GetByID", ctx, agentID).Return(existingAgent, nil)

		newName := "New Name"

		mockRepo.On("Update", ctx, mock.MatchedBy(func(a *domain.Agent) bool {
			return a.Name == newName
		})).Return(nil)

		// CreateVersion should NOT be called

		agent, err := service.UpdateAgent(ctx, agentID, userID, newName, config, domain.AgentStatusActive)

		assert.NoError(t, err)
		assert.Equal(t, newName, agent.Name)
		mockRepo.AssertNotCalled(t, "CreateVersion")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetByID", ctx, agentID).Return(nil, nil)
		_, err := service.UpdateAgent(ctx, agentID, userID, "name", nil, domain.AgentStatusActive)
		assert.Error(t, err)
	})
}

func TestAgentService_GetAgent(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetByID", ctx, agentID).Return(&domain.Agent{ID: agentID}, nil)
		agent, err := service.GetAgent(ctx, agentID)
		assert.NoError(t, err)
		assert.NotNil(t, agent)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetByID", ctx, agentID).Return(nil, nil)
		_, err := service.GetAgent(ctx, agentID)
		assert.Error(t, err)
		assert.Equal(t, "agent not found", err.Error())
	})
}

func TestAgentService_ListAgentResources(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedResources := []domain.Resource{
			{ID: uuid.New(), Name: "Resource 1"},
			{ID: uuid.New(), Name: "Resource 2"},
		}

		mockRepo.On("GetResources", ctx, agentID).Return(expectedResources, nil)

		resources, err := service.ListAgentResources(ctx, agentID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedResources), len(resources))
		assert.Equal(t, expectedResources[0].Name, resources[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetResources", ctx, agentID).Return(nil, errors.New("db error"))

		_, err := service.ListAgentResources(ctx, agentID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_ListAssignedUsers(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedUsers := []domain.User{
			{ID: uuid.New(), Email: "user1@example.com"},
			{ID: uuid.New(), Email: "user2@example.com"},
		}

		mockRepo.On("GetAssignedUsers", ctx, agentID).Return(expectedUsers, nil)

		users, err := service.ListAssignedUsers(ctx, agentID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedUsers), len(users))
		assert.Equal(t, expectedUsers[0].Email, users[0].Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetAssignedUsers", ctx, agentID).Return([]domain.User{}, errors.New("db error"))

		_, err := service.ListAssignedUsers(ctx, agentID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_GetAgentLLMs(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedLLMs := []domain.AgentLLM{
			{ID: uuid.New(), AgentID: agentID, LLMModel: domain.LLMModel{FamilyName: "GPT-4"}},
		}
		mockRepo.On("GetAssignedLLMs", ctx, agentID).Return(expectedLLMs, nil)

		llms, err := service.GetAgentLLMs(ctx, agentID)
		assert.NoError(t, err)
		assert.Len(t, llms, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_ListAssignedApplications(t *testing.T) {
	ctx := context.Background()
	agentID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedApps := []domain.Application{
			{Name: "App A"},
		}
		mockRepo.On("GetAssignedApplications", ctx, agentID).Return(expectedApps, nil)

		apps, err := service.ListAssignedApplications(ctx, agentID)
		assert.NoError(t, err)
		assert.Len(t, apps, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetAssignedApplications", ctx, agentID).Return([]domain.Application{}, errors.New("db error"))

		_, err := service.ListAssignedApplications(ctx, agentID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_ListAssignedAgents(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedAgents := []domain.Agent{
			{Name: "Agent X"},
		}
		mockRepo.On("GetAssignedAgents", ctx, userID).Return(expectedAgents, nil)

		agents, err := service.ListAssignedAgents(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, agents, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		mockRepo.On("GetAssignedAgents", ctx, userID).Return([]domain.Agent{}, errors.New("db error"))

		_, err := service.ListAssignedAgents(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_ListAgentsByStatus(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		expectedAgents := []domain.Agent{
			{Name: "Active Agent", Status: domain.AgentStatusActive},
		}
		mockRepo.On("ListByStatus", ctx, orgID, domain.AgentStatusActive).Return(expectedAgents, nil)

		agents, err := service.ListAgentsByStatus(ctx, orgID, domain.AgentStatusActive)
		assert.NoError(t, err)
		assert.Len(t, agents, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestAgentService_GetActiveMonthlyCost(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAgentRepository)
		service := NewAgentService(mockRepo)

		activeAgents := []domain.Agent{
			{Name: "Monthly Agent", Status: domain.AgentStatusActive, BillingCycle: domain.BillingCycleMonthly, CostAmount: 100.0},
			{Name: "Yearly Agent", Status: domain.AgentStatusActive, BillingCycle: domain.BillingCycleYearly, CostAmount: 1200.0},
			// OneTime should be ignored
			{Name: "OneTime Agent", Status: domain.AgentStatusActive, BillingCycle: domain.BillingCycleOneTime, CostAmount: 500.0},
		}
		mockRepo.On("ListByStatus", ctx, orgID, domain.AgentStatusActive).Return(activeAgents, nil)

		cost, err := service.GetActiveMonthlyCost(ctx, orgID)
		assert.NoError(t, err)
		// 100 (Monthly) + 1200/12 (Yearly=100) + 0 (OneTime) = 200
		assert.Equal(t, 200.0, cost)
		mockRepo.AssertExpectations(t)
	})
}
