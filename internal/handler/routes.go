package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, agentHandler *AgentHandler) {
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/invite", authHandler.InviteUser)
			auth.POST("/accept-invitation", authHandler.AcceptInvitation)
			auth.POST("/login", authHandler.Login)
		}

		agents := api.Group("/agents")
		{
			agents.POST("", agentHandler.CreateAgent)
			agents.GET("", agentHandler.ListAgents)
			agents.GET("/:id", agentHandler.GetAgent)
			agents.PUT("/:id", agentHandler.UpdateAgent)
			agents.DELETE("/:id", agentHandler.DeleteAgent)
			agents.GET("/:id/resources", agentHandler.ListAgentResources)
			agents.GET("/:id/users", agentHandler.ListAssignedUsers)
			agents.GET("/:id/applications", agentHandler.ListAssignedApplications)
			agents.GET("/:id/certifications", agentHandler.ListAgentCertifications)
		}
	}
}
