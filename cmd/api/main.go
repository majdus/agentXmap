package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agentXmap/internal/handler"
	"agentXmap/internal/repository"
	"agentXmap/internal/service"
	"agentXmap/pkg/config"
	"agentXmap/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title agentXmap API
// @version 1.0.0
// @description Backend API for agentXmap project.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init Logger
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Encoding); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Log.Info("Starting agentXmap API", zap.String("version", cfg.App.Version), zap.String("env", cfg.Server.Mode))

	// 3. Database Connection
	db, err := repository.InitDB(*cfg)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Migrate database schema
	if err := repository.AutoMigrate(db); err != nil {
		logger.Log.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 4. Init Repositories & Services
	userRepo := repository.NewUserRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	identityService := service.NewIdentityService(userRepo, invitationRepo)

	// 5. Initial Admin Creation
	if cfg.InitialAdmin.Email != "" && cfg.InitialAdmin.Password != "" {
		ctx := context.Background()
		logger.Log.Info("Attempting to create initial admin user...", zap.String("email", cfg.InitialAdmin.Email))
		user, err := identityService.SignUp(ctx, cfg.InitialAdmin.Email, cfg.InitialAdmin.Password)
		if err != nil {
			if err.Error() == "user already exists" {
				logger.Log.Info("Initial admin user already exists", zap.String("email", cfg.InitialAdmin.Email))
			} else {
				logger.Log.Error("Failed to create initial admin user", zap.Error(err))
				// Optional: Exit on failure? For now, we continue.
			}
		} else {
			logger.Log.Info("Initial admin user created successfully", zap.String("id", user.ID.String()))
		}
	}

	// 6. Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	// TODO: Add custom logger middleware

	// 7. Handlers & Routes
	authHandler := handler.NewAuthHandler(identityService)
	handler.RegisterRoutes(r, authHandler)

	// Add Health Check (Keep simple health check here or move to a SystemHandler)
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"version":   cfg.App.Version,
			"timestamp": time.Now().Unix(),
		})
	})

	// 8. Start Server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Listen:", zap.Error(err))
		}
	}()
	logger.Log.Info("Server listening", zap.String("port", cfg.Server.Port))

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Log.Info("Server exiting")
}
