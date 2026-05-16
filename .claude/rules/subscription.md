---
paths:
  - "internal/subscription/**/*.go"
---

# Subscription Rules

- `internal/subscription` owns the subscription conversation workflow and application-level subscription operations.
- Keep pending in-memory workflow state, command handling, street search transitions, subscription save, unsubscribe, and current-subscription display behavior in this package.
- Preserve `/start` recovery when current subscription lookup fails: prompt for street entry and attach the lookup error so adapters can log it.
- Treat pending conversation expiry as workflow-owned lazy TTL behavior, defaulting to 30 minutes; do not add a `BotRunner` cleanup goroutine unless a future task explicitly changes that tradeoff.
- Return structured `Response` values (including all Ukrainian text) for adapters to render; keep Telegram bot API calls, reply-keyboard construction, and transport details in `internal/telegram`.
- Use `internal/users` entities, validation, repositories, and matching behavior instead of duplicating user/address/street rules here.
