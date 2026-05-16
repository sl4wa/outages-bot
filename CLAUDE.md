# CLAUDE.md

outages-bot is a Go monorepo with two binaries built from one module (`github.com/sl4wa/outages-bot`):

- `cmd/outage-notification/` — Cobra app exposing `bot`, `notifier`, `outages`, `users` subcommands. Owns the subscription Telegram bot and per-user outage notification flow.
- `cmd/schedule-notification/` — polling daemon that broadcasts LOE schedule changes to every user in the shared users directory.

Repo-specific Claude Code guidance lives in `.claude/rules/`.

- Start with `.claude/rules/core.md` and `.claude/rules/commands.md`.
- For outage-flow work, load `.claude/rules/outage-flow.md` (auto-scoped to `cmd/outage-notification/**`, `internal/outage/**`, `test/integration/outage_*.go`).
- For schedule-flow work, load `.claude/rules/schedule-flow.md` (auto-scoped to `cmd/schedule-notification/**`, `internal/schedule/**`, `test/integration/schedule_*.go`).
- **Boundary**: `internal/outage/` and `internal/schedule/` must not import each other. Both may import `internal/shared/`.
- Both binaries run from the repo root (supervisord uses `directory=%(here)s`; Makefile uses `go run ./cmd/...`). They load env via a single `godotenv.Load()` and default file paths are root-relative.
