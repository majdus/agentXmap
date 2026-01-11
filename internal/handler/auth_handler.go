package handler

import (
	"net/http"

	"agentXmap/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	identityService service.IdentityService
}

func NewAuthHandler(identityService service.IdentityService) *AuthHandler {
	return &AuthHandler{
		identityService: identityService,
	}
}

// RegisterRequest DTO
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handler
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.identityService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// TODO: Better error mapping (e.g. check if user already exists)
		RespondError(c, http.StatusInternalServerError, "Registration failed", err.Error())
		return
	}

	RespondCreated(c, gin.H{
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
		RespondError(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	// TODO: Generate JWT Token here. For now, return User ID.
	RespondSuccess(c, gin.H{
		"token": "mock-jwt-token-" + user.ID.String(),
		"user":  user,
	})
}
