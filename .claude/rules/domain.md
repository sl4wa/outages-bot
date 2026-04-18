---
paths:
  - "internal/domain/**/*.go"
---

# Domain Rules

- Keep outage, address, period, description, and user rules in domain types.
- Use constructors and validation helpers so invalid domain state is rejected at the boundary.
- Do not move repository, Telegram, CLI, or file-format concerns into domain code.
- Preserve existing matching and subscription behavior unless the task explicitly changes it.
