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
	ID           string          `gorm:"type:varchar(50);primaryKey"`
	Name         string          `gorm:"type:varchar(100);not null"`
	ConfigSchema json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	SecretSchema json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	IsActive     bool            `gorm:"default:true"`
	CreatedAt    time.Time       `gorm:"default:now()"`
	UpdatedAt    time.Time       `gorm:"default:now()"`
}

type Resource struct {
	ID                uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID    uuid.UUID       `gorm:"type:uuid;not null"`
	TypeID            string          `gorm:"type:varchar(50);not null"`
	Name              string          `gorm:"type:varchar(255);not null"`
	ConnectionDetails json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	CreatedAt         time.Time       `gorm:"default:now()"`
	UpdatedAt         time.Time       `gorm:"default:now()"`
	DeletedAt         gorm.DeletedAt  `gorm:"index"`

	Organization Organization   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type         ResourceType   `gorm:"foreignKey:TypeID"`
	Secret       ResourceSecret `gorm:"foreignKey:ResourceID"`
}

type ResourceSecret struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ResourceID           uuid.UUID `gorm:"type:uuid;not null;unique"`
	EncryptedCredentials string    `gorm:"type:text;not null"`
	KeyVersionID         string    `gorm:"type:varchar(50)"`
	CreatedAt            time.Time `gorm:"default:now()"`
	UpdatedAt            time.Time `gorm:"default:now()"`
}

type AgentResourceAccess struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AgentID    uuid.UUID   `gorm:"type:uuid;not null"`
	ResourceID uuid.UUID   `gorm:"type:uuid;not null"`
	Permission AccessLevel `gorm:"type:access_level;default:'read_only'"`
	GrantedAt  time.Time   `gorm:"default:now()"`

	Agent    Agent    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
