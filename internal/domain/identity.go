package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enums
type UserRole string

const (
	UserRoleManager UserRole = "manager"
	UserRoleAdmin   UserRole = "admin"
	UserRoleUser    UserRole = "user"
)

// Organization represents a tenant or company using the platform.
type Organization struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Slug      string         `gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time      `gorm:"default:now()"`
	UpdatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Users     []User     `gorm:"foreignKey:OrganizationID"`
	Agents    []Agent    `gorm:"foreignKey:OrganizationID"`
	Resources []Resource `gorm:"foreignKey:OrganizationID"`
}

// User represents a system user belonging to an organization.
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null"`
	Email          string         `gorm:"type:varchar(255);not null;unique"`
	PasswordHash   string         `gorm:"type:varchar(255);not null"`
	Role           UserRole       `gorm:"type:user_role;default:'user';not null"`
	FirstName      string         `gorm:"type:varchar(100)"`
	LastName       string         `gorm:"type:varchar(100)"`
	CreatedAt      time.Time      `gorm:"default:now()"`
	UpdatedAt      time.Time      `gorm:"default:now()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relations
	Organization Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
