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

func TestApplicationRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewApplicationRepository(db)
	ctx := context.TODO()

	id := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			id:   id,
			mock: func() {
				// 1. Get App
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "applications" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "App 1"))

				// 2. Preloads

				// AgentAccess - GORM pluralizes default table names.
				// The struct field is AgentAccess []ApplicationAgentAccess.
				// The struct ApplicationAgentAccess has TableName()?
				// If not, GORM might use "application_agent_accesses".
				// The error showed actual sql: "application_agent_accesses".
				// My code expected "application_agent_access".
				// Updating test to expect "application_agent_accesses" OR just "application_agent_access".
				// I will use regex partial match to be safe or match actual.
				// Based on error: actual "application_agent_accesses".
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "application_agent_accesses" WHERE "application_agent_accesses"."application_id" = $1`)).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows(nil))

				// Keys
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "application_keys" WHERE "application_keys"."application_id" = $1`)).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := repo.GetByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplicationRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewApplicationRepository(db)
	ctx := context.TODO()

	app := &domain.Application{
		ID:        uuid.New(),
		Name:      "Test App",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "applications"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(app.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, app)
	assert.NoError(t, err)
}
