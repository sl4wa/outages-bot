# CLAUDE.md

outages-bot is a Go Telegram bot for outage subscriptions and notifications.

Repo-specific Claude Code guidance lives in `.claude/rules/`.

- Start with `.claude/rules/core.md` and `.claude/rules/commands.md`.
- Claude Code also loads path-scoped rules from `.claude/rules/` when working with matching files.
- Current package layout:
  - `internal/cli`: command execution helpers and stdout formatting
  - `internal/loe`: outage API transport, normalization, and cache behavior
  - `internal/notifier`: notification orchestration
  - `internal/outage`: outage models, validation, conversion, and snapshot interfaces
  - `internal/persistence`: file-backed user, street, and outage snapshot storage
  - `internal/subscription`: subscription conversation workflow, command handling, pending state, and application operations
  - `internal/telegram`: Telegram bot runner, command mapping, response rendering, and Telegram transport helpers
  - `internal/users`: user, address, street entities, validation, repositories, and reusable matching/listing behavior
- `main.go` is the composition root. It constructs Cobra commands and wires repositories, providers, notifier flow, and Telegram integrations.
- Treat notifier semantics, subscription conversation flow, file persistence, and outage normalization rules as intentional behavior unless the task explicitly changes them.
