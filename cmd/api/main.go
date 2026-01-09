package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// 3. Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	// TODO: Add custom logger middleware

	// 4. Routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"version":   cfg.App.Version,
				"timestamp": time.Now().Unix(),
			})
		})
	}

	// 5. Start Server
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
