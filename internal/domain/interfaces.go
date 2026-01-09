package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository interface defines methods the persistence layer must implement.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// InvitationRepository defines access to Invitations.
type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
	GetByToken(ctx context.Context, token string) (*Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
}

// OrganizationRepository defines access to Organizations.
type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	GetBySlug(ctx context.Context, slug string) (*Organization, error)
}

// AgentRepository defines agent persistence.
type AgentRepository interface {
	Create(ctx context.Context, agent *Agent) error
	GetByID(ctx context.Context, id uuid.UUID) (*Agent, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]Agent, error)
	ListByStatus(ctx context.Context, orgID uuid.UUID, status AgentStatus) ([]Agent, error)
	Update(ctx context.Context, agent *Agent) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Versioning
	CreateVersion(ctx context.Context, version *AgentVersion) error

	// Resources
	GetResources(ctx context.Context, agentID uuid.UUID) ([]Resource, error)

	// Assignments
	GetAssignedUsers(ctx context.Context, agentID uuid.UUID) ([]User, error)
	GetAssignedAgents(ctx context.Context, userID uuid.UUID) ([]Agent, error)
	GetAssignedLLMs(ctx context.Context, agentID uuid.UUID) ([]AgentLLM, error)
	GetAssignedApplications(ctx context.Context, agentID uuid.UUID) ([]Application, error)
}

// LLMRepository defines interactions with LLM configurations.
type LLMRepository interface {
	ListProviders(ctx context.Context) ([]LLMProvider, error)
	GetModel(ctx context.Context, id uuid.UUID) (*LLMModel, error)
	ListModels(ctx context.Context, providerID uuid.UUID) ([]LLMModel, error)
	ListAgentsUsingModel(ctx context.Context, modelID uuid.UUID) ([]Agent, error)
}

// AuditRepository for compliance logging.
type AuditRepository interface {
	CreateLog(ctx context.Context, log *SystemAuditLog) error
	CreateExecution(ctx context.Context, exec *AgentExecution) error
}

// ApplicationRepository defines access to Applications.
type ApplicationRepository interface {
	Create(ctx context.Context, app *Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	GetAssignedAgents(ctx context.Context, appID uuid.UUID) ([]Agent, error)
	CreateKey(ctx context.Context, key *ApplicationKey) error
}

// ResourceRepository defines access to Resources.
type ResourceRepository interface {
	Create(ctx context.Context, res *Resource) error
	GetByID(ctx context.Context, id uuid.UUID) (*Resource, error)
	ListAgentsWithAccess(ctx context.Context, resourceID uuid.UUID) ([]Agent, error)
}
