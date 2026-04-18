# CLAUDE.md

outages-bot is a Go Telegram bot for outage subscriptions and notifications.

Repo-specific Claude Code guidance lives in `.claude/rules/`.

- Start with `.claude/rules/core.md` and `.claude/rules/commands.md`.
- Claude Code also loads path-scoped rules from `.claude/rules/` when working with matching files.
- Treat notifier semantics, bot session flow, file persistence, and outage normalization rules as intentional behavior unless the task explicitly changes them.
