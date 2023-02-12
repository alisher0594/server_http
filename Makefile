include .env
export

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

help: ### help: print this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-up-app: ### Run app on docker
	docker-compose up --build -d postgres app && docker-compose logs -f
.PHONY: compose-up-app

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

test: ### run unit tests
	@echo 'Running unit tests...'
	go test -v -cover -race ./internal/...
	@echo 'Running end to end tests...'
	go test -v -cover -race ./cmd/api/...
.PHONY: test

run: ### run app
	go mod tidy && go mod download && \
	go run -tags migrate ./cmd/api -db-dsn='postgres://user:pass@localhost:5432/postgres?sslmode=disable'
.PHONY: run

migrate-create:  ### create new migration
	migrate create -seq -ext=.sql -dir=./migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path=./migrations -database=$(PG_URL) up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path=./migrations -database=$(PG_URL) down 1
.PHONY: migrate-down
