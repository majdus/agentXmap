package repository

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuditRepository_CreateLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewAuditRepository(gormDB)
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		log := &domain.SystemAuditLog{
			OrganizationID: orgID,
			ActorUserID:    &userID,
			EntityType:     "agent",
			EntityID:       entityID,
			Action:         domain.AuditActionCreate,
			Changes:        json.RawMessage(`{}`),
			IPAddress:      "127.0.0.1",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "system_audit_logs"`)).
			WithArgs(orgID, userID, "agent", entityID, "create", log.Changes, "127.0.0.1").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := repo.CreateLog(ctx, log)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuditRepository_CreateExecution(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewAuditRepository(gormDB)
	ctx := context.Background()
	orgID := uuid.New()
	agentID := uuid.New()
	versionID := uuid.New()
	modelID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		exec := &domain.AgentExecution{
			OrganizationID:   orgID,
			AgentID:          agentID,
			AgentVersionID:   versionID,
			LLMModelID:       modelID,
			Status:           "completed",
			LatencyMs:        100,
			TokenUsageInput:  50,
			TokenUsageOutput: 50,
			IsPIIDetected:    false,
			SafetyScore:      1.0,
			CreatedAt:        time.Now(),
		}

		mock.ExpectBegin()
		// GORM with Postgres uses Query for INSERT ... RETURNING
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agent_executions"`)).
			WithArgs(orgID, agentID, versionID, modelID, nil, nil, "completed", 100, 50, 50, false, 1.0, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(uuid.New(), time.Now()))
		mock.ExpectCommit()

		err := repo.CreateExecution(ctx, exec)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
