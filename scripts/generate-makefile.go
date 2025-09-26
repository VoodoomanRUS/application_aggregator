package main

import (
	"os"
	"text/template"
)

const makefileTemplate = `# .PHONY — указывает, что цели не являются файлами
.PHONY: db-build db-up db-down goose  
.PHONY: migration-status migration-up migration-up-by-one migration-up-to migration-down migration-down-to migration-redo 
.PHONY: migration-create-sql migration-create-go create-core debug-env
.PHONY: test-db-up test-migrate-up test-status migrate-check
.PHONY: swagger-gen run-dev debug-env

# ========== Настройки ==========
MIGRATIONS_DIR := ./migrations

COMPOSE_ARGS := --env-file .env.local -f docker-compose.yml
COMPOSE      := docker compose $(COMPOSE_ARGS)
GOOSE_SERVICE := goose


# ========== Базовые команды ==========
db-build:
	$(COMPOSE) build

db-up: db-build
	$(COMPOSE) up -d

db-down:
	$(COMPOSE) down


# Статус миграций
status:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) status

# Применить все миграции
up:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up

# Применить одну миграцию
up-by-one:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up-by-one

# Применить до указанной версии
up-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Укажите VERSION: make up-to VERSION=5"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) up-to $(VERSION)

# Откатить последнюю миграцию
down:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) down

# Откатить до указанной версии
down-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Укажите VERSION: make down-to VERSION=3"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) down-to $(VERSION)

# Перезапустить последнюю миграцию
redo:
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

# Создать защищённую SQL-миграцию
create-core:
	@if [ -z "$$(NAME)" ]; then \
		echo "❌ Укажите NAME: make create-core NAME=create_users"; \
		exit 1; \
	fi
	@echo "🔒 Создание core-миграции: $$(NAME)"
	$(COMPOSE) run --rm $(GOOSE_SERVICE) goose -dir $(MIGRATIONS_DIR) create $$(NAME) sql
	@MIG_FILE=$$(ls -t $(MIGRATIONS_DIR)/[0-9]*_$$(NAME).sql | head -n1); \
	if [ -n "$$$$MIG_FILE" ]; then \
		cat scripts/create-core-template.sql > $$$$MIG_FILE; \
		echo "✅ Создана защищённая миграция: $$$$MIG_FILE"; \
	else \
		echo "❌ Не удалось найти созданный файл миграции"; \
		exit 1; \
	fi

# Убедимся что переменные из .env действительно попали в контейнер.
debug-env:
	$(COMPOSE) run --rm $(GOOSE_SERVICE) env | grep -E "GOOSE|POSTGRES"

# ========== Тестовая БД и валидация миграций ==========
test-db-up:
	@cp -n .env.test.example .env.test 2>/dev/null || true
	@docker compose --env-file .env.test -f docker-compose.test.yml up -d postgres-test

test-db-down:
	@docker compose --env-file .env.test -f docker-compose.test.yml down -v

migrate-check:
	@echo "🔍 Валидация миграций (накат → откат → накат)..."
	@cp -n .env.test.example .env.test 2>/dev/null || true
	@docker compose --env-file .env.test -f docker-compose.test.yml up -d postgres-test
	@sleep 3
	@echo "→ Накат всех миграций..."
	@docker compose --env-file .env.test -f docker-compose.test.yml run --rm goose-test goose -dir ./migrations up
	@echo "→ Откат всех миграций..."
	@docker compose --env-file .env.test -f docker-compose.test.yml run --rm goose-test goose -dir ./migrations down-to 0
	@echo "→ Повторный накат..."
	@docker compose --env-file .env.test -f docker-compose.test.yml run --rm goose-test goose -dir ./migrations up
	@echo "✅ Валидация пройдена! Все миграции идемпотентны."
	@docker compose --env-file .env.test -f docker-compose.test.yml down -v
`

func main() {
	tmpl := template.Must(template.New("makefile").Parse(makefileTemplate))
	f, err := os.Create("../Makefile")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, nil)
	if err != nil {
		panic(err)
	}
	println("✅ Makefile успешно сгенерирован!")
}
