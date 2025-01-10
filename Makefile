DEV_ENV_SETUP_FOLDER ?= ./dev-env
DOCKER_COMPOSE_FILE ?= $(DEV_ENV_SETUP_FOLDER)/docker-compose.yml

VERSION ?= $(shell git rev-parse --short HEAD)

help:
	@echo "make version to get the current version"
	@echo "make start to start go-api server"
	@echo "make build"
	@echo "make rebuild-docker"
	@echo "make logs"
	@echo "make down to remove docker containers"
	@echo "make test to run the unit test"

version:
	@echo $(VERSION)

start:
	@echo "Starting app..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

down:
	@echo "Stopping app..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

build:
	go build -v ./...

rebuild-docker:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build --no-cache
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

watch: start
	@echo "Watching for file changes..."
	@docker-compose watch

logs:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

test:
	go test -v ./... -short