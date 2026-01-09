package domain

import (
	"time"

	"github.com/google/uuid"
)

type LLMProvider struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string     `gorm:"type:varchar(255);not null;unique" json:"name" example:"OpenAI"`
	WebsiteURL string     `gorm:"type:varchar(255)" json:"website_url" example:"https://openai.com"`
	Models     []LLMModel `gorm:"foreignKey:ProviderID" json:"models,omitempty"`
}

type LLMModel struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProviderID           uuid.UUID `gorm:"type:uuid;not null" json:"provider_id"`
	FamilyName           string    `gorm:"type:varchar(100);not null" json:"family_name" example:"GPT-4"`
	VersionName          string    `gorm:"type:varchar(100);not null" json:"version_name" example:"Turbo"`
	ApiModelName         string    `gorm:"type:varchar(255);not null" json:"api_model_name" example:"gpt-4-turbo"`
	IsLocal              bool      `gorm:"default:false" json:"is_local"`
	BaseURL              string    `gorm:"type:varchar(255)" json:"base_url,omitempty"`
	APIKeyEnvVar         string    `gorm:"type:varchar(255)" json:"api_key_env_var,omitempty" example:"OPENAI_API_KEY"`
	ContextWindowSize    int       `gorm:"type:int" json:"context_window_size" example:"128000"`
	CostPerMillionTokens float64   `gorm:"type:decimal(10,4)" json:"cost_per_million_tokens" example:"10.00"`
	IsActive             bool      `gorm:"default:true" json:"is_active"`

	Provider LLMProvider `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"provider,omitempty"`
}

type AgentLLM struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgentID     uuid.UUID `gorm:"type:uuid;not null" json:"agent_id"`
	LLMModelID  uuid.UUID `gorm:"type:uuid;not null" json:"llm_model_id"`
	IsPrimary   bool      `gorm:"default:false" json:"is_primary"`
	Temperature float64   `gorm:"type:float;default:0.7" json:"temperature"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`

	Agent    Agent    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agent,omitempty"`
	LLMModel LLMModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"llm_model,omitempty"`
}
