/* * ============================================================
 * SCRIPT: 00_reset.sql
 * PURPOSE: Drop everything to start fresh (DEV ONLY)
 * WARNING: DESTRUCTIVE ACTION
 * ============================================================
 */

BEGIN;

-- 1. Drop Tables (Cascade handles FKs automatically)
DROP TABLE IF EXISTS application_certifications CASCADE;
DROP TABLE IF EXISTS agent_certifications CASCADE;
DROP TABLE IF EXISTS llm_model_certifications CASCADE;
DROP TABLE IF EXISTS certifications CASCADE;

DROP TABLE IF EXISTS agent_executions CASCADE; -- Will drop partitions too
DROP TABLE IF EXISTS system_audit_logs CASCADE;

DROP TABLE IF EXISTS agent_resource_access CASCADE;
DROP TABLE IF EXISTS resource_secrets CASCADE;
DROP TABLE IF EXISTS resources CASCADE;
DROP TABLE IF EXISTS resource_types CASCADE;

DROP TABLE IF EXISTS application_agent_access CASCADE;
DROP TABLE IF EXISTS application_keys CASCADE;
DROP TABLE IF EXISTS applications CASCADE;

DROP TABLE IF EXISTS agent_llms CASCADE;
DROP TABLE IF EXISTS llm_models CASCADE;
DROP TABLE IF EXISTS llm_providers CASCADE;

DROP TABLE IF EXISTS agent_assignments CASCADE;
DROP TABLE IF EXISTS agent_versions CASCADE;
DROP TABLE IF EXISTS agents CASCADE;

DROP TABLE IF EXISTS users CASCADE;

-- 2. Drop Enums (Types)
DROP TYPE IF EXISTS audit_action CASCADE;
DROP TYPE IF EXISTS access_level CASCADE;
DROP TYPE IF EXISTS billing_cycle CASCADE;
DROP TYPE IF EXISTS agent_status CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;

-- 3. Drop Functions
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;

COMMIT;

-- Verification
SELECT 'DATABASE CLEARED SUCCESSFULLY' as status;
