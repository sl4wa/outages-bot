---
paths:
  - "**/*_test.go"
  - "internal/integration/**"
---

# Test Rules

- Behavior changes should get targeted Go tests in the nearest affected package.
- Prefer narrow test updates over broad fixture or snapshot rewrites.
- Keep integration-test contracts stable unless the task explicitly changes external behavior.
- When changing notifier, repository, bot TTL/session flow, or outage parsing and deduplication, review both focused unit tests and integration coverage.
