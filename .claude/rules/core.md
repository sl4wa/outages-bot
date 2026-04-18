# Core Repo Rules

outages-bot is a Go service and Telegram bot for outage lookup, subscriptions, and notifications.

- Keep changes small and preserve current behavior unless the task explicitly changes it.
- Keep business rules in `internal/domain`; keep CLI, Telegram, HTTP, and file I/O concerns out of domain code.
- Use existing constructors and validation helpers for domain values instead of assembling invalid state directly.
- `internal/cmd/bot.go` intentionally owns the in-memory Telegram conversation state machine, cleanup loop, and default 30-minute session TTL.
- In `NotifyUsers`, unchanged outage snapshots short-circuit notification work, users are removed only for `ErrRecipientUnavailable`, and user outage state is saved only after a successful send.
- Preserve file-backed persistence behavior, especially one-user-per-file storage and `FileUserRepository.Save` atomic temp-file-plus-rename writes.
- Preserve outage API normalization and deduplication behavior before downstream matching, including the street/buildings/start/end dedup key.
- Keep user-facing bot text in Ukrainian unless the task explicitly changes product copy.
- Add or update targeted Go tests when behavior changes.
- Load the relevant path-scoped rules from `.claude/rules/` before making broader changes in a subsystem.
