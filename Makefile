SHELL := /bin/sh

# Image and compose settings
IMAGE ?= poweron-app
TAG ?= latest
COMPOSE ?= docker compose -f docker-compose.yml

.PHONY: help build up down ps logs logs-bot logs-notifier restart shell-bot shell-notifier notifier-once clean test clean-all

help:
	@echo "Available targets:";
	@echo "  build           Build Docker image $(IMAGE):$(TAG)";
	@echo "  up              Start services with docker compose";
	@echo "  down            Stop and remove services";
	@echo "  ps              List compose services";
	@echo "  logs            Tail logs for all services";
	@echo "  logs-bot        Tail logs for bot";
	@echo "  logs-notifier   Tail logs for notifier";
	@echo "  restart         Restart services";
	@echo "  shell-bot       Open a shell in the bot container";
	@echo "  shell-notifier  Open a shell in the notifier container";
	@echo "  notifier-once   Run notifier once (no cron) via bot container";
	@echo "  test            Build image and run PHPUnit in notifier container";
	@echo "  clean           Remove local image $(IMAGE):$(TAG)";
	@echo "  clean-all       Stop, clean image, build, and start services";

# Image operations
build:
	docker build -f Dockerfile -t $(IMAGE):$(TAG) .

# Compose operations
up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f

logs-bot:
	$(COMPOSE) logs -f bot

logs-notifier:
	$(COMPOSE) logs -f notifier

restart:
	$(COMPOSE) restart

# Interactive helpers
shell-bot:
	$(COMPOSE) exec bot sh

shell-notifier:
	$(COMPOSE) exec notifier sh

notifier-once:
	$(COMPOSE) exec bot php bin/console app:notifier

clean:
	@docker rmi $(IMAGE):$(TAG) || true
	@echo "Removed image if it existed: $(IMAGE):$(TAG)"

# Tests (run in Docker using notifier service image)
test:
	$(COMPOSE) build notifier
	$(COMPOSE) run --rm notifier vendor/bin/phpunit tests --bootstrap tests/bootstrap.php --colors=always

# Compound command
clean-all: down clean build up
	@echo "All clean, build, and up complete."
