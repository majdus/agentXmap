package domain

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name" example:"My App"`
	Description string    `gorm:"type:text" json:"description" example:"An integration app"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`

	Owner       User                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner,omitempty"`
	Keys        []ApplicationKey         `gorm:"foreignKey:ApplicationID" json:"keys,omitempty"`
	AgentAccess []ApplicationAgentAccess `gorm:"foreignKey:ApplicationID" json:"agent_access,omitempty"`
}

type ApplicationKey struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ApplicationID uuid.UUID  `gorm:"type:uuid;not null" json:"application_id"`
	KeyHash       string     `gorm:"type:varchar(255);not null" json:"-"` // Never expose hash
	KeyPrefix     string     `gorm:"type:varchar(8);not null" json:"key_prefix" example:"sk-live-"`
	Name          string     `gorm:"type:varchar(100)" json:"name" example:"Production Key"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
}

type ApplicationAgentAccess struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null" json:"application_id"`
	AgentID       uuid.UUID `gorm:"type:uuid;not null" json:"agent_id"`
	CanInvoke     bool      `gorm:"default:true" json:"can_invoke"`
	RateLimit     *int      `json:"rate_limit,omitempty"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agent,omitempty"`
}
