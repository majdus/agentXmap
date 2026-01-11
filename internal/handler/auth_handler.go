package handler

import (
	"fmt"
	"net/http"

	"agentXmap/internal/domain"
	"agentXmap/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	identityService service.IdentityService
}

func NewAuthHandler(identityService service.IdentityService) *AuthHandler {
	return &AuthHandler{
		identityService: identityService,
	}
}

// InviteUserRequest DTO
type InviteUserRequest struct {
	Email string          `json:"email" binding:"required,email"`
	Role  domain.UserRole `json:"role" binding:"required"`
}

// AcceptInvitationRequest DTO
type AcceptInvitationRequest struct {
	Token     string `json:"token" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// InviteUser handler (Admin only)
func (h *AuthHandler) InviteUser(c *gin.Context) {
	// TODO: Get Admin ID from context (Middleware should set this)
	// For now, we assume a header "X-User-ID" exists or we mock it for the test
	// In a real scenario, this comes from the JWT Middleare: userID := c.GetString("userID")

	// Mocking Admin ID extraction for vertical slice demonstration
	// This MUST be replaced by middleware later
	adminIDStr := c.GetHeader("X-Admin-ID")
	if adminIDStr == "" {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "Missing X-Admin-ID header")
		return
	}
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "Invalid Admin ID")
		return
	}

	var req InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	invitations, err := h.identityService.InviteUsers(c.Request.Context(), adminID, []string{req.Email}, req.Role)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "Invitation failed", err.Error())
		return
	}

	RespondCreated(c, invitations)
}

// AcceptInvitation handler
func (h *AuthHandler) AcceptInvitation(c *gin.Context) {
	var req AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.identityService.AcceptInvitation(c.Request.Context(), req.Token, req.Password, req.FirstName, req.LastName)
	if err != nil {
		fmt.Printf("ACCEPT INVITATION FAILED: Token=%s, Err=%v\n", req.Token, err)
		RespondError(c, http.StatusBadRequest, "Accept invitation failed", err.Error())
		return
	}

	RespondSuccess(c, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login handler
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.identityService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// DEBUG logging for E2E test failure
		fmt.Printf("LOGIN FAILED: Email=%s, Err=%v\n", req.Email, err)
		RespondError(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	// TODO: Generate JWT Token here. For now, return User ID.
	RespondSuccess(c, gin.H{
		"token": "mock-jwt-token-" + user.ID.String(),
		"user":  user,
	})
}
