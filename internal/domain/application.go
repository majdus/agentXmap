package domain

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(255);not null;unique"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`

	Owner       User                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Keys        []ApplicationKey         `gorm:"foreignKey:ApplicationID"`
	AgentAccess []ApplicationAgentAccess `gorm:"foreignKey:ApplicationID"`
}

type ApplicationKey struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null"`
	KeyHash       string    `gorm:"type:varchar(255);not null"`
	KeyPrefix     string    `gorm:"type:varchar(8);not null"`
	Name          string    `gorm:"type:varchar(100)"`
	LastUsedAt    *time.Time
	ExpiresAt     *time.Time
	CreatedAt     time.Time `gorm:"default:now()"`
}

type ApplicationAgentAccess struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null"`
	AgentID       uuid.UUID `gorm:"type:uuid;not null"`
	CanInvoke     bool      `gorm:"default:true"`
	RateLimit     *int
	CreatedAt     time.Time `gorm:"default:now()"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
