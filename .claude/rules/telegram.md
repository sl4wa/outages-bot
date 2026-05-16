---
paths:
  - "internal/outage/telegram/**/*.go"
---

# Telegram Rules

- `internal/outage/telegram` adapts Telegram updates to `internal/outage/subscription` commands and renders `subscription.Response` values back to Telegram messages.
- Keep command mapping, reply-keyboard construction, and Telegram send/logging behavior in this package.
- Keep subscription workflow state, street search transitions, save/unsubscribe/display operations, outage normalization, and persistence formats in their respective packages.
