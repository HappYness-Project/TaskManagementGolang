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
