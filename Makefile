# Load .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variables
SQL_DIR := database/schemas
DB_CONTAINER := agentxmap_db
DB_NAME_IN_DOCKER := agentxmap_db
DB_USER_IN_DOCKER := postgres

# Helper to run psql inside docker
PSQL_CMD := docker exec -i $(DB_CONTAINER) psql -U $(DB_USER_IN_DOCKER) -d $(DB_NAME_IN_DOCKER)

# ==============================================================================
# TARGETS
# ==============================================================================

.PHONY: help db-check db-reset db-schema db-seed db-refresh docker-up docker-down run build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  docker-up   : Start DB container"
	@echo "  docker-down : Stop DB container"
	@echo "  db-check    : Verify connection"
	@echo "  db-reset    : ⚠️  DROP EVERYTHING"
	@echo "  db-schema   : Create tables (01_schema.sql)"
	@echo "  db-seed     : Populate with providers/certs (02_seed.sql)"
	@echo "  db-refresh  : FULL RESET -> SCHEMA -> SEED"
	@echo "  run         : Run API server"
	@echo "  build       : Build API server"

# ==============================================================================
# DATABASE (via Docker)
# ==============================================================================

db-check:
	@echo "Checking database connection..."
	@$(PSQL_CMD) -c "\conninfo" || (echo "❌ Connection failed"; exit 1)

db-reset:
	@echo "🧹 Cleaning database..."
	@cat $(SQL_DIR)/00_reset.sql | $(PSQL_CMD)
	@echo "✅ Database cleared."

db-schema:
	@echo "🏗️  Applying schema..."
	@cat $(SQL_DIR)/01_schema.sql | $(PSQL_CMD)
	@echo "✅ Schema applied."

db-seed:
	@echo "🌱 Seeding data..."
	@cat $(SQL_DIR)/02_seed.sql | $(PSQL_CMD)
	@echo "✅ Data seeded."

db-refresh: db-reset db-schema db-seed
	@echo "🚀 Database is fresh, seeded and ready for dev!"

# ==============================================================================
# DOCKER INFRA
# ==============================================================================

docker-up:
	@echo "🐳 Starting database container..."
	@docker compose up -d
	@echo "⏳ Waiting for database to be ready..."
	@sleep 2
	@echo "✅ Database is up."

docker-down:
	@echo "🛑 Stopping database container..."
	@docker compose down
	@echo "✅ Database stopped."

# ==============================================================================
# GO COMMANDS
# ==============================================================================

run:
	go run cmd/api/main.go

build:
	go build -o bin/server cmd/api/main.go
