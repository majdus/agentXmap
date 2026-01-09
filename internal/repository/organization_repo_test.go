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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestOrganizationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewOrganizationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		org := &domain.Organization{
			Name: "Test Org",
			Slug: "test-org",
		}

		mock.ExpectBegin()
		// GORM might insert CreatedAt, UpdatedAt as well if they are zero, but here it seems it only inserted Name, Slug, DeletedAt.
		// Adjusting regex to be more flexible and args to match actual observation or use AnyArg efficiently.
		// Observed: INSERT INTO "organizations" ("name","slug","deleted_at") VALUES ($1,$2,$3)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organizations"`)).
			WithArgs(org.Name, org.Slug, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), org)
		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		org := &domain.Organization{
			Name: "Test Org",
			Slug: "test-org",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "organizations"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), org)
		assert.Error(t, err)
	})
}

func TestOrganizationRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewOrganizationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "name", "slug", "created_at", "updated_at"}).
			AddRow(id, "Test Org", "test-org", time.Now(), time.Now())

		// GORM adds deleted_at check and limit
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE id = $1 AND "organizations"."deleted_at" IS NULL`)).
			WithArgs(id, 1). // Added limit arg
			WillReturnRows(rows)

		org, err := repo.GetByID(context.Background(), id)
		assert.NoError(t, err)
		if org != nil {
			assert.Equal(t, id, org.ID)
			assert.Equal(t, "Test Org", org.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE id = $1`)).
			WithArgs(id, 1). // Added limit arg
			WillReturnError(gorm.ErrRecordNotFound)

		org, err := repo.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, org)
	})
}

func TestOrganizationRepository_GetBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewOrganizationRepository(gormDB)

	t.Run("success", func(t *testing.T) {
		slug := "test-org"
		rows := sqlmock.NewRows([]string{"id", "name", "slug", "created_at", "updated_at"}).
			AddRow(uuid.New(), "Test Org", slug, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE slug = $1 AND "organizations"."deleted_at" IS NULL`)).
			WithArgs(slug, 1). // Added limit arg
			WillReturnRows(rows)

		org, err := repo.GetBySlug(context.Background(), slug)
		assert.NoError(t, err)
		if org != nil {
			assert.Equal(t, slug, org.Slug)
		}
	})

	t.Run("not found", func(t *testing.T) {
		slug := "non-existent"
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE slug = $1`)).
			WithArgs(slug, 1). // Added limit arg
			WillReturnError(gorm.ErrRecordNotFound)

		org, err := repo.GetBySlug(context.Background(), slug)
		assert.Error(t, err)
		assert.Nil(t, org)
	})
}
