package repository

import (
	"context"
	"errors"
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

	tests := []struct {
		name    string
		input   *domain.SystemAuditLog
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			input: log,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "system_audit_logs"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // 8 args
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(log.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Error",
			input: log,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "system_audit_logs"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.CreateLog(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
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

	tests := []struct {
		name    string
		input   *domain.AgentExecution
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			input: exec,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agent_executions"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(exec.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Error",
			input: exec,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agent_executions"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.CreateExecution(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
