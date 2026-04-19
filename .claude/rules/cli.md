---
paths:
  - "internal/cli/**/*.go"
---

# CLI Rules

- `internal/cli` owns command execution helpers and stdout formatting only.
- Do not move Cobra command construction or dependency wiring into `internal/cli`; keep that in `main.go`.
- Do not move Telegram transport, outage API transport, or notification orchestration into this package.
- Preserve current command output behavior unless the task explicitly changes the CLI contract.
