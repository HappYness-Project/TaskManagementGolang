ifeq (,$(wildcard ./.env))
$(error .env file not found)
else
include .env
export $(shell sed 's/=.*//' .env)
endif


VERSION ?= $(shell git rev-parse --short HEAD)

version:
	@echo $(VERSION)

start:
	@echo "Starting app..."
	@docker-compose up --build -d

stop:
	@echo "Stopping app..."
	@docker-compose down

watch: start
	@echo "Watching for file changes..."
	@docker-compose watch

logs:
	@docker-compose logs -f

build:
	go build -v ./...

test:
	go test -v ./... -short