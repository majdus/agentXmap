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

	// 3. Scenario: Register and Login
	email := "e2e_test@example.com"
	password := "securepassword123"

	// Step A: Register
	e.POST("/api/v1/auth/register").
		WithJSON(gin.H{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("success").ValueEqual("success", true).
		ContainsKey("data").
		Path("$.data.email").String().Equal(email)

	// Step B: Login
	e.POST("/api/v1/auth/login").
		WithJSON(gin.H{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").ValueEqual("success", true).
		ContainsKey("data").
		Path("$.data.token").String().NotEmpty()
}
