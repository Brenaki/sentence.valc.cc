# Sentence API (Go)

Quote-of-the-day API following Clean Architecture, SOLID and TDD.

## Layout

```
cmd/api                      entrypoint (composition root)
internal/domain              entities + ports (interfaces) + sentinel errors
internal/usecase             business rules (GetQuoteOfTheDay, ReactToQuote)
internal/infra/http          delivery: handlers, router, CORS
internal/infra/provider/ninja  api-ninjas client (QuoteProvider)
internal/infra/repository/mysql  MySQL adapter (QuoteRepository)
internal/config              env config
migrations                   schema (auto-loaded by docker mysql initdb)
```

Dependencies point inward: `infra` and `usecase` depend on `domain` interfaces;
`domain` depends on nothing. Use cases receive their dependencies via
constructors (dependency inversion), so each layer is unit-tested in isolation.

## Endpoints

- `GET /quote-of-the-day` — returns today's quote. Looks in MySQL for a row
  created today; if absent, fetches from api-ninjas, persists it, returns it.
- `POST /quotes/{id}/reactions` — body `{"reaction": 0|1}` (0 = dislike, 1 = like).
  Increments the matching counter. `204` on success, `400` invalid, `404` unknown id.
- `GET /healthz` — liveness.

## Run

```bash
# 1. config: .env must contain API_KEY_NINJA=...
# 2. database
docker compose up -d
# 3. api (mise provides go 1.25)
set -a; source .env; set +a
go run ./cmd/api
```

Env vars (defaults in `internal/config`): `HTTP_ADDR` (`:8080`),
`MYSQL_DSN` (`root:root@tcp(127.0.0.1:3306)/sentence?parseTime=true...`),
`API_KEY_NINJA` (required), `ALLOW_ORIGINS` (`*`).

## Test

```bash
go test ./...
```

Use-case, handler, provider (httptest) and repository (sqlmock) layers are all
covered without needing a live database or network.
