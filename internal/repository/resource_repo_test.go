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

func TestResourceRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewResourceRepository(db)
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
				// 1. Get Resource
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resources" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type_id"}).AddRow(id, "postgres"))

				// 2. Preloads

				// Secret
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_secrets" WHERE "resource_secrets"."resource_id" = $1`)).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows(nil))

				// Type
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resource_types" WHERE "resource_types"."id" = $1`)).
					WithArgs("postgres").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("postgres", "PostgreSQL"))
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   id,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "resources" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
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
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestResourceRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewResourceRepository(db)
	ctx := context.TODO()

	res := &domain.Resource{
		ID:        uuid.New(),
		Name:      "Test Resource",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		input   *domain.Resource
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			input: res,
			mock: func() {
				// GORM behavior: 7 args seen in error.
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resources"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // 7 args
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(res.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Error",
			input: res,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "resources"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
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
