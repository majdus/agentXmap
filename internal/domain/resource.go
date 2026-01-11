package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessLevel string

const (
	AccessLevelReadOnly  AccessLevel = "read_only"
	AccessLevelReadWrite AccessLevel = "read_write"
)

type ResourceType struct {
	ID           string          `gorm:"type:varchar(50);primaryKey" json:"id" example:"postgres-db"`
	Name         string          `gorm:"type:varchar(100);not null" json:"name" example:"PostgreSQL Database"`
	ConfigSchema json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"config_schema" swaggertype:"string"`
	SecretSchema json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"secret_schema" swaggertype:"string"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time       `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"default:now()" json:"updated_at"`
}

type Resource struct {
	ID                uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TypeID            string          `gorm:"type:varchar(50);not null" json:"type_id"`
	Name              string          `gorm:"type:varchar(255);not null" json:"name" example:"Production DB"`
	ConnectionDetails json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"connection_details" swaggertype:"string"`
	CreatedAt         time.Time       `gorm:"default:now()" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"default:now()" json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`

	Type   ResourceType   `gorm:"foreignKey:TypeID" json:"type,omitempty"`
	Secret ResourceSecret `gorm:"foreignKey:ResourceID" json:"-"` // Never expose secret relation directly
}

type ResourceSecret struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResourceID           uuid.UUID `gorm:"type:uuid;not null;unique" json:"resource_id"`
	EncryptedCredentials string    `gorm:"type:text;not null" json:"-"` // Explicitly hide
	KeyVersionID         string    `gorm:"type:varchar(50)" json:"key_version_id"`
	CreatedAt            time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt            time.Time `gorm:"default:now()" json:"updated_at"`
}

type AgentResourceAccess struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgentID    uuid.UUID   `gorm:"type:uuid;not null" json:"agent_id"`
	ResourceID uuid.UUID   `gorm:"type:uuid;not null" json:"resource_id"`
	Permission AccessLevel `gorm:"type:access_level;default:'read_only'" json:"permission"`
	GrantedAt  time.Time   `gorm:"default:now()" json:"granted_at"`

	Agent    Agent    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agent,omitempty"`
	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"resource,omitempty"`
}
