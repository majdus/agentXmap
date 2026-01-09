# Services & Interfaces Report

This document provides a comprehensive overview of the backend services within the **AgentXmap** platform. It details the responsibility of each service and lists the interfaces available for interaction.

---

## 1. Identity Service

**Responsibility**: Manages user authentication, organization creation, and user invitations. It handles the secure onboarding of new tenants and users.

### Interfaces

- **`SignUp(ctx, orgName, email, password)`**
  - Creates a new Organization and the initial Admin User.
  - Returns: `*User`, `error`
- **`Login(ctx, email, password)`**
  - Authenticates a user using email and password. Generates a secure session/token context (logic implied).
  - Returns: `*User`, `error`
- **`InviteUsers(ctx, invitorID, emails, role)`**
  - Sends email invitations to new users to join an existing organization. Requires Admin/Manager permissions.
  - Returns: `[]*Invitation`, `error`
- **`AcceptInvitation(ctx, token, password, firstName, lastName)`**
  - Completes the user registration process using a valid invitation token.
  - Returns: `*User`, `error`

---

## 2. Agent Service

**Responsibility**: The core service for managing AI Agents. It handles lifecycle (CRUD), configuration versioning, resource assignments, and billing calculations.

### Interfaces

- **`CreateAgent(ctx, orgID, userID, name, config)`**
  - Creates a new Agent and initializes its first configuration version.
  - Returns: `*domain.Agent`, `error`
- **`GetAgent(ctx, id)`**
  - Retrieves detailed information about a specific Agent.
  - Returns: `*domain.Agent`, `error`
- **`ListAgents(ctx, orgID)`**
  - Lists all agents belonging to a specific Organization.
  - Returns: `[]domain.Agent`, `error`
- **`ListAgentsByStatus(ctx, orgID, status)`**
  - Filters agents by their status (e.g., Active, Inactive).
  - Returns: `[]domain.Agent`, `error`
- **`UpdateAgent(ctx, id, userID, name, config, status)`**
  - Updates an Agent's details. Automatically creates a new `AgentVersion` if the configuration changes.
  - Returns: `*domain.Agent`, `error`
- **`DeleteAgent(ctx, id)`**
  - Soft-deletes an Agent.
  - Returns: `error`
- **`GetActiveMonthlyCost(ctx, orgID)`**
  - Calculates the total projected monthly cost for all active agents in an organization.
  - Returns: `float64`, `error`
- **`ListAgentResources(ctx, agentID)`**
  - Lists external resources (DBs, APIs) assigned to an Agent.
  - Returns: `[]domain.Resource`, `error`
- **`ListAssignedUsers(ctx, agentID)`**
  - Lists users who have permission to manage or use this Agent.
  - Returns: `[]domain.User`, `error`
- **`ListAssignedAgents(ctx, userID)`**
  - Lists agents assigned to a specific user.
  - Returns: `[]domain.Agent`, `error`
- **`GetAgentLLMs(ctx, agentID)`**
  - Retrieves the LLM Models (e.g., GPT-4) configured for this Agent.
  - Returns: `[]domain.AgentLLM`, `error`
- **`ListAssignedApplications(ctx, agentID)`**
  - Lists external Applications that are authorized to invoke this Agent.
  - Returns: `[]domain.Application`, `error`
- **`ListAgentCertifications(ctx, agentID)`**
  - Lists compliance certifications (e.g., ISO 27001) associated with this Agent.
  - Returns: `[]domain.Certification`, `error`

---

## 3. Application Service

**Responsibility**: Manages external Applications (API Consumers) that integrate with the platform. Handles API Key generation and access control.

### Interfaces

- **`CreateApplication(ctx, ownerID, name, description)`**
  - Registers a new external Application.
  - Returns: `*domain.Application`, `error`
- **`GetApplication(ctx, id)`**
  - Retrieves details of an Application including its keys.
  - Returns: `*domain.Application`, `error`
- **`CreateAPIKey(ctx, appID, name)`**
  - Generates a new secure API Key (`sk-live-...`) for an Application. The raw key is returned only once.
  - Returns: `rawKey string`, `*domain.ApplicationKey`, `error`
- **`ListAssignedAgents(ctx, appID)`**
  - Lists Agents that this Application is authorized to access.
  - Returns: `[]domain.Agent`, `error`
- **`ListApplicationCertifications(ctx, appID)`**
  - Lists compliance certifications associated with this Application.
  - Returns: `[]domain.Certification`, `error`

---

## 4. LLM Service

**Responsibility**: Manages the catalog of Large Language Models (LLMs) and Providers. Used to configure what models are available to Agents.

### Interfaces

- **`ListProviders(ctx)`**
  - Lists available LLM Providers (e.g., OpenAI, Anthropic).
  - Returns: `[]domain.LLMProvider`, `error`
- **`ListModels(ctx, providerID)`**
  - Lists models available under a specific Provider.
  - Returns: `[]domain.LLMModel`, `error`
- **`GetModel(ctx, id)`**
  - Retrieves details of a specific LLM Model.
  - Returns: `*domain.LLMModel`, `error`
- **`ListAgentsUsingModel(ctx, modelID)`**
  - Finds all Agents currently configured to use a specific LLM Model. Useful for impact analysis (e.g., deprecating a model).
  - Returns: `[]domain.Agent`, `error`
- **`ListModelCertifications(ctx, modelID)`**
  - Lists compliance certifications associated with a specific LLM Model.
  - Returns: `[]domain.Certification`, `error`

---

## 5. Resource Service

**Responsibility**: Manages external resources (Databases, Third-party APIs) that Agents interact with.

### Interfaces

- **`CreateResource(ctx, orgID, typeID, name, config)`**
  - Registers a new Resource (e.g., a Postgres DB connection).
  - Returns: `*domain.Resource`, `error`
- **`GetResource(ctx, id)`**
  - Retrieves details of a Resource.
  - Returns: `*domain.Resource`, `error`
- **`ListAgentsWithAccess(ctx, resourceID)`**
  - Lists all Agents that have been granted access to this Resource.
  - Returns: `[]domain.Agent`, `error`

---

## 6. Audit Service

**Responsibility**: Handles immutable logging for compliance and security. Tracks system actions and agent executions.

### Interfaces

- **`LogAction(ctx, orgID, actorUserID, entityType, entityID, action, changes, ipAddress)`**
  - Records a system event (e.g., User X updated Agent Y).
  - Returns: `error`
- **`RecordExecution(ctx, exec)`**
  - Logs the execution of an Agent, including latency, token usage, and safety scores.
  - Returns: `error`
