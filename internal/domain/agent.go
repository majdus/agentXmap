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
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string          `gorm:"type:varchar(255);not null;unique" json:"name"`
	Status        AgentStatus     `gorm:"type:agent_status;default:'active'" json:"status"`
	CostAmount    float64         `gorm:"type:decimal(10,2);default:0.00" json:"cost_amount"`
	CostCurrency  string          `gorm:"type:varchar(3);default:'EUR'" json:"cost_currency"`
	BillingCycle  BillingCycle    `gorm:"type:billing_cycle;default:'monthly'" json:"billing_cycle"`
	Configuration json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"configuration" swaggertype:"string" example:"{\"model\": \"gpt-4\"}"`
	CreatedBy     *uuid.UUID      `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy     *uuid.UUID      `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt     time.Time       `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relations

	Versions    []AgentVersion    `json:"versions,omitempty"`
	Assignments []AgentAssignment `json:"assignments,omitempty"`
}

// AgentVersion represents an immutable snapshot of an agent's configuration for compliance (EU AI Act).
type AgentVersion struct {
	ID                    uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgentID               uuid.UUID       `gorm:"type:uuid;not null" json:"agent_id"`
	VersionNumber         int             `gorm:"not null" json:"version_number"`
	ConfigurationSnapshot json.RawMessage `gorm:"type:jsonb;not null" json:"configuration_snapshot" swaggertype:"string"`
	ReasonForChange       string          `gorm:"type:varchar(255)" json:"reason_for_change"`
	CreatedBy             *uuid.UUID      `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt             time.Time       `gorm:"default:now()" json:"created_at"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// AgentAssignment links an agent to a user.
type AgentAssignment struct {
	AgentID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"agent_id"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	AssignedAt time.Time `gorm:"default:now()" json:"assigned_at"`

	Agent Agent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"agent,omitempty"`
	User  User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}
