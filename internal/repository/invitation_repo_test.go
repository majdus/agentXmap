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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestInvitationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewInvitationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		invitation := &domain.Invitation{
			InvitorID: uuid.New(),
			Email:     "test@example.com",
			Token:     "secure-token",
			Role:      domain.UserRoleUser,
			Status:    domain.InvitationStatusPending,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "invitations"`)).
			WithArgs(
				invitation.InvitorID,
				invitation.Email,
				invitation.Token,
				invitation.Role,
				invitation.Status,
				invitation.ExpiresAt,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), invitation)
		assert.NoError(t, err)
	})
}

func TestInvitationRepository_GetByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewInvitationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		token := "secure-token"
		invitationID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "invitor_id", "email", "token", "role", "status", "expires_at", "created_at", "updated_at",
		}).AddRow(
			invitationID, uuid.New(), "test@example.com", token, "user", "pending", time.Now(), time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invitations" WHERE token = $1`)).
			WithArgs(token, 1). // Limit 1
			WillReturnRows(rows)

		// Observed behavior: GORM queries Users (Invitor).
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		inv, err := repo.GetByToken(context.Background(), token)
		assert.NoError(t, err)
		if inv != nil {
			assert.Equal(t, token, inv.Token)
		}
	})

	t.Run("not found", func(t *testing.T) {
		token := "invalid-token"
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invitations" WHERE token = $1`)).
			WithArgs(token, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		inv, err := repo.GetByToken(context.Background(), token)
		assert.Error(t, err)
		assert.Nil(t, inv)
	})
}

func TestInvitationRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewInvitationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		invitation := &domain.Invitation{
			ID:     uuid.New(),
			Token:  "secure-token",
			Status: domain.InvitationStatusAccepted,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations"`)).
			WithArgs(
				sqlmock.AnyArg(), // InvitorID
				sqlmock.AnyArg(), // Email
				invitation.Token,
				sqlmock.AnyArg(), // Role
				invitation.Status,
				sqlmock.AnyArg(), // ExpiresAt
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				invitation.ID,    // ID (WHERE clause)
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), invitation)
		assert.NoError(t, err)
	})
}
