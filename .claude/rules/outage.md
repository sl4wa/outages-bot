---
paths:
  - "internal/outage/**/*.go"
---

# Outage Rules

- `internal/outage/outage` owns outage entities, constructors, validation, raw DTO conversion, equality helpers, and snapshot store interfaces.
- Use existing constructors such as `NewPeriod`, `NewAddress`, and `NewDescription` instead of assembling invalid state directly.
- The path filter `internal/outage/**` above scopes this rule to the outage app; for the schedule app see the schedule rules.
- Keep transport, Telegram, and file-format concerns out of this package.
- Preserve normalized outage matching behavior unless the task explicitly changes it.
