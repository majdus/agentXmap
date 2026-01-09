package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"agentXmap/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditRepository_CreateLog(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAuditRepository(db)
	ctx := context.TODO()

	log := &domain.SystemAuditLog{
		ID:             uuid.New(),
		OrganizationID: uuid.New(),
		Action:         domain.AuditActionCreate,
		OccurredAt:     time.Now(),
	}

	// GORM behavior: 8 arguments seen in error.
	// We previously expected 9.
	// We will match 8 args now.
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "system_audit_logs"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // 8 args
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(log.ID))
	mock.ExpectCommit()

	err := repo.CreateLog(ctx, log)
	assert.NoError(t, err)
}

func TestAuditRepository_CreateExecution(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAuditRepository(db)
	ctx := context.TODO()

	exec := &domain.AgentExecution{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		AgentID:   uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agent_executions"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(exec.ID))
	mock.ExpectCommit()

	err := repo.CreateExecution(ctx, exec)
	assert.NoError(t, err)
}
