include .env
export

MIGRATION_DIR       := db/migrations
TEST_DOCKER_COMPOSE := docker compose --file tests/docker-compose.yml

lint:
	golangci-lint run --allow-parallel-runners --config .golangci.yml ./...

# Example: make goose-create db=store name=init
goose-create:
	$(if $(value db),,$(error Database is not specified. Use "make goose-create db=yourdb name=yourname"))
	$(if $(value name),,$(error Migration name is not specified. Use "make goose-create db=yourdb name=yourname"))

	goose -v -dir $(MIGRATION_DIR)/$(db) create $(name) sql

goose-up-all:
	for db in $(POSTGRES_DATABASES); do (make goose-up db="$$db"); done

# Example: make goose-up db=store
goose-up:
	$(if $(value db),,$(error Database is not specified. Use "make goose-create db=yourdb"))

	-goose -v -dir $(MIGRATION_DIR)/$(db) postgres "postgresql://store:store@localhost:5432/$(db)?sslmode=disable" up
	goose -v -dir $(MIGRATION_DIR)/$(db) postgres "postgresql://store:store@localhost:5432/$(db)?sslmode=disable" status

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Запуск тестовой инфраструктуры для интеграционных тестов
test-infrastructure: test-infrastructure-down
	$(TEST_DOCKER_COMPOSE) up --detach
	$(TEST_DOCKER_COMPOSE) logs #--follow

# Завершение тестовой инфраструктуры
test-infrastructure-down:
	$(TEST_DOCKER_COMPOSE) down --remove-orphans

test:
	go test -v -coverprofile tests/cover.out -coverpkg=./... -race ./...
	grep -v "\.gen\.go\>" tests/cover.out | grep -v '_test\>' | grep -v '\<tests\>' > tests/cover.skipgen.out
	#go tool cover -func=tests/cover.skipgen.out #go tool cover -html=tests/cover.skipgen.out

# Example: make test-target name=TestGetDataV1DBFileURL
test-target:
	$(if $(value name),,$(error Test name is not specified. Example "make test-target name=TestAuth/TestLogout"))
	go test -v -race -run TestIntegration/$(name) ./tests/...
	@cat assets/succeeded.ascii