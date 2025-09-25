# .PHONY — указывает, что цели не являются файлами
.PHONY: db-build db-up db-down goose migration-status migration-up migration-up-by-one migration-up-to migration-down migration-down-to migration-redo migration-create-sql migration-create-go debug-env

# ========== Настройки ==========
COMPOSE_ARGS := --env-file .env.local -f docker-compose.yml
COMPOSE      := docker compose $(COMPOSE_ARGS)
MIGRATIONS_DIR := ./migrations
GOOSE_SERVICE := goose

# ========== Базовые команды ==========
db-build:
	$(COMPOSE) build

db-up: db-build
	$(COMPOSE) up -d

db-down:
	$(COMPOSE) down

# Статус миграций
migration-status:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) status

# Применить все миграции
migration-up:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up

# Применить одну миграцию
migration-up-by-one:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up-by-one

# Применить до указанной версии
migration-up-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Ошибка: укажите VERSION, например: make up-to VERSION=5"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up-to $(VERSION)

# Откатить последнюю миграцию
migration-down:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) down

# Откатить до указанной версии
migration-down-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Ошибка: укажите VERSION, например: make down-to VERSION=3"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) down-to $(VERSION)

# Перезапустить последнюю миграцию
migration-redo:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) redo

# ========== Создание миграций ==========

# Создать SQL-миграцию
migration-create-sql:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Ошибка: укажите NAME, например: make create-sql NAME=add_users_table"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Создать Go-миграцию
migration-create-go:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Ошибка: укажите NAME, например: make create-go NAME=add_users_table"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) create $(NAME) go

# Убедимся что переменные из .env действительно попали в контейнер.
debug-env:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) env