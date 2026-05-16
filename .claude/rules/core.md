# Core Repo Rules

outages-bot is a Go monorepo with two binaries built from one module: an outage-lookup/subscription Telegram bot and a schedule broadcaster. Both share helpers under `internal/shared`.

- Think Before Coding: for broad or ambiguous work, clarify the goal, assumptions, affected packages, and verification path before editing.
- Simplicity: prefer the smallest implementation that fully solves the request using existing project patterns.
- Surgical Changes: keep edits scoped to the requested behavior; avoid unrelated refactors, rewrites, formatting churn, and dependency changes.
- Verification: decide the narrowest useful check before or during implementation, then run it or explain why it could not run.
- Module path: `github.com/sl4wa/outages-bot`. Single `go.mod` at the repo root.
- Repository layout:
  - `cmd/outage-notification/` — entrypoint for the bot/notifier/outages/users Cobra app
  - `cmd/schedule-notification/` — entrypoint for the schedule broadcaster
  - `internal/outage/...` — owns the outage app's parsing and notification path
    - `cli`, `loe`, `notifier`, `outage`, `persistence`, `subscription`, `telegram`, `users`
  - `internal/schedule/...` — owns the schedule app's parsing and broadcast path
    - `loe`, `message`, `notifier`, `persistence`, `schedule`, `telegram`
  - `internal/shared/...` — code reused across both apps: `httpcache`, `subscribers`, `telegram`
  - `test/integration/` — cross-package integration tests (outage_*.go for the outage flow, schedule_*.go for the schedule flow)
- **Boundary convention**: code under `internal/outage/` and `internal/schedule/` MUST NOT import each other. Both may import `internal/shared/`. Cross-tree imports are a smell; raise them in review.
- Both binaries are expected to run from the repo root. Supervisord uses `directory=%(here)s` and the Makefile uses `go run ./cmd/...`; only one `godotenv.Load()` (no `../.env` fallback).
- Load the relevant path-scoped rules from `.claude/rules/` before making broader changes in a subsystem. Key flow-level rules: `outage-flow.md` and `schedule-flow.md`.
