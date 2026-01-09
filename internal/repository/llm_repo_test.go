package repository

import (
	"context"
	"regexp"
	"testing"

	"agentXmap/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLLMRepository_ListProviders(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewLLMRepository(db)
	ctx := context.TODO()

	providerID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				// 1. SELECT * FROM providers
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(providerID, "OpenAI")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "llm_providers"`)).
					WillReturnRows(rows)

				// 2. Preload Models
				// Error was: "failed to assign association ... make sure foreign fields exists".
				// LLMModel struct has ProviderID.
				// GORM needs the FK column in the result set of the query?
				// Result set "id", "family_name". Missing "provider_id".
				// GORM Association needs FK to map back.
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "llm_models" WHERE "llm_models"."provider_id" = $1`)).
					WithArgs(providerID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "family_name", "provider_id"}).AddRow(uuid.New(), "GPT-4", providerID))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := repo.ListProviders(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, 1) // Expect 1 provider
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLLMRepository_GetModel(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewLLMRepository(db)
	ctx := context.TODO()

	modelID := uuid.New()
	providerID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    *domain.LLMModel
		wantErr bool
	}{
		{
			name: "Found",
			id:   modelID,
			mock: func() {
				// 1. SELECT * FROM models
				rows := sqlmock.NewRows([]string{"id", "provider_id", "family_name"}).
					AddRow(modelID, providerID, "GPT-4")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "llm_models" WHERE id = $1`)).
					WithArgs(modelID, 1).
					WillReturnRows(rows)

				// 2. Preload Provider
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "llm_providers" WHERE "llm_providers"."id" = $1`)).
					WithArgs(providerID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(providerID, "OpenAI"))
			},
			want:    &domain.LLMModel{ID: modelID, FamilyName: "GPT-4"},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   modelID,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "llm_models" WHERE id = $1`)).
					WithArgs(modelID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := repo.GetModel(ctx, tt.id)
			if tt.name == "Not Found" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
