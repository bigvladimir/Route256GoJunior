POSTGRES_SETUP := user=user password=1 dbname=pvzdb host=localhost port=5432 sslmode=disable
POSTGRES_SETUP_TEST := user=user password=1 dbname=test host=localhost port=5433 sslmode=disable

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=1.2.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

INTEGRATION_TEST_FILES=$(CURDIR)/tests
DOCKER_COMPOSER_TEST_FILE=docker-compose-test.yaml


.PHONY: build-docker
build-docker:
	docker-compose build

.PHONY: up-all
up-all:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down


.PHONY: build-server
build-server:
	go build -o ./cmd/httpserver/pvz.exe ./cmd/httpserver

.PHONY: run-server
run-server:
	./cmd/httpserver/pvz.exe


.PHONY: test
test:
	$(info Running tests...)	
	go test ./...


.PHONY: integration-test-run
integration-test-run:
	$(info Running integration tests...)
	go test -tags=integration $(INTEGRATION_TEST_FILES)


.PHONY: test-build-docker
test-build-docker:
	docker-compose -f $(DOCKER_COMPOSER_TEST_FILE) build

.PHONY: test-compose-up
test-compose-up:
	docker-compose -f $(DOCKER_COMPOSER_TEST_FILE) up -d

.PHONY: test-compose-down
test-compose-down:
	docker-compose -f $(DOCKER_COMPOSER_TEST_FILE) down

.PHONY: wait-for-containers
wait-for-containers:
	$(info Waiting for containers...)
	timeout 3


.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down


.PHONY: integration-test
integration-test: test-build-docker test-compose-up wait-for-containers test-migration-up integration-test-run test-migration-down test-compose-down

.PHONY: lint
lint:
	golangci-lint run --config=.golangci.yaml ./...
