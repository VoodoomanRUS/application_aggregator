DEV_COMPOSE_ARGS=--env-file .env.local -f Docker-compose.yml
DEV_COMPOSE_ENV=docker compose $(DEV_COMPOSE_ARGS)
DEV_COMPOSE=docker compose $(DEV_COMPOSE_ARGS)

MIGRATE_CMD = go run main.go

db-build:
	$(DEV_COMPOSE) build
db-up: db-build
	$(DEV_COMPOSE) --env-file .env.local up -d
db-down:
	$(DEV_COMPOSE) down
migration-up:
	@echo "Applying migrations..."
	${MIGRATE_CMD} up
migration-down:
	@echo "Rolling back migrations..."
	${MIGRATE_CMD} down
status:
	@echo "Checking db status..."
	${MIGRATE_CMD} status