package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"agentXmap/internal/domain"
	"agentXmap/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgentHandler struct {
	agentService service.AgentService
	llmService   service.LLMService
}

func NewAgentHandler(agentService service.AgentService, llmService service.LLMService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		llmService:   llmService,
	}
}

// CreateAgentRequest DTO
type CreateAgentRequest struct {
	Name          string          `json:"name" binding:"required" example:"Support Agent"`
	Configuration json.RawMessage `json:"configuration"`
}

// UpdateAgentRequest DTO
type UpdateAgentRequest struct {
	Name          string             `json:"name" example:"New Agent Name"`
	Configuration json.RawMessage    `json:"configuration"`
	Status        domain.AgentStatus `json:"status" example:"active"`
}

// CreateAgent godoc
// @Summary      Create a new agent
// @Description  Create a new AI agent with a specific configuration.
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        X-User-ID  header    string  false  "User ID for attribution"
// @Param        request    body      CreateAgentRequest  true  "Agent details"
// @Success      201  {object}  Response{data=domain.Agent}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	// TODO: Get UserID from context
	// userID := c.GetString("userID")
	// Mocking UserID for now as we don't have the middleware setup in this context
	userID := uuid.Nil // Or generate a new one, or require header. Let's assume NIL or header for now.
	if idStr := c.GetHeader("X-User-ID"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			userID = id
		}
	}

	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent, err := h.agentService.CreateAgent(c.Request.Context(), userID, req.Name, req.Configuration)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to create agent", err.Error())
		return
	}

	RespondCreated(c, agent)
}

// GetAgent godoc
// @Summary      Get agent details
// @Description  Retrieve full details of an agent by ID.
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response{data=domain.Agent}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id} [get]
func (h *AgentHandler) GetAgent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	agent, err := h.agentService.GetAgent(c.Request.Context(), id)
	if err != nil {
		// Differentiate between 404 and 500
		RespondError(c, http.StatusInternalServerError, "Failed to get agent", err.Error())
		return
	}

	RespondSuccess(c, agent)
}

// ListAgents godoc
// @Summary      List all agents
// @Description  Retrieve a list of all available agents.
// @Tags         agents
// @Produce      json
// @Success      200  {object}  Response{data=[]domain.Agent}
// @Failure      500  {object}  Response
// @Router       /agents [get]
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents, err := h.agentService.ListAgents(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list agents", err.Error())
		return
	}
	RespondSuccess(c, agents)
}

// UpdateAgent godoc
// @Summary      Update an agent
// @Description  Modify name, configuration or status of an existing agent.
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        id         path      string  true  "Agent ID"
// @Param        X-User-ID  header    string  false  "User ID for attribution"
// @Param        request    body      UpdateAgentRequest  true  "Update details"
// @Success      200  {object}  Response{data=domain.Agent}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id} [patch]
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	// TODO: Get UserID from context
	userID := uuid.Nil
	if idStr := c.GetHeader("X-User-ID"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			userID = id
		}
	}

	var req UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	agent, err := h.agentService.UpdateAgent(c.Request.Context(), id, userID, req.Name, req.Configuration, req.Status)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to update agent", err.Error())
		return
	}

	RespondSuccess(c, agent)
}

// DeleteAgent godoc
// @Summary      Delete an agent
// @Description  Remove an agent from the system.
// @Tags         agents
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id} [delete]
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	if err := h.agentService.DeleteAgent(c.Request.Context(), id); err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to delete agent", err.Error())
		return
	}

	RespondSuccess(c, nil)
}

// ListAgentResources godoc
// @Summary      List agent resources
// @Description  Retrieve all resources (DBs, APIs) accessible by this agent.
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response{data=[]domain.Resource}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id}/resources [get]
func (h *AgentHandler) ListAgentResources(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	resources, err := h.agentService.ListAgentResources(c.Request.Context(), id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list agent resources", err.Error())
		return
	}

	RespondSuccess(c, resources)
}

// ListAssignedUsers godoc
// @Summary      List assigned users
// @Description  Retrieve all users who have access to this agent.
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response{data=[]domain.User}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id}/users [get]
func (h *AgentHandler) ListAssignedUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	users, err := h.agentService.ListAssignedUsers(c.Request.Context(), id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list assigned users", err.Error())
		return
	}

	RespondSuccess(c, users)
}

// ListAssignedApplications godoc
// @Summary      List assigned applications
// @Description  Retrieve all external applications authorized to use this agent.
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response{data=[]domain.Application}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id}/applications [get]
func (h *AgentHandler) ListAssignedApplications(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	apps, err := h.agentService.ListAssignedApplications(c.Request.Context(), id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list assigned applications", err.Error())
		return
	}

	RespondSuccess(c, apps)
}

// ListAgentCertifications godoc
// @Summary      List agent certifications
// @Description  Retrieve certifications for both the agent and its associated LLM models.
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  Response{data=object{agent_certifications=[]domain.Certification,llm_certifications=map[string][]domain.Certification}}
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /agents/{id}/certifications [get]
func (h *AgentHandler) ListAgentCertifications(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid Agent ID", err.Error())
		return
	}

	ctx := c.Request.Context()

	// 1. Get Agent Certifications
	agentCerts, err := h.agentService.ListAgentCertifications(ctx, id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list agent certifications", err.Error())
		return
	}

	// 2. Get Agent LLMs
	agentLLMs, err := h.agentService.GetAgentLLMs(ctx, id)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Failed to list agent LLMs", err.Error())
		return
	}

	// 3. Get Certifications for each LLM
	llmCertsMap := make(map[string][]domain.Certification)
	for _, agentLLM := range agentLLMs {
		certs, err := h.llmService.ListModelCertifications(ctx, agentLLM.LLMModelID)
		if err != nil {
			// Log error but continue? Or fail? Let's log and continue for partial result or fail.
			// Ideally we shouldn't fail everything if one LLM lookup fails, but for now strict consistency.
			// Actually, let's just log (fmt.Printf for now) and continue.
			fmt.Printf("Failed to get certs for model %s: %v\n", agentLLM.LLMModelID, err)
			continue
		}

		modelName := agentLLM.LLMModel.ApiModelName
		if modelName == "" {
			// If model is not loaded with details, might need to fetch it?
			// The Repo GetAssignedLLMs does Preload("LLMModel"). So it should be there.
			modelName = "Unknown Model"
		}
		llmCertsMap[modelName] = certs
	}

	RespondSuccess(c, gin.H{
		"agent_certifications": agentCerts,
		"llm_certifications":   llmCertsMap,
	})
}
