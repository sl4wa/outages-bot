-include .env
export

.PHONY: help build test run-bot run-notifier run-outages run-users clean

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build          Build the binary to bin/outages-bot"
	@echo "  test           Run all tests"
	@echo "  run-bot        Run the Telegram bot"
	@echo "  run-notifier   Run the notifier"
	@echo "  run-outages    Print current outages"
	@echo "  run-users      List subscribed users"
	@echo "  clean          Remove build artifacts"
	@echo "  help           Show this help"

build:
	go build -o bin/outages-bot .

test:
	go test ./...

run-bot:
	go run . bot

run-notifier:
	go run . notifier

run-outages:
	go run . outages

run-users:
	go run . users

clean:
	rm -rf bin/
