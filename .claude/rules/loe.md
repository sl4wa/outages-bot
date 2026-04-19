---
paths:
  - "internal/loe/**/*.go"
---

# LOE Rules

- Keep outage API code transport-focused.
- Preserve HTTP cache behavior, response normalization, and array-or-string building parsing.
- Preserve deduplication by street, buildings, start, and end, with later duplicates replacing earlier rows.
- Do not move subscription, notifier, or Telegram conversation logic into this package.
