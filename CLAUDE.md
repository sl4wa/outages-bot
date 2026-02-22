# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Telegram bot that notifies users about power outages in Lviv, Ukraine. Users subscribe via Telegram with their street and building number; a periodic notifier checks an external API and sends Telegram messages when a relevant outage is found.

## Common Commands

```bash
# Run all tests
go test ./...

# Run a single test file
go test ./internal/domain/ -run TestOutage

# Build the binary (use a temp dir for verification builds, not the project root)
go build -o /tmp/outages-bot .

# CLI commands
./outages-bot bot        # Run the Telegram bot (long-running)
./outages-bot notifier   # Fetch outages and send notifications (once)
./outages-bot notifier --interval=60s  # Run in a loop every 60s
./outages-bot outages    # Print a table of current outages
./outages-bot users      # List all subscribed users
```

## Architecture

The codebase follows a layered architecture under `internal/`:

- **domain** — Pure business logic with no external dependencies. Types (`User`, `Outage`, `Street`), value objects (`UserAddress`, `OutageAddress`, `OutagePeriod`, `OutageDescription`, `OutageInfo`), and domain services (`OutageFinder`).
- **application** — Application use cases and service orchestration. Organized by feature: `subscription/` (street search, save/show subscription), `notification/` (outage fetching, formatting, notification dispatch), `admin/` (CLI listing DTOs). Defines port interfaces in `ports.go` and shared types in `types.go`.
- **repository** — Data persistence: `FileUserRepository` and `CachedUserRepository` for users, `FileStreetRepository` for streets. All stored as flat files.
- **client** — External service clients: `outageapi/` fetches from the Lviv power outage API, `telegram/` contains notification sender and user info provider.
- **cmd** — Command runners: CLI commands (notifier, outages, users) and the Telegram bot runner with conversation state machine.

### Key Data Flow

1. `outages-bot notifier` (cron every 5 min) → `OutageFetchService` fetches outages from API → `NotificationService` iterates users, uses `OutageFinder` to match outages to user addresses, sends via `NotificationSender`, and records what was notified to avoid duplicates.
2. `outages-bot bot` (long-running) → Telegram bot handlers manage user subscriptions (street search, address saving, stop/unsubscribe).

### Persistence

No database — users are stored as individual text files (`data/users/{chatId}.txt`) with key-value lines. Streets are loaded from `data/streets.csv`. Outages are fetched live from the external API on each notifier run.

### Testing

Tests use in-memory test doubles and are co-located with the code they test (`_test.go` files). Integration tests are in `internal/integration/`.

## Environment

- Go 1.25+
- `TELEGRAM_BOT_TOKEN` env var — Telegram bot token (required for bot, notifier, users commands)
- `DATA_DIR` env var — path to data directory (default: `data`)
- supervisord runs two processes: `bot` (long-running Telegram bot) and `notifier` (long-running with `--interval` flag)
