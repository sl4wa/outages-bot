---
paths:
  - "cmd/schedule-notification/**"
  - "internal/schedule/**"
  - "test/integration/schedule_*.go"
---

# Schedule Flow Rules

The schedule app fetches schedule HTML from the LOE menu API, selects the latest schedule for today and tomorrow, broadcasts Telegram HTML notifications when schedule text changes, and persists current state to CSV.

## Broadcaster semantics

- `internal/schedule` and `cmd/schedule-notification` form an intentionally fire-and-forget broadcaster. It is **not** a subscription bot.
- Reads the shared users dir (`data/users` by default) and broadcasts schedule changes to every user file it finds.
- Does **not** own subscription state, subscribe/unsubscribe users, or prune blocked users on `ErrRecipientUnavailable` — that ownership lives in the outage app.
- Do not flag the absence of subscription handling, blocked-user removal, or per-user state persistence as inconsistency — it is the intended scope boundary.

## Key invariants

- Save-before-notify order in `internal/schedule/notifier/runner.go`: state is saved first so a notify failure does not re-broadcast on retry. See `TestRunnerSavesBeforeNotifyingAndPropagatesNotifyError`.
- HTTP cache in `internal/schedule/loe/cache.go` defers `Save` until after parse succeeds — a failed parse leaves the cache unchanged.
- Boundary: must not import `internal/outage/...`.

## Real concerns to still raise

- Swallowed I/O errors.
- Message formatting (`internal/schedule/message`) and parsing correctness (`internal/schedule/loe/provider.go`).

## Environment

- `SCHEDULE_API_URL` — LOE schedule API endpoint.
- `SCHEDULE_STATE_PATH` — CSV state file path (defaults to `data/outages.csv`).
- `SCHEDULE_HTTP_CACHE_PATH` — HTTP cache path (derived from state path when unset).
- `TELEGRAM_USERS_DIR` — shared users directory (defaults to `data/users`).
