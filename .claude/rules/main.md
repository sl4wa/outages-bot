---
paths:
  - "cmd/**/main.go"
---

# Main Rules

- Each `cmd/<app>/main.go` is the composition root for its binary.
- Keep Cobra command construction and dependency wiring centralized in `main.go`; do not push wiring into `internal/` packages.
- Prefer wiring existing package APIs together rather than re-implementing business logic in `main.go`.
- Both binaries are expected to run from the repo root: a single `godotenv.Load()` (no `../.env` fallback), and default file paths are root-relative (e.g. `data/users`, `data/outages.csv`).
- Preserve environment and data-directory behavior unless the task explicitly changes runtime configuration.
