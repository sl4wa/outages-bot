---
paths:
  - "internal/users/**/*.go"
---

# Users Rules

- `internal/users` owns user, address, and street entities plus validation, outage matching helpers, and user listing behavior.
- Preserve current user validation, subscription matching behavior, and repository interface expectations.
- Keep subscription conversation workflows, street search transitions, subscription display, unsubscribe, Telegram transport, outage API transport, and file-format concerns in their respective packages.
