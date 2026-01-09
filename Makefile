# ... (VARIABLES & CONFIG - inchangé) ...
SQL_DIR := database/schemas

# ==============================================================================
# TARGETS
# ==============================================================================

.PHONY: help db-check db-reset db-schema db-seed db-refresh

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  db-check    : Verify connection"
	@echo "  db-reset    : ⚠️  DROP EVERYTHING"
	@echo "  db-schema   : Create tables (01_schema.sql)"
	@echo "  db-seed     : Populate with providers/certs (02_seed.sql)"
	@echo "  db-refresh  : FULL RESET -> SCHEMA -> SEED"

db-check:
	@echo "Checking database connection..."
	@psql "$(DB_DSN)" -c "\conninfo" || (echo "❌ Connection failed"; exit 1)

db-reset:
	@echo "🧹 Cleaning database..."
	@psql "$(DB_DSN)" -f $(SQL_DIR)/00_reset.sql
	@echo "✅ Database cleared."

db-schema:
	@echo "🏗️  Applying schema..."
	@psql "$(DB_DSN)" -f $(SQL_DIR)/01_schema.sql
	@echo "✅ Schema applied."

# NOUVELLE COMMANDE
db-seed:
	@echo "🌱 Seeding data..."
	@psql "$(DB_DSN)" -f $(SQL_DIR)/02_seed.sql
	@echo "✅ Data seeded."

# MISE A JOUR DE LA COMMANDE
db-refresh: db-reset db-schema db-seed
	@echo "🚀 Database is fresh, seeded and ready for dev!"

# ==============================================================================
# GO COMMANDS (Bonus)
# ==============================================================================
.PHONY: run build

run:
	go run cmd/api/main.go

build:
	go build -o bin/server cmd/api/main.go
