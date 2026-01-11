package repository

import (
	"fmt"
	"time"

	"agentXmap/internal/domain"
	"agentXmap/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection using GORM.
func InitDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password, // Ensure this matches struct field
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	var logLevel logger.LogLevel
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection Pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate applies schema changes.
// Note: In production, prefer using the SQL scripts via Makefile.
// This is useful for dev environments or testing.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Agent{},
		&domain.AgentVersion{},
		&domain.AgentAssignment{},
		&domain.LLMProvider{},
		&domain.LLMModel{},
		&domain.AgentLLM{},
		&domain.Application{},
		&domain.ApplicationKey{},
		&domain.ApplicationAgentAccess{},
		&domain.ResourceType{},
		&domain.Resource{},
		&domain.ResourceSecret{},
		&domain.AgentResourceAccess{},
		&domain.Certification{},
		&domain.AgentCertification{},
		&domain.SystemAuditLog{},
		// &domain.AgentExecution{}, // Partitioned table often skipped in auto-migrate or handled carefully
	)
}
