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
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name             string    `gorm:"type:varchar(255);not null;unique"`
	IssuingAuthority string    `gorm:"type:varchar(255);not null"`
	Description      string    `gorm:"type:text"`
	BadgeURL         string    `gorm:"type:varchar(255)"`
	OfficialLink     string    `gorm:"type:varchar(255)"`
	CreatedAt        time.Time `gorm:"default:now()"`
}

type AgentCertification struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AgentID         uuid.UUID `gorm:"type:uuid;not null"`
	CertificationID uuid.UUID `gorm:"type:uuid;not null"`
	ReferenceNumber string    `gorm:"type:varchar(100)"`
	ObtainedAt      time.Time `gorm:"default:current_date"`
	ExpiresAt       *time.Time

	Agent         Agent         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Certification Certification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type SystemAuditLog struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID       `gorm:"type:uuid;not null"`
	ActorUserID    *uuid.UUID      `gorm:"type:uuid"`
	EntityType     string          `gorm:"type:varchar(50);not null"`
	EntityID       uuid.UUID       `gorm:"type:uuid;not null"`
	Action         AuditAction     `gorm:"type:audit_action;not null"`
	Changes        json.RawMessage `gorm:"type:jsonb"`
	IPAddress      string          `gorm:"type:varchar(45)"`
	OccurredAt     time.Time       `gorm:"default:now()"`
}

// AgentExecution matches the partitioned table.
// GORM handling of partitioning requires care, often just insert/read.
type AgentExecution struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt        time.Time  `gorm:"default:now();primaryKey"` // Partition Key
	OrganizationID   uuid.UUID  `gorm:"type:uuid;not null"`
	AgentID          uuid.UUID  `gorm:"type:uuid;not null"`
	AgentVersionID   uuid.UUID  `gorm:"type:uuid;not null"`
	LLMModelID       uuid.UUID  `gorm:"type:uuid;not null"`
	UserID           *uuid.UUID `gorm:"type:uuid"`
	ApplicationID    *uuid.UUID `gorm:"type:uuid"`
	Status           string     `gorm:"type:varchar(50)"`
	LatencyMs        int        `gorm:"type:int"`
	TokenUsageInput  int        `gorm:"type:int"`
	TokenUsageOutput int        `gorm:"type:int"`
	IsPIIDetected    bool       `gorm:"default:false"`
	SafetyScore      float64    `gorm:"type:float"`
}
