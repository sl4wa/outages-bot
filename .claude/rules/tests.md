---
paths:
  - "**/*_test.go"
  - "test/integration/**"
---

# Test Rules

- Behavior changes should get targeted Go tests in the nearest affected package.
- Docs-only or guidance-only changes should use read-only verification such as `rg`, `git diff`, or line-count review instead of Go tests.
- Prefer narrow test updates over broad fixture or snapshot rewrites.
- Keep integration-test contracts stable unless the task explicitly changes external behavior.
- When changing notifier behavior, review both unit and integration coverage for snapshot short-circuiting, recipient removal, and saved outage info ordering.
- When changing subscription or Telegram adapter behavior, review tests around subscription workflow state, command mapping, reply-keyboard handling, and Ukrainian user-facing copy.
- For subscription pending-state expiry behavior, inject or control `Workflow.now` in tests instead of using sleeps or real-time waits.
- When changing persistence, review tests for TOML user storage, CSV street and outage snapshot formats, and atomic temp-file-plus-rename writes.
- When changing outage parsing or `loe` behavior, review normalization, building parsing, deduplication, and downstream matching coverage.
