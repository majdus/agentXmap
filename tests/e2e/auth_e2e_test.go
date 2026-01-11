package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"agentXmap/internal/handler"
	"agentXmap/internal/repository"
	"agentXmap/internal/service"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// setupRouter initializes the app with the Test DB
func setupRouter() *gin.Engine {
	// Initialize Repos & Services using the Test DB from setup_test.go
	db := GetTestDB()
	userRepo := repository.NewUserRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	identityService := service.NewIdentityService(userRepo, invitationRepo)
	authHandler := handler.NewAuthHandler(identityService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	handler.RegisterRoutes(r, authHandler)
	return r
}

func TestAuth_E2E(t *testing.T) {
	// 1. Setup Server
	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// 2. Setup Expect Client
	e := httpexpect.Default(t, server.URL)

	// 3. Scenario: Invitation Flow
	// A. Bootstrap Admin (Done manualy here, usually via config/migration)
	db := GetTestDB()
	adminEmail := "admin@example.com"
	adminID := uuid.New()
	// Create Admin directly in DB to simulate bootstrap
	db.Exec("INSERT INTO users (id, email, password_hash, role) VALUES (?, ?, ?, 'admin')",
		adminID, adminEmail, "$2a$14$P.u.something.hash") // Mock hash

	// B. Admin Invites User
	// We need to pass X-Admin-ID header as per our temporary handler implementation
	userEmail := "bob@example.com"

	obj := e.POST("/api/v1/auth/invite").
		WithHeader("X-Admin-ID", adminID.String()).
		WithJSON(gin.H{
			"email": userEmail,
			"role":  "user",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("data").Array().First().Object()

	obj.ContainsKey("email").ValueEqual("email", userEmail)
	obj.ContainsKey("token").Value("token").String().NotEmpty()

	inviteToken := obj.Value("token").String().Raw()

	// C. User Accepts Invitation
	password := "newsecurepassword"
	e.POST("/api/v1/auth/accept-invitation").
		WithJSON(gin.H{
			"token":      inviteToken,
			"password":   password,
			"first_name": "Bob",
			"last_name":  "Builder",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").ValueEqual("success", true)

	// D. User Logins
	e.POST("/api/v1/auth/login").
		WithJSON(gin.H{
			"email":    userEmail,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").ValueEqual("success", true).
		ContainsKey("data").
		Path("$.data.token").String().NotEmpty()
}
