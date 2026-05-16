# Core Repo Rules

outages-bot is a Go service and Telegram bot for outage lookup, subscriptions, and notifications.

- Think Before Coding: for broad or ambiguous work, clarify the goal, assumptions, affected packages, and verification path before editing.
- Simplicity: prefer the smallest implementation that fully solves the request using existing project patterns.
- Surgical Changes: keep edits scoped to the requested behavior; avoid unrelated refactors, rewrites, formatting churn, and dependency changes.
- Verification: decide the narrowest useful check before or during implementation, then run it or explain why it could not run.
- Keep package boundaries aligned with the current architecture:
  - `main.go` composes dependencies and Cobra commands
  - `internal/cli` runs command helpers and formats stdout output
  - `internal/loe` handles outage API transport, normalization, and cache behavior
  - `internal/outage` owns outage entities, validation, conversion, and snapshot comparison interfaces
  - `internal/subscription` owns subscription conversation workflow, command handling, pending state, and application operations
  - `internal/users` owns user, address, street entities, validation, repositories, and reusable matching/listing behavior
  - `internal/notifier` owns notification orchestration
  - `internal/telegram` owns Telegram transport, command mapping, response rendering, and reply keyboards
  - `internal/persistence` owns file-backed repositories and on-disk formats
- Use existing constructors and validation helpers in `internal/outage` and `internal/users` instead of assembling invalid state directly.
- In `NotifyUsers`, unchanged outage snapshots short-circuit notification work, users are removed only for `ErrRecipientUnavailable`, and user outage state is saved only after a successful send.
- File-backed persistence uses one user per `*.toml` file under the configured `DATA_DIR/users` directory, defaults to `data/users`, and saves users/outage snapshots with temp-file-plus-rename writes.
- Outage API normalization and deduplication happen before downstream matching, including the street/buildings/start/end dedup key.
- Keep user-facing bot text in Ukrainian unless the task explicitly changes product copy.
- Load the relevant path-scoped rules from `.claude/rules/` before making broader changes in a subsystem.
