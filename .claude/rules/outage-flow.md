---
paths:
  - "cmd/outage-notification/**"
  - "internal/outage/**"
  - "test/integration/outage_*.go"
---

# Outage Flow Rules

The outage app is a Cobra CLI with four subcommands: `bot` (subscription Telegram bot), `notifier` (per-user outage notification daemon), `outages` (list outages), `users` (list users).

## Subscription bot

- Owns the full subscription conversation: street search, address selection, confirmation, `/start` recovery when lookup fails.
- Pending conversation state is in-memory with a lazy 30-minute TTL; inject or control `Workflow.now` in tests instead of using real-time waits.
- Keep all user-facing copy in Ukrainian unless the task explicitly changes product copy.
- Return structured `Response` values from `internal/outage/subscription`; Telegram API calls live in `internal/outage/telegram`.

## Outage notification pipeline

- Flow: fetch → normalize → deduplicate → match users → notify.
- Normalization and deduplication (street/buildings/start/end dedup key) happen in `internal/outage/loe` before downstream matching.
- Use existing constructors (`NewPeriod`, `NewAddress`, `NewDescription`) from `internal/outage/outage`; do not assemble invalid state directly.
- Use existing constructors and validation helpers from `internal/outage/users` for user/address/street entities.

## NotifyUsers semantics

- Unchanged outage snapshots short-circuit all notification work for that user.
- Users are removed only on `ErrRecipientUnavailable`; any other error leaves the user file intact.
- User outage state is saved only after a successful send.

## File persistence

- One user per `*.toml` file under `DATA_DIR/users` (defaults to `data/users`).
- All writes use atomic temp-file-plus-rename to avoid partial state.
- Outage snapshots use CSV format.

## Environment

- `OUTAGE_API_URL` — outage API endpoint.
- `DATA_DIR` — root data directory (defaults to `data`).
