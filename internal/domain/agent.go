package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentStatus string
type BillingCycle string

const (
	AgentStatusActive      AgentStatus = "active"
	AgentStatusInactive    AgentStatus = "inactive"
	AgentStatusMaintenance AgentStatus = "maintenance"
	AgentStatusDeprecated  AgentStatus = "deprecated"

	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
	BillingCycleOneTime BillingCycle = "one_time"
	BillingCycleCustom  BillingCycle = "custom"
)

// Agent represents an AI Agent.
type Agent struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID       `gorm:"type:uuid;not null"`
	Name           string          `gorm:"type:varchar(255);not null"`
	Status         AgentStatus     `gorm:"type:agent_status;default:'active'"`
	CostAmount     float64         `gorm:"type:decimal(10,2);default:0.00"`
	CostCurrency   string          `gorm:"type:varchar(3);default:'EUR'"`
	BillingCycle   BillingCycle    `gorm:"type:billing_cycle;default:'monthly'"`
	Configuration  json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	CreatedBy      *uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy      *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt      time.Time       `gorm:"default:now()"`
	UpdatedAt      time.Time       `gorm:"default:now()"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`

	// Relations
	Organization Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Versions     []AgentVersion
	Assignments  []AgentAssignment
}

// AgentVersion represents an immutable snapshot of an agent's configuration for compliance (EU AI Act).
type AgentVersion struct {
	ID                    uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AgentID               uuid.UUID       `gorm:"type:uuid;not null"`
	VersionNumber         int             `gorm:"not null"`
	ConfigurationSnapshot json.RawMessage `gorm:"type:jsonb;not null"`
	ReasonForChange       string          `gorm:"type:varchar(255)"`
	CreatedBy             *uuid.UUID      `gorm:"type:uuid"`
	CreatedAt             time.Time       `gorm:"default:now()"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// AgentAssignment links an agent to a user.
type AgentAssignment struct {
	AgentID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	AssignedAt time.Time `gorm:"default:now()"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User  User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
