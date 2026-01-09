/* * ============================================================
 * SCRIPT: 02_seed.sql
 * PURPOSE: Populate DB with "State of the Art" Data (Jan 2026)
 * SOURCE: Strategic Report - Jan 2026
 * AUTHOR: Tech Lead
 * ============================================================
 */

BEGIN;

-- ============================================================
-- 1. SEED: LLM PROVIDERS
-- ============================================================
INSERT INTO llm_providers (name, website_url) VALUES
    ('OpenAI', 'https://openai.com'),
    ('Google DeepMind', 'https://deepmind.google'),
    ('Anthropic', 'https://www.anthropic.com'),
    ('Meta AI', 'https://ai.meta.com'),
    ('Mistral AI', 'https://mistral.ai'),
    ('DeepSeek', 'https://www.deepseek.com'),
    ('Groq', 'https://groq.com'), -- Added for SaaS inference of Open Weights
    ('Ollama (Local)', 'https://ollama.com')
ON CONFLICT (name) DO NOTHING;

-- ============================================================
-- 2. SEED: LLM MODELS (The 2026 Landscape)
-- ============================================================

-- --- OpenAI (The Reasoning Era) ---
-- Source: Report Section 2.1 [cite: 43]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, context_window_size, cost_per_million_tokens) VALUES
    -- GPT-5.2: The "System 2" Reasoning Model
    ((SELECT id FROM llm_providers WHERE name = 'OpenAI'), 'GPT-5', '5.2 (Reasoning)', 'gpt-5.2-2025-12', false, 400000, 1.75),
    -- GPT-5-mini: High Throughput / Cost Killer
    ((SELECT id FROM llm_providers WHERE name = 'OpenAI'), 'GPT-5', 'Mini', 'gpt-5-mini', false, 400000, 0.25);

-- --- Google DeepMind (Multimodal Native) ---
-- Source: Report Section 2.2 [cite: 66]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, context_window_size, cost_per_million_tokens) VALUES
    -- Gemini 3 Pro: Video/Audio Native
    ((SELECT id FROM llm_providers WHERE name = 'Google DeepMind'), 'Gemini 3', 'Pro', 'gemini-3-pro', false, 1000000, 2.00),
    -- Gemini 2.5 Flash: Massive Context Cheap
    ((SELECT id FROM llm_providers WHERE name = 'Google DeepMind'), 'Gemini 2', '2.5 Flash', 'gemini-2.5-flash', false, 1000000, 0.10);

-- --- Anthropic (Computer Use) ---
-- Source: Report Section 2.3 [cite: 90]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, context_window_size, cost_per_million_tokens) VALUES
    -- Claude 4.5 Opus: The Premium Writer
    ((SELECT id FROM llm_providers WHERE name = 'Anthropic'), 'Claude 4.5', 'Opus', 'claude-4.5-opus', false, 200000, 5.00),
    -- Claude 4.5 Sonnet: The Coding Workhorse
    ((SELECT id FROM llm_providers WHERE name = 'Anthropic'), 'Claude 4.5', 'Sonnet', 'claude-4.5-sonnet', false, 200000, 3.00);

-- --- Meta Llama 4 (Mixture of Experts) ---
-- Source: Report Section 3.1 [cite: 136]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, base_url, context_window_size, cost_per_million_tokens) VALUES
    -- Llama 4 Scout (Hosted on Groq for Speed)
    ((SELECT id FROM llm_providers WHERE name = 'Groq'), 'Llama 4', 'Scout (SaaS)', 'llama4-scout-16x', false, 'https://api.groq.com/openai/v1', 10000000, 0.11),
    -- Llama 4 Maverick (Self-Hosted for Sovereignty)
    ((SELECT id FROM llm_providers WHERE name = 'Meta AI'), 'Llama 4', 'Maverick (Local)', 'llama4:maverick', true, 'http://localhost:11434', 1000000, 0.00);

-- --- Mistral AI (European Sovereignty) ---
-- Source: Report Section 3.2 [cite: 136]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, context_window_size, cost_per_million_tokens) VALUES
    ((SELECT id FROM llm_providers WHERE name = 'Mistral AI'), 'Mistral', 'Large 3', 'mistral-large-3', false, 256000, 0.50);

-- --- DeepSeek (Price Disruption) ---
-- Source: Report Section 3.3 [cite: 129, 136]
INSERT INTO llm_models (provider_id, family_name, version_name, api_model_name, is_local, context_window_size, cost_per_million_tokens) VALUES
    -- DeepSeek R1: The Reasoning Price Killer ($0.14!)
    ((SELECT id FROM llm_providers WHERE name = 'DeepSeek'), 'DeepSeek', 'R1 (Reasoning)', 'deepseek-reasoner', false, 128000, 0.14);


-- ============================================================
-- 3. SEED: CERTIFICATIONS (The 2026 Compliance Wall)
-- Source: Report Section 5 [cite: 146, 163, 187, 190]
-- ============================================================

INSERT INTO certifications (name, issuing_authority, description, official_link) VALUES
    (
        'ISO/IEC 42001:2023',
        'ISO',
        'The global standard for Artificial Intelligence Management Systems (AIMS). Mandatory for Enterprise Trust in 2026.',
        'https://www.iso.org/standard/81230.html'
    ),
    (
        'EU AI Act Compliant',
        'European Commission',
        'Full compliance with the High-Risk obligations (Article 6) applicable from Aug 2026.',
        'https://artificialintelligenceact.eu/'
    ),
    (
        'NIST AI RMF',
        'NIST (USA)',
        'Map, Measure, Manage, Govern framework. Required by US partners.',
        'https://www.nist.gov/itl/ai-risk-management-framework'
    ),
    (
        'SecNumCloud',
        'ANSSI (France)',
        'Sovereign Cloud qualification. Required for French Public Sector and OIV.',
        'https://www.ssi.gouv.fr/'
    ),
    (
        'SOC 2 Type II',
        'AICPA',
        'Audits controls relevant to security, availability, and processing integrity.',
        'https://www.aicpa.org/'
    ),
    (
        'OWASP LLM Top 10',
        'OWASP',
        'Technical hardening against Prompt Injection and Data Leakage.',
        'https://owasp.org/www-project-top-10-for-large-language-models/'
    )
ON CONFLICT (name) DO NOTHING;

-- ============================================================
-- 4. SEED: RESOURCE TYPES (Standard Tools)
-- ============================================================
INSERT INTO resource_types (id, name, config_schema, secret_schema) VALUES
    ('postgres_db', 'PostgreSQL Database', '{"host": "string", "port": "integer", "dbname": "string", "sslmode": "string"}', '{"username": "string", "password": "string"}'),
    ('aws_s3', 'AWS S3 Bucket', '{"bucket_name": "string", "region": "string"}', '{"access_key_id": "string", "secret_access_key": "string"}'),
    ('rest_api', 'REST API Endpoint', '{"base_url": "string", "timeout_seconds": "integer"}', '{"api_key": "string", "bearer_token": "string"}')
ON CONFLICT (id) DO NOTHING;

COMMIT;

SELECT '✅ DATA SEEDED: State of the Art 2026 (GPT-5.2, Llama 4, ISO 42001)' as status;
