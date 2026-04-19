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
  - `internal/telegram`: Telegram bot runner and Telegram transport helpers
  - `internal/users`: subscription and user-facing bot flows
- `main.go` is the composition root. It constructs Cobra commands and wires repositories, providers, notifier flow, and Telegram integrations.
- Treat notifier semantics, bot session flow, file persistence, and outage normalization rules as intentional behavior unless the task explicitly changes them.
