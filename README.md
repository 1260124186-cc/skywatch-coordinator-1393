# Skywatch Coordinator

Skywatch Coordinator is a small HTTP service used by field-science teams to coordinate night-sky observation campaigns. A coordinator opens a campaign, stations operate observation shifts, and reviewers release a coherent observation set once every submitted frame has been assessed.

## Run

```bash
go run ./cmd/skywatch
```

The service listens on `127.0.0.1:18087` by default. Set `SKYWATCH_ADDR` to use another address.

## Useful endpoints

- `GET /health`
- `POST /v1/campaigns`
- `POST /v1/shifts`
- `POST /v1/observations`
- `POST /v1/shifts/{id}/close`
- `POST /v1/campaigns/{id}/release`
- `GET /v1/campaigns/{id}/summary`

## Test

```bash
go test ./...
```

The in-memory service starts with a small public demonstration campaign so a running instance has an immediately inspectable workflow.
