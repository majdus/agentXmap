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
	Email string          `json:"email" binding:"required,email" example:"user@example.com"`
	Role  domain.UserRole `json:"role" binding:"required" example:"user"`
}

// AcceptInvitationRequest DTO
type AcceptInvitationRequest struct {
	Token     string `json:"token" binding:"required" example:"abc-123-token"`
	Password  string `json:"password" binding:"required,min=8" example:"securepassword123"`
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// InviteUser godoc
// @Summary      Invite a new user
// @Description  Create a new invitation for a user with a specific role. Admin only.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Admin-ID  header    string  true  "Admin User ID"
// @Param        request     body      InviteUserRequest  true  "Invitation details"
// @Success      201  {object}  Response{data=[]domain.Invitation}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router       /auth/invite [post]
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

// AcceptInvitation godoc
// @Summary      Accept an invitation
// @Description  Complete user registration by providing a password and personal details using the invitation token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      AcceptInvitationRequest  true  "Acceptance details"
// @Success      200  {object}  Response{data=object{id=string,email=string}}
// @Failure      400  {object}  Response
// @Router       /auth/accept-invitation [post]
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

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return a mock JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200  {object}  Response{data=object{token=string,user=domain.User}}
// @Failure      400  {object}  Response
// @Failure      401  {object}  Response
// @Router       /auth/login [post]
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
