run: stop up

mod:
	# This make rule requires Go 1.11+
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

up:
	docker-compose -f docker-compose.yaml up -d --build

stop:
	docker-compose -f docker-compose.yaml stop

down:
	docker-compose -f docker-compose.yaml down

# Regenerates OPA data from rego files
HAVE_GO_BINDATA := $(shell command -v mockgen 2> /dev/null)
generate:
ifndef HAVE_GO_BINDATA
	@echo "requires 'mockgen' (GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4)"
	@exit 1 # fail
else
	go generate ./...
endif

## test: run unit test
test:
	# DEBUG=true bash -c "go test -v github.com/qeek-dev/retailbase/<package-name> -run ..."
	go test -v -race -cover -coverprofile unit_cover.out ./...

## integration: run integration test
integration:
	go test -v -race -tags=integration -coverprofile integration_cover.out ./...

test-db-up:
	docker-compose -f docker-compose.test.yaml up --build

test-db-down:
	docker-compose -f docker-compose.test.yaml down --volumes

.PHONY: generate test