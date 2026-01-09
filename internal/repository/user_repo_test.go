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
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)

	ctx := context.TODO()
	user := &domain.User{
		ID:             uuid.New(),
		OrganizationID: uuid.New(),
		Email:          "test@example.com",
		PasswordHash:   "hashedpassword",
		Role:           domain.UserRoleUser,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tests := []struct {
		name    string
		input   *domain.User
		mock    func()
		wantErr bool
	}{
		{
			name:  "Success",
			input: user,
			mock: func() {
				mock.ExpectBegin()
				// GORM + Postgres = Query with RETURNING
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Database Error",
			input: user,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
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

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	ctx := context.TODO()

	email := "test@example.com"
	userID := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    *domain.User
		wantErr bool
	}{
		{
			name:  "Found",
			email: email,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "organization_id", "password_hash", "role", "created_at", "updated_at"}).
					AddRow(userID, email, orgID, "hash", "user", time.Now(), time.Now())

				// 1. Query User
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs(email, 1).
					WillReturnRows(rows)

				// 2. Preload Org
				orgRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(orgID, "Test Org")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1`)).
					WithArgs(orgID).
					WillReturnRows(orgRows)
			},
			want:    &domain.User{ID: userID, Email: email},
			wantErr: false,
		},
		{
			name:  "Not Found",
			email: "unknown@example.com",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs("unknown@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetByEmail(ctx, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Email, got.Email)
			} else {
				assert.Nil(t, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)
	ctx := context.TODO()

	id := uuid.New()
	orgID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    *domain.User
		wantErr bool
	}{
		{
			name: "Found",
			id:   id,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "organization_id", "password_hash", "role", "created_at", "updated_at"}).
					AddRow(id, "test@example.com", orgID, "hash", "user", time.Now(), time.Now())

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnRows(rows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1`)).
					WithArgs(orgID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(orgID, "Org"))
			},
			want:    &domain.User{ID: id, Email: "test@example.com"},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   id,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
			} else {
				assert.Nil(t, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
