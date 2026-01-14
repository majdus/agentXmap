package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"agentXmap/internal/domain"
	"agentXmap/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAgentService
type MockAgentService struct {
	mock.Mock
}

func (m *MockAgentService) CreateAgent(ctx context.Context, userID uuid.UUID, name string, config json.RawMessage) (*domain.Agent, error) {
	args := m.Called(ctx, userID, name, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agent), args.Error(1)
}

func (m *MockAgentService) GetAgent(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agent), args.Error(1)
}

func (m *MockAgentService) ListAgents(ctx context.Context) ([]domain.Agent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentService) ListAgentsByStatus(ctx context.Context, status domain.AgentStatus) ([]domain.Agent, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentService) GetActiveMonthlyCost(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAgentService) ListAgentResources(ctx context.Context, agentID uuid.UUID) ([]domain.Resource, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.Resource), args.Error(1)
}

func (m *MockAgentService) ListAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]domain.User, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockAgentService) ListAssignedAgents(ctx context.Context, userID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockAgentService) GetAgentLLMs(ctx context.Context, agentID uuid.UUID) ([]domain.AgentLLM, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.AgentLLM), args.Error(1)
}

func (m *MockAgentService) ListAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.Application), args.Error(1)
}

func (m *MockAgentService) ListAgentCertifications(ctx context.Context, agentID uuid.UUID) ([]domain.Certification, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]domain.Certification), args.Error(1)
}

func (m *MockAgentService) UpdateAgent(ctx context.Context, id, userID uuid.UUID, name string, config json.RawMessage, status domain.AgentStatus) (*domain.Agent, error) {
	args := m.Called(ctx, id, userID, name, config, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agent), args.Error(1)
}

func (m *MockAgentService) DeleteAgent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLLMService
type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) ListProviders(ctx context.Context) ([]domain.LLMProvider, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.LLMProvider), args.Error(1)
}

func (m *MockLLMService) ListModels(ctx context.Context, providerID uuid.UUID) ([]domain.LLMModel, error) {
	args := m.Called(ctx, providerID)
	return args.Get(0).([]domain.LLMModel), args.Error(1)
}

func (m *MockLLMService) GetModel(ctx context.Context, id uuid.UUID) (*domain.LLMModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LLMModel), args.Error(1)
}

func (m *MockLLMService) ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]domain.Agent, error) {
	args := m.Called(ctx, modelID)
	return args.Get(0).([]domain.Agent), args.Error(1)
}

func (m *MockLLMService) ListModelCertifications(ctx context.Context, modelID uuid.UUID) ([]domain.Certification, error) {
	args := m.Called(ctx, modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Certification), args.Error(1)
}

func TestAgentHandler_CreateAgent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.POST("/agents", h.CreateAgent)

	userID := uuid.New()
	agentName := "Test Agent"
	config := json.RawMessage(`{"key":"value"}`)

	expectedAgent := &domain.Agent{
		ID:            uuid.New(),
		Name:          agentName,
		Configuration: config,
		CreatedBy:     &userID,
	}

	mockAgentService.On("CreateAgent", mock.Anything, mock.Anything, agentName, config).Return(expectedAgent, nil)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"name":          agentName,
		"configuration": config,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	// Header for UserID mock
	req.Header.Set("X-User-ID", userID.String())

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	dataMap := response.Data.(map[string]interface{})
	assert.Equal(t, agentName, dataMap["name"])
}

