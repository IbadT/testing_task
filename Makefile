include .env
export


GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RED    := $(shell tput -Txterm setaf 1)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# DATABASE_URL="host=localhost user=postgres password=postgres dbname=testingtask port=5432 sslmode=disable"
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/testingtask?sslmode=disable
MIGRATIONS_DIR=./migrations

# Docker переменные
DOCKER_COMPOSE = docker compose
DOCKER_SERVICE_APP = app
DOCKER_SERVICE_POSTGRES = postgres
# DOCKER_DB_URL = postgres://postgres:postgres@postgres:5432/testingtask?sslmode=disable
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${HOST_POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}

run:
	go run cmd/main.go

build:
	go build cmd/.

gen:
	oapi-codegen -config openapi/.openapi -include-tags subscriptions -package subscriptions openapi/openapi.yaml > ./internal/web/subscriptions/api.gen.go

gen-docs:
	pwd
	swag init -g cmd/main.go

rm-gen-docs:
	rm -rf ./docs

lint:
	golangci-lint run --color=always

migrate-new:
	migrate create -ext sql -dir ${MIGRATIONS_DIR} ${NAME}

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-reup:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

db-create:
	docker-compose exec -T postgres createdb -U postgres testingtask 2>/dev/null || true

# Docker команды
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-up-build:
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_COMPOSE) down

docker-down-clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-logs-app:
	$(DOCKER_COMPOSE) logs -f $(DOCKER_SERVICE_APP)

docker-logs-postgres:
	$(DOCKER_COMPOSE) logs -f $(DOCKER_SERVICE_POSTGRES)

docker-ps:
	$(DOCKER_COMPOSE) ps

docker-restart:
	$(DOCKER_COMPOSE) restart

docker-restart-app:
	$(DOCKER_COMPOSE) restart $(DOCKER_SERVICE_APP)

docker-build:
	$(DOCKER_COMPOSE) build

docker-build-app:
	$(DOCKER_COMPOSE) build $(DOCKER_SERVICE_APP)

docker-exec:
	$(DOCKER_COMPOSE) exec $(DOCKER_SERVICE_APP) sh

docker-exec-postgres:
	$(DOCKER_COMPOSE) exec $(DOCKER_SERVICE_POSTGRES) psql -U postgres -d testingtask


# Миграции в Docker
docker-migrate-up:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" up

docker-migrate-down:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" down 1

docker-migrate-down-all:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" down

docker-migrate-force:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" force $(VERSION)

docker-migrate-version:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" version

docker-migrate-new:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate create -ext sql -dir /app/migrations $(NAME)

docker-migrate-reup:
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DOCKER_DB_URL)" down
	$(DOCKER_COMPOSE) run --rm $(DOCKER_SERVICE_APP) migrate -path /app/migrations -database "$(DATABASE_URL)" up

# Очистка Docker
docker-clean:
	docker system prune -f

docker-clean-all:
	docker system prune -a -f --volumes

# Полный цикл разработки
# dev-up: docker-up-build docker-migrate-up
# 	@echo "$(GREEN)✅ Приложение запущено и миграции применены$(RESET)"

dev-up: docker-up-build
	@if $(MAKE) docker-migrate-up; then \
		echo "$(GREEN)✅ Приложение запущено и миграции применены$(RESET)"; \
	else \
		echo "$(RED)❌ Ошибка: миграции не применились!$(RESET)"; \
		exit 1; \
	fi


dev-down: docker-down
	@echo "$(YELLOW)🛑 Приложение остановлено$(RESET)"

dev-rebuild: docker-down docker-up-build docker-migrate-up
	@echo "$(GREEN)✅ Приложение пересобрано и запущено$(RESET)"



# Проверка здоровья
docker-health:
	@echo "$(YELLOW)Проверка статуса контейнеров:$(RESET)"
	$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "$(YELLOW)Проверка доступности БД:$(RESET)"
	$(DOCKER_COMPOSE) exec $(DOCKER_SERVICE_POSTGRES) pg_isready -U postgres

.PHONY: run build gen gen-docs rm-gen-docs lint db-create
.PHONY: docker-up docker-up-build docker-down docker-down-clean docker-logs docker-logs-app docker-logs-postgres docker-ps docker-restart docker-restart-app docker-build docker-build-app docker-exec docker-exec-postgres
.PHONY: docker-migrate-up docker-migrate-down docker-migrate-down-all docker-migrate-force docker-migrate-version docker-migrate-new docker-migrate-reup
.PHONY: docker-clean docker-clean-all dev-up dev-down dev-rebuild docker-health






# GREEN  := $(shell tput -Txterm setaf 2)
# YELLOW := $(shell tput -Txterm setaf 3)
# WHITE  := $(shell tput -Txterm setaf 7)
# RESET  := $(shell tput -Txterm sgr0)

# # DATABASE_URL="host=localhost user=postgres password=postgres dbname=testingtask port=5432 sslmode=disable"
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/testingtask?sslmode=disable
# MIGRATIONS_DIR=./migrations

# run:
# 	go run cmd/main.go

# build:
# 	go build cmd/.

# gen:
# 	oapi-codegen -config openapi/.openapi -include-tags subscriptions -package subscriptions openapi/openapi.yaml > ./internal/web/subscriptions/api.gen.go

# gen-docs:
# 	pwd
# 	swag init -g cmd/main.go

# rm-gen-docs:
# 	rm -rf ./docs

# lint:
# 	golangci-lint run --color=always

# db-create:
# 	docker-compose exec -T postgres createdb -U postgres testingtask 2>/dev/null || true


# migrate-new:
# 	migrate create -ext sql -dir ${MIGRATIONS_DIR} ${NAME}

# migrate-up:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# migrate-down:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# migrate-down-all:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

# migrate-force:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)

# migrate-version:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# migrate-reup:
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down
# 	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up