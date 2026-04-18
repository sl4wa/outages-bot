---
paths:
  - "internal/client/**/*.go"
---

# Client Rules

- Keep outage API and Telegram code transport-focused.
- Do not move subscription, outage matching, or notification decision logic into clients.
- In `internal/client/outageapi`, preserve comment normalization, building parsing from array-or-string input, and deduplication by street/buildings/start/end with later duplicates replacing earlier rows.
- Preserve existing request, response, and formatting behavior unless the task explicitly changes an external contract.
- When client behavior changes, verify the nearest application or integration tests that depend on it.
