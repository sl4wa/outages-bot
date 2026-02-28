# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Telegram bot that notifies users about power outages in Lviv, Ukraine. Users subscribe via Telegram by selecting their street and building number; the notifier periodically fetches outage data from the Lviv power outage API and sends Telegram messages to affected subscribers. All user-facing text is in Ukrainian.

## Commands

```bash
make build          # Build binary to bin/outages-bot
make test           # Run all tests (go test ./...)
go test ./internal/domain/...   # Run tests for a single package
go test -run TestName ./internal/domain/  # Run a single test
```

Four CLI subcommands (cobra): `bot` (long-running Telegram bot), `notifier` (fetch outages & send notifications, optional `--interval` for looping), `outages` (print outages table), `users` (list subscribers).

## Architecture

Hexagonal architecture with strict dependency direction: `domain` ← `application` ← `cmd`/`client`/`repository`.

- **`internal/domain/`** — Core entities (`Outage`, `User`, `Street`, `UserAddress`, `OutageAddress`, `OutagePeriod`, `OutageDescription`, `OutageInfo`), repository interfaces (`UserRepository`, `StreetRepository`), and `FindOutageForNotification` (matches a user to their next unnotified outage). Pure logic, no external dependencies.
- **`internal/application/`** — Use cases and port interfaces. `ports.go` defines `OutageProvider`, `NotificationSender`, `UserInfoProvider`. `types.go` defines DTOs. Sub-packages: `notification/` (fetch + notify services), `subscription/` (search street, show/save subscription, unsubscribe), `admin/` (list users).
- **`internal/client/`** — External adapters. `outageapi/` fetches from Lviv power API. `telegram/` implements notification sender and user info provider.
- **`internal/repository/`** — File-based persistence. Users stored as individual YAML files (`data/users/<chatID>.yml`). Streets loaded from `data/streets.csv`. `NewFileUserRepository` runs a live `.txt` → `.yml` migration on startup.
- **`internal/cmd/`** — CLI command runners that wire everything together.
- **`internal/integration/`** — Integration tests using testify suites with real repositories and mocked Telegram API.
- **`main.go`** — Cobra root command setup and dependency wiring.

## Key Patterns

- Domain entities use constructor functions with validation (e.g., `NewUserAddress`, `NewOutagePeriod`) — always use these, don't construct structs directly.
- Repository `Save` uses atomic writes (temp file + rename).
- The bot uses an in-memory conversation state machine (`StepSearchStreet` → `StepSaveSubscription`) with TTL-based expiry; the zero state is simply an absent map entry.
- `NotificationService` auto-removes users who have blocked the bot (HTTP 403).
- Duplicate outages from the API are deduplicated by a composite key (street ID + buildings + time range).

## Environment Variables

- `TELEGRAM_BOT_TOKEN` — required for `bot`, `notifier`, and `users` commands
- `OUTAGE_API_URL` — required for `notifier` and `outages` commands
- `DATA_DIR` — data directory path (default: `data`)

## Testing

Tests use `github.com/stretchr/testify` (assert/require/suite). Integration tests in `internal/integration/` use testify suites with `SetupTest` for per-test isolation via `t.TempDir()`. The bot integration tests mock the Telegram API with `httptest.NewServer`.
