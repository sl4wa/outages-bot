---
paths:
  - "internal/persistence/**/*.go"
---

# Persistence Rules

- Keep repositories data-focused; do not move business rules or orchestration into persistence code.
- Preserve file-backed storage contracts:
  - one user per `*.yml` file under the configured `DATA_DIR/users` directory, with default `data/users`
  - CSV-backed street and outage snapshot formats
- Preserve atomic temp-file-plus-rename writes for both user files and outage snapshot files.
- When persistence behavior changes, review affected callers for missing-file, ordering, and snapshot assumptions.
