.PHONY: help build test run-bot run-notifier run-outages run-users run-schedule run-schedule-loop clean

-include .env
export

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build both binaries to bin/"
	@echo "  test               Run all tests"
	@echo "  run-bot            Run the Telegram bot"
	@echo "  run-notifier       Run the outage notifier"
	@echo "  run-outages        Print current outages"
	@echo "  run-users          List subscribed users"
	@echo "  run-schedule       Run the schedule notifier once"
	@echo "  run-schedule-loop  Run the schedule notifier every 60s"
	@echo "  clean              Remove build artifacts"

build:
	mkdir -p bin
	go build -o bin/outage-notification ./cmd/outage-notification
	go build -o bin/schedule-notification ./cmd/schedule-notification

test:
	go test ./...

run-bot:
	go run ./cmd/outage-notification bot

run-notifier:
	go run ./cmd/outage-notification notifier

run-outages:
	go run ./cmd/outage-notification outages

run-users:
	go run ./cmd/outage-notification users

run-schedule:
	go run ./cmd/schedule-notification

run-schedule-loop:
	go run ./cmd/schedule-notification --interval=60s

clean:
	rm -rf bin/
