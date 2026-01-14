package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agentXmap/internal/domain"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestAgent_E2E(t *testing.T) {
	// 1. Setup Server
	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// 2. Setup Expect Client
	e := httpexpect.Default(t, server.URL)

	// 3. Seed Data
	db := GetTestDB()
	agentID := seedAgentData(t, db)

	// 4. Test Scenarios

	// A. Create Agent
	// Mock Header UserID
	userID := uuid.New()

	e.POST("/api/v1/agents").
		WithHeader("X-User-ID", userID.String()).
		WithJSON(map[string]interface{}{
			"name": "New E2E Agent",
			"configuration": map[string]interface{}{
				"setting": "enabled",
			},
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("data").Object().
		ContainsKey("id").
		ContainsKey("name").ValueEqual("name", "New E2E Agent")

	// B. Get Agent (Seeded)
	e.GET("/api/v1/agents/"+agentID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		ContainsKey("name").ValueEqual("name", "Seeded Agent")

	// C. List Certifications
	// Expecting seeded certs
	certsObj := e.GET("/api/v1/agents/" + agentID.String() + "/certifications").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object()

	certsObj.ContainsKey("agent_certifications").
		Value("agent_certifications").Array().NotEmpty()

	certsObj.ContainsKey("llm_certifications").
		Value("llm_certifications").Object().
		ContainsKey("gpt-4-turbo").Value("gpt-4-turbo").Array().NotEmpty()
}

func seedAgentData(t *testing.T, db *gorm.DB) uuid.UUID {
	// Create Agent
	agent := domain.Agent{
		ID:            uuid.New(),
		Name:          "Seeded Agent",
		Configuration: []byte("{}"),
		Status:        domain.AgentStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&agent).Error; err != nil {
		t.Fatalf("Failed to seed agent: %v", err)
	}

	// Create Certification
	cert := domain.Certification{
		ID:               uuid.New(),
		Name:             "E2E Agent Cert",
		IssuingAuthority: "E2E Auth",
	}
	if err := db.Create(&cert).Error; err != nil {
		t.Fatalf("Failed to seed cert: %v", err)
	}

	// Link Agent Cert
	agentCert := domain.AgentCertification{
		AgentID:         agent.ID,
		CertificationID: cert.ID,
	}
	if err := db.Create(&agentCert).Error; err != nil {
		t.Fatalf("Failed to seed agent cert: %v", err)
	}

	// Create LLM Provider
	llmProvider := domain.LLMProvider{
		ID:   uuid.New(),
		Name: "OpenAI",
	}
	if err := db.Create(&llmProvider).Error; err != nil {
		t.Fatalf("Failed to seed llm provider: %v", err)
	}

	// Create LLM Model & Cert
	llmModel := domain.LLMModel{
		ID:           uuid.New(),
		ProviderID:   llmProvider.ID,
		FamilyName:   "GPT-4",
		VersionName:  "Turbo",
		ApiModelName: "gpt-4-turbo",
	}
	if err := db.Create(&llmModel).Error; err != nil {
		t.Fatalf("Failed to seed llm model: %v", err)
	}

	llmCert := domain.Certification{
		ID:               uuid.New(),
		Name:             "E2E LLM Cert",
		IssuingAuthority: "E2E Auth",
	}
	if err := db.Create(&llmCert).Error; err != nil {
		t.Fatalf("Failed to seed llm cert: %v", err)
	}

	llmModelCert := domain.LLMModelCertification{
		LLMModelID:      llmModel.ID,
		CertificationID: llmCert.ID,
	}
	if err := db.Create(&llmModelCert).Error; err != nil {
		t.Fatalf("Failed to seed llm model cert: %v", err)
	}

	// Link Agent to LLM
	agentLLM := domain.AgentLLM{
		AgentID:    agent.ID,
		LLMModelID: llmModel.ID,
		IsPrimary:  true,
	}
	if err := db.Create(&agentLLM).Error; err != nil {
		t.Fatalf("Failed to seed agent llm: %v", err)
	}

	return agent.ID
}
