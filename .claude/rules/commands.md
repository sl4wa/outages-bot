# Run And Test Commands

- `make test` runs the full test suite.
- `make build` builds `bin/outages-bot`.
- `go test ./...` runs all Go tests.
- Useful focused suites:
  - `go test ./internal/domain/...`
  - `go test ./internal/application/...`
  - `go test ./internal/repository/...`
  - `go test ./internal/client/...`
  - `go test ./internal/cmd/... -run TestName`
- Useful local commands:
  - `go run . bot`
  - `go run . notifier`
  - `go run . notifier --interval=60s`
  - `go run . outages`
  - `go run . users`
- Runtime commands may require `TELEGRAM_BOT_TOKEN`, `OUTAGE_API_URL`, and `DATA_DIR`.
- `DATA_DIR` defaults to `data` when unset.
