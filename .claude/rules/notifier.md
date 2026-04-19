---
paths:
  - "internal/notifier/**/*.go"
---

# Notifier Rules

- Keep this package focused on notification orchestration.
- Preserve current semantics: unchanged outage snapshots short-circuit work, only `ErrRecipientUnavailable` removes a user, and user outage info is saved only after a successful send.
- Do not move outage transport, Telegram transport, or persistence format logic into this package.
