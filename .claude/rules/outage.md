---
paths:
  - "internal/outage/**/*.go"
---

# Outage Rules

- `internal/outage` owns outage entities, constructors, validation, raw DTO conversion, equality helpers, and snapshot store interfaces.
- Use existing constructors such as `NewPeriod`, `NewAddress`, and `NewDescription` instead of assembling invalid state directly.
- Keep transport, Telegram, and file-format concerns out of this package.
- Preserve normalized outage matching behavior unless the task explicitly changes it.
