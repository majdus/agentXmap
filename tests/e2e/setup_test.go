package e2e

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"agentXmap/internal/repository"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 2. Get Connection String
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// 3. Init GORM from the container connection string
	// Direct connection for test avoids parsing the DSN back into config parts
	testDB, err = gorm.Open(repository.GetPostgresDialector(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test db: %s", err)
	}

	// 4. Run Migrations
	// Create Enums first
	testDB.Exec("CREATE TYPE user_role AS ENUM ('admin', 'manager', 'user')")
	testDB.Exec("CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'revoked')")
	testDB.Exec("CREATE TYPE agent_status AS ENUM ('active', 'inactive', 'maintenance', 'deprecated')")
	testDB.Exec("CREATE TYPE billing_cycle AS ENUM ('monthly', 'yearly', 'one_time', 'custom')")
	testDB.Exec("CREATE TYPE access_level AS ENUM ('read_only', 'read_write')")
	testDB.Exec("CREATE TYPE audit_action AS ENUM ('create', 'update', 'delete', 'login', 'export_data')")

	// Ensure we are in project root or can find migrations if they are SQL files.
	// Since we use AutoMigrate from GORM models, we just call it.
	if err := repository.AutoMigrate(testDB); err != nil {
		log.Fatalf("failed to migrate test db: %s", err)
	}

	// 5. Run Tests
	code := m.Run()

	// 6. Cleanup
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

// GetTestDB returns the gorm instance connected to the test container
func GetTestDB() *gorm.DB {
	return testDB
}
