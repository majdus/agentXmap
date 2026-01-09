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
	"gorm.io/gorm"
)

func TestAgentRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAgentRepository(db)
	ctx := context.TODO()

	agent := &domain.Agent{
		ID:             uuid.New(),
		OrganizationID: uuid.New(),
		Name:           "Test Agent",
		Status:         domain.AgentStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tests := []struct {
		name    string
		input   *domain.Agent
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			input: agent,
			mock: func() {
				mock.ExpectBegin()
				// GORM inserted 12 args in the failed log.
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agents"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(agent.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Fail",
			input: agent,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "agents"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Create(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAgentRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewAgentRepository(db)
	ctx := context.TODO()

	id := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    *domain.Agent
		wantErr bool
	}{
		{
			name: "Found",
			id:   id,
			mock: func() {
				// 1. Main Query
				rows := sqlmock.NewRows([]string{"id", "name", "organization_id", "status", "cost_amount", "cost_currency", "billing_cycle", "configuration"}).
					AddRow(id, "Agent 007", orgID, "active", 0.0, "EUR", "monthly", []byte("{}"))

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "agents" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnRows(rows)

				// 2. Preloads
				// Organization is usually fetched if FK is present.
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1`)).
					WithArgs(orgID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(orgID, "Org"))

				// Versions
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "agent_versions" WHERE "agent_versions"."agent_id" = $1`)).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			want:    &domain.Agent{ID: id, Name: "Agent 007"},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   id,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "agents" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				if assert.NotNil(t, got) {
					assert.Equal(t, tt.want.ID, got.ID)
				}
			} else {
				assert.Nil(t, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
