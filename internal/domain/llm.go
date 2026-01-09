package domain

import (
	"time"

	"github.com/google/uuid"
)

type LLMProvider struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string     `gorm:"type:varchar(255);not null;unique"`
	WebsiteURL string     `gorm:"type:varchar(255)"`
	Models     []LLMModel `gorm:"foreignKey:ProviderID"`
}

type LLMModel struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProviderID           uuid.UUID `gorm:"type:uuid;not null"`
	FamilyName           string    `gorm:"type:varchar(100);not null"`
	VersionName          string    `gorm:"type:varchar(100);not null"`
	ApiModelName         string    `gorm:"type:varchar(255);not null"`
	IsLocal              bool      `gorm:"default:false"`
	BaseURL              string    `gorm:"type:varchar(255)"`
	APIKeyEnvVar         string    `gorm:"type:varchar(255)"`
	ContextWindowSize    int       `gorm:"type:int"`
	CostPerMillionTokens float64   `gorm:"type:decimal(10,4)"`
	IsActive             bool      `gorm:"default:true"`

	Provider LLMProvider `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AgentLLM struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AgentID     uuid.UUID `gorm:"type:uuid;not null"`
	LLMModelID  uuid.UUID `gorm:"type:uuid;not null"`
	IsPrimary   bool      `gorm:"default:false"`
	Temperature float64   `gorm:"type:float;default:0.7"`
	CreatedAt   time.Time `gorm:"default:now()"`

	Agent    Agent    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LLMModel LLMModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
