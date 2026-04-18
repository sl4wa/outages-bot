---
paths:
  - "internal/repository/**/*.go"
---

# Repository Rules

- Keep repositories data-focused; do not move domain rules or orchestration into repository code.
- Preserve compatibility with the current file-backed storage layout under `data/`, especially one YAML file per user under `data/users`.
- Keep CSV and on-disk formats stable unless the task explicitly changes persistence contracts.
- `FileUserRepository.Save` must keep the temp-file-plus-rename atomic write pattern and clean up the temp file if rename fails.
- When persistence behavior changes, review affected callers for ordering, missing-file, and snapshot assumptions.
