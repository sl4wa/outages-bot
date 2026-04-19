---
paths:
  - "internal/users/**/*.go"
---

# Users Rules

- `internal/users` owns subscription, street search, subscription display, unsubscribe, and user listing flows.
- Keep user-facing bot copy in Ukrainian unless the task explicitly changes product text.
- Preserve current subscription matching behavior and repository interface expectations.
- Do not move Telegram transport, outage API transport, or file-format concerns into this package.
