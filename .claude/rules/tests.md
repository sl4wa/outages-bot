---
paths:
  - "**/*_test.go"
  - "internal/integration/**"
---

# Test Rules

- Behavior changes should get targeted Go tests in the nearest affected package.
- Prefer narrow test updates over broad fixture or snapshot rewrites.
- Keep integration-test contracts stable unless the task explicitly changes external behavior.
- When changing notifier behavior, review both unit and integration coverage for snapshot short-circuiting, recipient removal, and saved outage info ordering.
- When changing Telegram or user flow behavior, review tests around bot TTL/session flow, reply-keyboard handling, and Ukrainian user-facing copy.
- When changing persistence, review tests for YAML user storage, CSV street and outage snapshot formats, and atomic temp-file-plus-rename writes.
- When changing outage parsing or `loe` behavior, review normalization, building parsing, deduplication, and downstream matching coverage.
