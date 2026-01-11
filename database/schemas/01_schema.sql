/* * ============================================================
 * SCRIPT: 01_schema.sql
 * PURPOSE: Create the Master B2B Platform Schema
 * AUTHOR: Tech Lead & Dev
 * ============================================================
 */

BEGIN;

-- Enable UUID extension if not already available (standard in modern PG)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- 0. UTILITY FUNCTIONS
-- ============================================================

-- Function to auto-update 'updated_at' columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ============================================================
-- 1. ENUMS
-- ============================================================

CREATE TYPE user_role AS ENUM ('manager', 'admin', 'user');
CREATE TYPE agent_status AS ENUM ('active', 'inactive', 'maintenance', 'deprecated');
CREATE TYPE billing_cycle AS ENUM ('monthly', 'yearly', 'one_time', 'custom');
CREATE TYPE access_level AS ENUM ('read_only', 'read_write');
CREATE TYPE audit_action AS ENUM ('create', 'update', 'delete', 'login', 'export_data');
CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'revoked');

-- ============================================================
-- 2. CORE: IDENTITY
-- ============================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Trigger for users.updated_at
CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    role user_role NOT NULL DEFAULT 'user',
    status invitation_status DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TRIGGER update_invitations_modtime BEFORE UPDATE ON invitations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- 3. AGENT DOMAIN
-- ============================================================

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status agent_status DEFAULT 'active',

    -- Financials
    cost_amount DECIMAL(10,2) DEFAULT 0.00,
    cost_currency VARCHAR(3) DEFAULT 'EUR',
    billing_cycle billing_cycle DEFAULT 'monthly',

    -- Mutable Configuration
    configuration JSONB DEFAULT '{}',

    -- Audit
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE TRIGGER update_agents_modtime BEFORE UPDATE ON agents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE agent_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    configuration_snapshot JSONB NOT NULL,
    reason_for_change VARCHAR(255),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(agent_id, version_number)
);
COMMENT ON TABLE agent_versions IS 'Immutable snapshot for AI Governance (EU AI Act)';

CREATE TABLE agent_assignments (
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (agent_id, user_id)
);

-- ============================================================
-- 4. LLM INFRASTRUCTURE
-- ============================================================

CREATE TABLE llm_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    website_url VARCHAR(255)
);

CREATE TABLE llm_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES llm_providers(id) ON DELETE CASCADE,

    family_name VARCHAR(100) NOT NULL,
    version_name VARCHAR(100) NOT NULL,
    api_model_name VARCHAR(255) NOT NULL,

    is_local BOOLEAN DEFAULT FALSE,
    base_url VARCHAR(255),
    api_key_env_var VARCHAR(255),

    context_window_size INT,
    cost_per_million_tokens DECIMAL(10, 4),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE agent_llms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    llm_model_id UUID NOT NULL REFERENCES llm_models(id) ON DELETE CASCADE,
    is_primary BOOLEAN DEFAULT FALSE,
    temperature FLOAT DEFAULT 0.7,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, llm_model_id)
);

-- ============================================================
-- 5. SOFTWARE APPLICATIONS
-- ============================================================

CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TRIGGER update_apps_modtime BEFORE UPDATE ON applications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE application_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(8) NOT NULL,
    name VARCHAR(100),
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE application_agent_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    can_invoke BOOLEAN DEFAULT TRUE,
    rate_limit INT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(application_id, agent_id)
);

-- ============================================================
-- 6. RESOURCES & SECRETS
-- ============================================================

CREATE TABLE resource_types (
    id VARCHAR(50) PRIMARY KEY, -- e.g. 'postgres'
    name VARCHAR(100) NOT NULL,
    config_schema JSONB DEFAULT '{}',
    secret_schema JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_id VARCHAR(50) NOT NULL REFERENCES resource_types(id) ON UPDATE CASCADE,
    name VARCHAR(255) NOT NULL,
    connection_details JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE TRIGGER update_resources_modtime BEFORE UPDATE ON resources FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE resource_secrets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL UNIQUE REFERENCES resources(id) ON DELETE CASCADE,
    encrypted_credentials TEXT NOT NULL,
    key_version_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE agent_resource_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    permission access_level DEFAULT 'read_only',
    granted_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, resource_id)
);

-- ============================================================
-- 7. TRUST & CERTIFICATIONS
-- ============================================================

CREATE TABLE certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    issuing_authority VARCHAR(255) NOT NULL,
    description TEXT,
    badge_url VARCHAR(255),
    official_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Linking Tables
CREATE TABLE llm_model_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    llm_model_id UUID NOT NULL REFERENCES llm_models(id) ON DELETE CASCADE,
    certification_id UUID NOT NULL REFERENCES certifications(id) ON DELETE CASCADE,
    reference_number VARCHAR(100),
    obtained_at DATE DEFAULT CURRENT_DATE,
    expires_at DATE,
    validation_url VARCHAR(255),
    UNIQUE(llm_model_id, certification_id)
);

CREATE TABLE agent_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    certification_id UUID NOT NULL REFERENCES certifications(id) ON DELETE CASCADE,
    reference_number VARCHAR(100),
    obtained_at DATE DEFAULT CURRENT_DATE,
    expires_at DATE,
    UNIQUE(agent_id, certification_id)
);

CREATE TABLE application_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    certification_id UUID NOT NULL REFERENCES certifications(id) ON DELETE CASCADE,
    reference_number VARCHAR(100),
    obtained_at DATE DEFAULT CURRENT_DATE,
    expires_at DATE,
    UNIQUE(application_id, certification_id)
);

-- ============================================================
-- 8. COMPLIANCE & BIG DATA (PARTITIONING)
-- ============================================================

CREATE TABLE system_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID, -- Nullable if system action
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action audit_action NOT NULL,
    changes_json JSONB,
    ip_address VARCHAR(45),
    occurred_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_audit_date ON system_audit_logs(occurred_at);
CREATE INDEX idx_audit_entity ON system_audit_logs(entity_id);

-- Partitioned Table for Executions
-- NOTE: In partitioned tables, the PK MUST include the partition key.
CREATE TABLE agent_executions (
    id UUID DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(), -- Partition Key

    agent_id UUID NOT NULL, -- Loose reference to keep data even if agent deleted? No, let's keep it safe.
    agent_version_id UUID NOT NULL,
    llm_model_id UUID NOT NULL,

    user_id UUID,
    application_id UUID,

    status VARCHAR(50),
    latency_ms INT,
    token_usage_input INT,
    token_usage_output INT,

    is_pii_detected BOOLEAN DEFAULT FALSE,
    safety_score FLOAT,

    PRIMARY KEY (id, created_at) -- Composite PK required for partitioning
) PARTITION BY RANGE (created_at);

-- Create Initial Partition (e.g., for the current year or month)
-- Usually managed by code/scripts, but here is a starter partition:
CREATE TABLE agent_executions_default PARTITION OF agent_executions
    DEFAULT;
    -- 'DEFAULT' partition catches everything not covered by specific ranges.
    -- In prod, use specific ranges (e.g., FROM '2024-01-01' TO '2024-02-01')

-- Indexes on the parent table (propagated to partitions)
CREATE INDEX idx_executions_agent ON agent_executions(agent_id);
CREATE INDEX idx_executions_date ON agent_executions(created_at);


COMMIT;

-- Final Check
SELECT 'DATABASE SCHEMA CREATED SUCCESSFULLY' as status;
