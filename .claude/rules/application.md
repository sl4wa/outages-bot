---
paths:
  - "internal/application/**/*.go"
  - "internal/cmd/**/*.go"
  - "main.go"
---

# Application And Command Rules

- `internal/application` owns use-case logic such as outage fetch conversion, subscription save/search/unsubscribe, notifications, and user listing.
- `internal/cmd` owns Cobra command helpers plus the Telegram `BotRunner` conversation state machine, step tracking, reply-keyboard flow, cleanup ticker, and session expiry handling.
- `main.go` wires dependencies, environment inputs, and subcommand registration.
- Reuse existing application services and command helpers before adding new abstractions.
- Preserve `/start`, `/stop`, and `/subscription` behavior and keep the default bot session TTL at 30 minutes unless the task explicitly changes it.
- In notifier flow, unchanged outage snapshots short-circuit work, only `ErrRecipientUnavailable` removes a user, and user outage info is saved only after a successful send.
- Keep bot-facing copy in Ukrainian unless the task explicitly changes product text.
