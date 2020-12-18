.PHONY: run
run: stop up ## docker-compose stop & up

.PHONY: mod
mod: ## tidy go mod
	# This make rule requires Go 1.11+
	GO111MODULE=on go mod tidy

.PHONY: up
up: ## docker-compose up
	docker-compose -f docker-compose.yaml up -d --build

.PHONY: stop
stop: ## docker-compose stop
	docker-compose -f docker-compose.yaml stop

.PHONY: down
down: ## docker-compose down
	docker-compose -f docker-compose.yaml down

.PHONY: generate
HAVE_GO_BINDATA := $(shell command -v mockgen 2> /dev/null)
generate: ## Regenerates OPA data from rego files
ifndef HAVE_GO_BINDATA
	@echo "requires 'mockgen' (GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4)"
	@exit 1 # fail
else
	go generate ./...
endif

.PHONY: test
test: ## test: run unit test
	# DEBUG=true bash -c "go test -v github.com/qeek-dev/retailbase/<package-name> -run ..."
	go test -v -race -cover -coverprofile unit_cover.out ./...

.PHONY: integration
integration: ## integration: run integration test
	go test -v -race -tags=integration -coverprofile integration_cover.out ./...

.PHONY: e2e
e2e: ## e2e: run e2e test
	go test -v -race -tags=e2e -coverprofile integration_cover.out ./...

.PHONY: integration-test-db-up
integration-test-db-up: ## docker-compose integration up
	docker-compose -f docker-compose.integration.yaml up --build

.PHONY: integration-test-db-down
integration-test-db-down: ## docker-compose integration down
	docker-compose -f docker-compose.integration.yaml down --volumes

.PHONY: e2e-test-db-up
e2e-test-db-up: ## docker-compose e2e up
	docker-compose -f docker-compose.e2e.yaml up --build

.PHONY: e2e-test-db-down
e2e-test-db-down: ## docker-compose e2e down
	docker-compose -f docker-compose.e2e.yaml down --volumes

.PHONY: help
help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help