#!make

include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

DATABASE_URL="postgres://${LS_POSTGRES_USER}:${LS_POSTGRES_PASSWORD}@${MIGRATION_HOST}:5432/${LS_POSTGRES_DB}?sslmode=disable&search_path=public"
MIGRATION_NAME ?= $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name')

.PHONY:

swagger:
	@echo "[SWAG] Generating OpenAPI 2.0 schema"
	@swag init -g main.go
	@make openapi
	@make redoc

openapi:
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@openapi-generator-cli generate -i ./docs/swagger.yaml -o ./docs/v3 -g openapi-yaml --minimal-update
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@openapi-generator-cli generate -s -i ./docs/swagger.json -o ./docs/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@mv -f ./docs/v3/openapi/openapi.json ./docs/v3 ; mv -f ./docs/v3/openapi/openapi.yaml ./docs/v3 ; rm -rf ./docs/v3/openapi

redoc:
	@echo "[REDOC] Building Redocly API preview"
	@redocly build-docs ./docs/v3/openapi.yaml --config ./docs/redocly.yaml --output ./docs/v3/index.html

up:
	docker compose up -d --build

createdb:
	@echo "[ATLAS] Generating SQL and applying schema"
	@atlas schema apply --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --env gorm

migration:
	@echo "[ATLAS] Generating migration and applying to DB"
	@atlas schema clean --env gorm --url $(DATABASE_URL) --auto-approve
	@atlas migrate diff $(MIGRATION_NAME) --dev-url $(DATABASE_URL) --env gorm
	@atlas migrate push license-server --env gorm --dev-url $(DATABASE_URL)
	@atlas schema apply --env gorm --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --auto-approve

push:
	@atlas migrate push license-server --env gorm --dev-url $(DATABASE_URL)
	# @atlas migrate push license-server --dev-url docker://postgres/15/dev?search_path=public

apply:
	@atlas schema apply --env gorm --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --auto-approve

schema:
	@echo "[ATLAS] Opening DB schema viewer"
	@atlas schema inspect -w --env gorm --url $(DATABASE_URL)

refresh:
	@echo "[ATLAS] Cleaning DB schema"
	@atlas schema clean --env gorm --url $(DATABASE_URL) --auto-approve
	@atlas schema apply --env gorm --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --auto-approve

lint:
	@atlas migrate lint --dev-url $(DATABASE_URL) -w

stop:
	@echo "[DOCKER] Stopping all containers"
	@docker compose stop

# schemadiff:
# 	@echo "[ATLAS] Showing schema diff"
# 	@atlas schema diff --env gorm --from $(DATABASE_URL) --to "./migrations?format=atlas&version="

# lint:
# 	@echo "[GOLANGCI LINT] Checking project files"
# 	@golangci-lint run
