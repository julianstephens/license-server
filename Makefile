#!make

include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.21-alpine
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

DATABASE_URL="postgres://${LS_POSTGRES_USER}:${LS_POSTGRES_PASSWORD}@${MIGRATION_HOST}:5432/${LS_POSTGRES_DB}?sslmode=disable&search_path=public"
MIGRATION_NAME ?= $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name')

.PHONY: swag up ci-swaggen

swag:
	watch -n 10 swag init -g main.go

up:
	docker compose up -d --build

openapi:
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "rm -rf /work/.openapi-generator"
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "cp -f /work/.openapi-generator-ignore /work/openapi"
	@echo "[OAS3] Converting Swagger 2-to-3 (json)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -s -i /work/swagger.json -o /work/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "mv -f /work/openapi/openapi.json /work ; mv -f /work/openapi/openapi.yaml /work ; rm -rf /work/openapi"

createdb:
	@echo "[ATLAS] Generating SQL and applying schema"
	@atlas schema apply --url $(DATABASE_URL) --dev-url $(DATABASE_URL) --env gorm

migration:
	@echo "[ATLAS] Generating migration and applying to DB"
	@atlas migrate diff $(MIGRATION_NAME) --dev-url $(DATABASE_URL) --env gorm
	# @atlas migrate push ls_app --dev-url $(DATABASE_URL)
	# @atlas migrate apply --env gorm --url $(DATABASE_URL) --tx-mode none

push:
	@atlas migrate push license-server --env gorm --dev-url $(DATABASE_URL)
	# @atlas migrate push license-server --dev-url docker://postgres/15/dev?search_path=public

apply:
	@atlas schema apply --env gorm --url $(DATABASE_URL) --dev-url $(DATABASE_URL)

schema:
	@echo "[ATLAS] Opening DB schema viewer"
	@atlas schema inspect -w --env gorm --url $(DATABASE_URL)

schemaclean:
	@echo "[ATLAS] Cleaning DB schema"
	@atlas schema clean --env gorm --url $(DATABASE_URL)

lint:
	@atlas migrate lint --dev-url $(DATABASE_URL) -w

# schemadiff:
# 	@echo "[ATLAS] Showing schema diff"
# 	@atlas schema diff --env gorm --from $(DATABASE_URL) --to "./migrations?format=atlas&version="

# lint:
# 	@echo "[GOLANGCI LINT] Checking project files"
# 	@golangci-lint run
