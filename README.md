# Outages Bot

A single Go module (`github.com/sl4wa/outages-bot`) producing two binaries for Lviv power outage Telegram notifications.

## Layout

- `cmd/outage-notification/` — subscription bot and current-outage notifier (Cobra subcommands: `bot`, `notifier`, `outages`, `users`).
- `cmd/schedule-notification/` — schedule poller and broadcaster.
- `internal/outage/` — outage app's domain code (cli, loe, notifier, outage, persistence, subscription, telegram, users).
- `internal/schedule/` — schedule app's domain code (loe, message, notifier, persistence, schedule, telegram).
- `internal/shared/` — code reused by both apps: `httpcache`, `subscribers`, `telegram`.
- `test/integration/` — outage-app integration tests.

Boundary: `internal/outage/` and `internal/schedule/` must not import each other. Both may import `internal/shared/`.

## Build and run

```bash
make build           # bin/outage-notification, bin/schedule-notification
make test            # go test ./...
make run-bot         # go run ./cmd/outage-notification bot
make run-notifier    # go run ./cmd/outage-notification notifier
make run-schedule    # go run ./cmd/schedule-notification
```

Both binaries are expected to run from the repo root. Configuration is read from `.env` plus environment variables:

- Outage app (`cmd/outage-notification`): `TELEGRAM_BOT_TOKEN`, `OUTAGE_API_URL`, `DATA_DIR` (defaults to `data`).
- Schedule app (`cmd/schedule-notification`): `TELEGRAM_BOT_TOKEN`, `SCHEDULE_API_URL`, `DATA_DIR` (uses `schedule.csv`, `users/`, and `schedule.http-cache` under that directory).

See `.env.example` for a starter file.
