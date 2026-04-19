---
paths:
  - "internal/telegram/**/*.go"
---

# Telegram Rules

- `internal/telegram` owns Telegram transport plus the bot runner conversation state machine.
- Preserve reply-keyboard flow, cleanup behavior, and the default 30-minute bot session TTL unless the task explicitly changes product behavior.
- Keep subscription logic, outage normalization, and persistence formats in their respective packages.
