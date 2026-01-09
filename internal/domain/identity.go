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

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

// Organization represents a tenant or company using the platform.
type Organization struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name" example:"Acme Corp"`
	Slug      string         `gorm:"type:varchar(255);not null;unique" json:"slug" example:"acme-corp"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Users     []User     `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
	Agents    []Agent    `gorm:"foreignKey:OrganizationID" json:"agents,omitempty"`
	Resources []Resource `gorm:"foreignKey:OrganizationID" json:"resources,omitempty"`
}

// User represents a system user belonging to an organization.
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null" json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email          string         `gorm:"type:varchar(255);not null;unique" json:"email" example:"john.doe@acme.com"`
	PasswordHash   string         `gorm:"type:varchar(255);not null" json:"-"` // Never export password hash
	Role           UserRole       `gorm:"type:user_role;default:'user';not null" json:"role" example:"admin"`
	FirstName      string         `gorm:"type:varchar(100)" json:"first_name" example:"John"`
	LastName       string         `gorm:"type:varchar(100)" json:"last_name" example:"Doe"`
	CreatedAt      time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Organization Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
}

// Invitation represents a pending or completed user invitation.
type Invitation struct {
	ID             uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;not null" json:"organization_id"`
	InvitorID      uuid.UUID        `gorm:"type:uuid;not null" json:"invitor_id"`
	Email          string           `gorm:"type:varchar(255);not null" json:"email"`
	Token          string           `gorm:"type:varchar(255);not null;unique" json:"-"` // Token might be sensitive if exposed in lists
	Role           UserRole         `gorm:"type:user_role;default:'user';not null" json:"role"`
	Status         InvitationStatus `gorm:"type:invitation_status;default:'pending'" json:"status"`
	ExpiresAt      time.Time        `gorm:"not null" json:"expires_at"`
	CreatedAt      time.Time        `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"default:now()" json:"updated_at"`

	// Relations
	Organization Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Invitor      User         `gorm:"foreignKey:InvitorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invitor,omitempty"`
}
