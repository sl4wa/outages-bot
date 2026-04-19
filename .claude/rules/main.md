---
paths:
  - "main.go"
---

# Main Rules

- `main.go` is the composition root.
- Keep Cobra command construction and dependency wiring centralized here.
- Prefer wiring existing package APIs together rather than re-implementing business logic in `main.go`.
- Preserve environment and data directory behavior unless the task explicitly changes runtime configuration.
