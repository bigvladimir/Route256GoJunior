POSTGRES_SETUP_TEST := user=postgres password=1 dbname=test host=localhost port=5432 sslmode=disable

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=1.2.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

INTEGRATION_TEST_FILES=$(CURDIR)/tests
DOCKER_COMPOSER_TEST_FILE=docker-compose-test.yaml

.PHONY: test
test:
	$(info Running tests...)	
	go test ./...


.PHONY: test-docker-compose-up
test-docker-compose-up:
	docker-compose -f $(DOCKER_COMPOSER_TEST_FILE) up -d

.PHONY: test-docker-compose-down
test-docker-compose-down:
	docker-compose -f $(DOCKER_COMPOSER_TEST_FILE) down

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: integration-test
integration-test:
	$(info Running integration tests...)
	go test -tags=integration $(INTEGRATION_TEST_FILES)


.PHONY: test-all
test-all: test test-docker-compose-up test-migration-up integration-test test-migration-down test-docker-compose-down
