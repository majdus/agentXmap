package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditAction string

const (
	AuditActionCreate     AuditAction = "create"
	AuditActionUpdate     AuditAction = "update"
	AuditActionDelete     AuditAction = "delete"
	AuditActionLogin      AuditAction = "login"
	AuditActionExportData AuditAction = "export_data"
)

type Certification struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(255);not null;unique" json:"name" example:"ISO 27001"`
	IssuingAuthority string    `gorm:"type:varchar(255);not null" json:"issuing_authority" example:"ISO"`
	Description      string    `gorm:"type:text" json:"description" example:"Information Security Management"`
	BadgeURL         string    `gorm:"type:varchar(255)" json:"badge_url,omitempty" example:"https://example.com/badge.png"`
	OfficialLink     string    `gorm:"type:varchar(255)" json:"official_link,omitempty" example:"https://iso.org"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
}

type AgentCertification struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgentID         uuid.UUID  `gorm:"type:uuid;not null" json:"agent_id"`
	CertificationID uuid.UUID  `gorm:"type:uuid;not null" json:"certification_id"`
	ReferenceNumber string     `gorm:"type:varchar(100)" json:"reference_number,omitempty" example:"CERT-12345"`
	ObtainedAt      time.Time  `gorm:"default:current_date" json:"obtained_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`

	Agent         Agent         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agent,omitempty"`
	Certification Certification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"certification,omitempty"`
}

type LLMModelCertification struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LLMModelID      uuid.UUID  `gorm:"type:uuid;not null" json:"llm_model_id"`
	CertificationID uuid.UUID  `gorm:"type:uuid;not null" json:"certification_id"`
	ReferenceNumber string     `gorm:"type:varchar(100)" json:"reference_number,omitempty" example:"LLM-CERT-001"`
	ObtainedAt      time.Time  `gorm:"default:current_date" json:"obtained_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	ValidationURL   string     `gorm:"type:varchar(255)" json:"validation_url,omitempty"`

	LLMModel      LLMModel      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"llm_model,omitempty"`
	Certification Certification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"certification,omitempty"`
}

type ApplicationCertification struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ApplicationID   uuid.UUID  `gorm:"type:uuid;not null" json:"application_id"`
	CertificationID uuid.UUID  `gorm:"type:uuid;not null" json:"certification_id"`
	ReferenceNumber string     `gorm:"type:varchar(100)" json:"reference_number,omitempty"`
	ObtainedAt      time.Time  `gorm:"default:current_date" json:"obtained_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`

	Application   Application   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"application,omitempty"`
	Certification Certification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"certification,omitempty"`
}

type SystemAuditLog struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActorUserID *uuid.UUID      `gorm:"type:uuid" json:"actor_user_id,omitempty"`
	EntityType  string          `gorm:"type:varchar(50);not null" json:"entity_type" example:"agent"`
	EntityID    uuid.UUID       `gorm:"type:uuid;not null" json:"entity_id"`
	Action      AuditAction     `gorm:"type:audit_action;not null" json:"action" example:"update"`
	Changes     json.RawMessage `gorm:"type:jsonb" json:"changes,omitempty"`
	IPAddress   string          `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	OccurredAt  time.Time       `gorm:"default:now()" json:"occurred_at"`
}

// AgentExecution matches the partitioned table.
// GORM handling of partitioning requires care, often just insert/read.
type AgentExecution struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt        time.Time  `gorm:"default:now();primaryKey" json:"created_at"` // Partition Key
	AgentID          uuid.UUID  `gorm:"type:uuid;not null" json:"agent_id"`
	AgentVersionID   uuid.UUID  `gorm:"type:uuid;not null" json:"agent_version_id"`
	LLMModelID       uuid.UUID  `gorm:"type:uuid;not null" json:"llm_model_id"`
	UserID           *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ApplicationID    *uuid.UUID `gorm:"type:uuid" json:"application_id,omitempty"`
	Status           string     `gorm:"type:varchar(50)" json:"status" example:"success"`
	LatencyMs        int        `gorm:"type:int" json:"latency_ms"`
	TokenUsageInput  int        `gorm:"type:int" json:"token_usage_input"`
	TokenUsageOutput int        `gorm:"type:int" json:"token_usage_output"`
	IsPIIDetected    bool       `gorm:"default:false" json:"is_pii_detected"`
	SafetyScore      float64    `gorm:"type:float" json:"safety_score"`
}
