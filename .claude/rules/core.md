# Core Repo Rules

outages-bot is a Go service and Telegram bot for outage lookup, subscriptions, and notifications.

- Keep changes small and preserve current behavior unless the task explicitly changes it.
- Keep package boundaries aligned with the current architecture:
  - `main.go` composes dependencies and Cobra commands
  - `internal/cli` runs command helpers and formats stdout output
  - `internal/loe` handles outage API transport, normalization, and cache behavior
  - `internal/outage` owns outage entities, validation, conversion, and snapshot comparison interfaces
  - `internal/users` owns subscription, lookup, and user-facing bot flow logic
  - `internal/notifier` owns notification orchestration
  - `internal/telegram` owns Telegram transport and the in-memory bot conversation state machine
  - `internal/persistence` owns file-backed repositories and on-disk formats
- Use existing constructors and validation helpers in `internal/outage` and `internal/users` instead of assembling invalid state directly.
- `internal/telegram/bot_runner.go` intentionally owns the in-memory Telegram conversation state machine, cleanup loop, and default 30-minute session TTL.
- In `NotifyUsers`, unchanged outage snapshots short-circuit notification work, users are removed only for `ErrRecipientUnavailable`, and user outage state is saved only after a successful send.
- Preserve file-backed persistence behavior, especially one user per `*.yml` file under the configured `DATA_DIR/users` directory, with default `data/users`, and atomic temp-file-plus-rename writes for both user and outage snapshot saves.
- Preserve outage API normalization and deduplication behavior before downstream matching, including the street/buildings/start/end dedup key.
- Keep user-facing bot text in Ukrainian unless the task explicitly changes product copy.
- Add or update targeted Go tests when behavior changes.
- Load the relevant path-scoped rules from `.claude/rules/` before making broader changes in a subsystem.
