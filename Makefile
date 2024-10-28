DEV_ENV_SETUP_FOLDER ?= ./dev-env
DOCKER_COMPOSE_FILE ?= $(DEV_ENV_SETUP_FOLDER)/docker-compose.yml

VERSION ?= $(shell git rev-parse --short HEAD)

help:
	@echo "make version to get the current version"
	@echo "make start to start go-api server"
	@echo "make test to run the unit test"

version:
	@echo $(VERSION)

start:
	@echo "Starting app..."
	@docker-compose up --build -d

env-down:
	@docker-compose down


build:
	go build -v ./...

rebuild-docker:
	@docker-compose -f down
	@docker-compose -f build --no-cache
	@docker-compose -f up -d

stop:
	@echo "Stopping app..."
	@docker-compose down

watch: start
	@echo "Watching for file changes..."
	@docker-compose watch

logs:
	@docker-compose -f logs -f

test:
	go test -v ./... -short