func TestAgentHandler_GetAgent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id", h.GetAgent)

	agentID := uuid.New()
	expectedAgent := &domain.Agent{
		ID:   agentID,
		Name: "Existing Agent",
	}

	mockAgentService.On("GetAgent", mock.Anything, agentID).Return(expectedAgent, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String(), nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestAgentHandler_GetAgent_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id", h.GetAgent)

	agentID := uuid.New()

	// Mock Service Error
	mockAgentService.On("GetAgent", mock.Anything, agentID).Return(nil, errors.New("agent not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String(), nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Note: Handler implementation currently maps all errors to 500. Ideally should be 404.
	// Verifying current behavior first.
}

func TestAgentHandler_ListAssignedUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id/users", h.ListAssignedUsers)

	agentID := uuid.New()
	expectedUsers := []domain.User{
		{ID: uuid.New(), Email: "user1@example.com"},
		{ID: uuid.New(), Email: "user2@example.com"},
	}

	mockAgentService.On("ListAssignedUsers", mock.Anything, agentID).Return(expectedUsers, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String()+"/users", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	dataList := response.Data.([]interface{})
	assert.Len(t, dataList, 2)
}

func TestAgentHandler_ListAssignedApplications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id/applications", h.ListAssignedApplications)

	agentID := uuid.New()
	expectedApps := []domain.Application{
		{ID: uuid.New(), Name: "App 1"},
	}

	mockAgentService.On("ListAssignedApplications", mock.Anything, agentID).Return(expectedApps, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String()+"/applications", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	dataList := response.Data.([]interface{})
	assert.Len(t, dataList, 1)
}

func TestAgentHandler_ListAgentCertifications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id/certifications", h.ListAgentCertifications)

	agentID := uuid.New()
	llmModelID := uuid.New()

	expectedAgentCerts := []domain.Certification{
		{Name: "ISO 27001", IssuingAuthority: "ISO"},
	}

	expectedAgentLLMs := []domain.AgentLLM{
		{
			LLMModelID: llmModelID,
			LLMModel: domain.LLMModel{
				ApiModelName: "GPT-4",
			},
		},
	}

	expectedLLMCerts := []domain.Certification{
		{Name: "SOC 2", IssuingAuthority: "AICPA"},
	}

	// Mock expectations
	mockAgentService.On("ListAgentCertifications", mock.Anything, agentID).Return(expectedAgentCerts, nil)
	mockAgentService.On("GetAgentLLMs", mock.Anything, agentID).Return(expectedAgentLLMs, nil)
	mockLLMService.On("ListModelCertifications", mock.Anything, llmModelID).Return(expectedLLMCerts, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String()+"/certifications", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	dataMap := response.Data.(map[string]interface{})

	// Check Agent Certs
	agentCerts := dataMap["agent_certifications"].([]interface{})
	assert.Len(t, agentCerts, 1)
	assert.Equal(t, "ISO 27001", agentCerts[0].(map[string]interface{})["name"])

	// Check LLM Certs
	llmCertsMap := dataMap["llm_certifications"].(map[string]interface{})
	assert.Len(t, llmCertsMap, 1)
	gpt4Certs := llmCertsMap["GPT-4"].([]interface{})
	assert.Len(t, gpt4Certs, 1)
	assert.Equal(t, "SOC 2", gpt4Certs[0].(map[string]interface{})["name"])
}

func TestAgentHandler_ListAgentCertifications_PartialLLMFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	mockLLMService := new(MockLLMService)
	h := handler.NewAgentHandler(mockAgentService, mockLLMService)

	r := gin.New()
	r.GET("/agents/:id/certifications", h.ListAgentCertifications)

	agentID := uuid.New()
	llmModelID := uuid.New()

	expectedAgentCerts := []domain.Certification{} // Empty agent certs

	expectedAgentLLMs := []domain.AgentLLM{
		{
			LLMModelID: llmModelID,
			LLMModel: domain.LLMModel{
				ApiModelName: "GPT-4",
			},
		},
	}

	// Mock expectations
	mockAgentService.On("ListAgentCertifications", mock.Anything, agentID).Return(expectedAgentCerts, nil)
	mockAgentService.On("GetAgentLLMs", mock.Anything, agentID).Return(expectedAgentLLMs, nil)

	// Mock LLM Service failure
	mockLLMService.On("ListModelCertifications", mock.Anything, llmModelID).Return(nil, errors.New("failed upstream"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agents/"+agentID.String()+"/certifications", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) // Should still succeed, just missing LLM certs

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	dataMap := response.Data.(map[string]interface{})

	// Check Agent Certs (Empty)
	agentCerts := dataMap["agent_certifications"].([]interface{})
	assert.Len(t, agentCerts, 0)

	// Check LLM Certs - Should contain the key but maybe empty list?
	// Our logic: `llmCertsMap[modelName] = certs`. If error, we continue.
	// So `llmCertsMap` will NOT have the entry for "GPT-4".
	llmCertsMap := dataMap["llm_certifications"].(map[string]interface{})
	_, exists := llmCertsMap["GPT-4"]
	assert.False(t, exists, "GPT-4 should be missing from certs map due to error")
}
