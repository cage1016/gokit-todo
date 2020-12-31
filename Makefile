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
	docker-compose -f docker-compose.yaml down --volumes

.PHONY: generate
HAVE_GO_BINDATA := $(shell command -v mockgen 2> /dev/null)
generate: ## Regenerates GRPC proto and gomock
ifndef HAVE_GO_BINDATA
	@echo "requires 'mockgen' (GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4)"
	@exit 1 # fail
else
	go generate ./...
endif

.PHONY: test
test: ## test: run unit test
	# DEBUG=true bash -c "go test -v github.com/qeek-dev/retailbase/<package-name> -run ..."
	go test -v -race -cover -coverprofile coverage.txt -covermode=atomic ./...

.PHONY: test-integration
test-integration: ## test-integration: run integration test
	docker-compose -f docker-compose.integration.yaml up --build --abort-on-container-exit --remove-orphans

.PHONY: test-e2e
test-e2e: ## test-e2e: run e2e test
	docker-compose -f docker-compose.e2e.yaml up --build  --abort-on-container-exit --remove-orphans

.PHONY: help
help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_0-9-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help