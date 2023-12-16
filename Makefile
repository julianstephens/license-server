#!make

include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

DATABASE_URL="postgres://${LS_DATABASE_USER}:${LS_DATABASE_PASSWORD}@${MIGRATION_HOST}:5432/${LS_DATABASE_DB}?sslmode=disable&search_path=public"
MIGRATION_NAME ?= $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name')

.PHONY: up down stop lint docs-% db-%

up:
	@echo "[DCKR] Building and starting all containers"
	@docker compose up -d --build

stop:
	@echo "[DCKR] Stopping all containers"
	@docker compose stop

down:
	@echo "[DCKR] Destroying all containers"
	@docker compose down

lint:
	@echo "[GOCI] Linting project files"
	@golangci-lint run

build:
	@echo "[GOCI] Building CLI executable"
	@go build -o app ./cmd/cli.go

docs-swagger:
	@echo "[SWAG] Generating OpenAPI 2.0 schema"
	@swag init -g main.go

docs-convert:
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@openapi-generator-cli generate -i ./docs/swagger.yaml -o ./docs/v3 -g openapi-yaml --minimal-update
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@openapi-generator-cli generate -s -i ./docs/swagger.json -o ./docs/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@mv -f ./docs/v3/openapi/openapi.json ./docs/v3 ; mv -f ./docs/v3/openapi/openapi.yaml ./docs/v3 ; rm -rf ./docs/v3/openapi

docs-build:
	@echo "[REDC] Building Redocly API preview"
	@redocly build-docs ./docs/v3/openapi.yaml --config ./docs/redocly.yaml --output ./docs/v3/index.html

docs-gen:
	@make docs-swagger
	@make docs-convert
	@make docs-build

db-clean:
	@echo "[ATLS] Clearing DB"
	@atlas schema clean --env gorm --url $(DATABASE_URL) --auto-approve

db-push:
	@echo "[ATLS] Pushing DB schema to Atlas cloud"
	@atlas migrate push license-server --env gorm --dev-url $(DATABASE_URL)

db-apply:
	@echo "[ATLS] Applying schema to DB"
	@atlas schema apply --env gorm --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --auto-approve

db-schema:
	@echo "[ATLS] Opening DB schema viewer"
	@atlas schema inspect -w --env gorm --url $(DATABASE_URL)

db-refresh:
	@echo "[ATLS] Recreating DB"
	@make db-clean
	@make db-apply

db-lint:
	@echo "[ATLS] Linting DB"
	@atlas migrate lint --dev-url $(DATABASE_URL) -w

db-diff:
	@echo "[ATLS] Showing schema diff"
	@atlas schema diff --env gorm --from $(DATABASE_URL) --to "./migrations?format=atlas&version="

db-migration:
	@echo "[ATLS] Generating migration and applying to DB"
	@make db-clean
	@atlas migrate diff $(MIGRATION_NAME) --dev-url $(DATABASE_URL) --env gorm
	@make db-push
	@make db-apply
