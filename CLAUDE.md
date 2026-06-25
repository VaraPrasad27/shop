# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Go HTTP backend for a shop. Currently exposes `GET /` (returns `shop api`) and `GET /products` (returns the product catalogue from PostgreSQL with optional `?limit=` and `?offset=` pagination). Frontend is not yet present in this repo.

## Build / run / test

All commands run from the `server/` directory (it's a standalone Go module — module path is `github.com/VaraPrasad27/shop/server`).

```sh
cd server

# Build
go build -o ./tmp/main ./cmd/.

# Run (loads DATABASE_URL from .env via godotenv; PORT defaults to 8080)
go run ./cmd/.

# Live-reload during development (uses the .air.toml config in this dir)
air

# Vet
go vet ./...

# Run a single test (no tests exist yet — create _test.go files next to the code)
go test ./internal/handlers/ -run TestGetAllProductsHandler

# Tidy
go mod tidy
```

## Layout

```
server/
├── cmd/main.go                    # entry point: load config → connect DB → mount routes → run server with graceful shutdown
├── .air.toml                      # live-reload config (builds ./tmp/main from ./cmd/)
├── .env                           # DATABASE_URL + PG creds (gitignored, copy locally)
├── go.mod / go.sum                # module: github.com/VaraPrasad27/shop/server
├── internal/
│   ├── config/                    # env loading (godotenv), validates DATABASE_URL, reads optional PORT
│   ├── db/                        # pgxpool.New + Ping(ctx) so startup fails loudly
│   ├── models/                    # Product struct (db: + json: tags)
│   ├── repositories/              # SQL via pgxpool; uses pgx.CollectRows + RowToStructByName
│   ├── services/                  # thin pass-through to repositories; place for future business logic
│   ├── handlers/                  # http.HandlerFunc factories returning handlers
│   └── routes/                    # New(pool) *chi.Mux — route mounting only
└── seed/
    ├── schema.sql                 # `products` table DDL (does NOT enable pgcrypto — see notes)
    └── seed.sql                   # sample rows (id, name, description, price_cents, currency, image_url, stock)
```

## Request flow (end-to-end)

`main.go` → `config.LoadConfig()` (reads `.env`, validates `DATABASE_URL`, reads `PORT`) → `db.Connect(ctx, url)` (creates pool + pings) → `routes.New(pool)` returns a `chi.Mux` → `&http.Server{...}.ListenAndServe()` in a goroutine → main blocks on `signal.Notify(Interrupt, SIGTERM)` → `srv.Shutdown(ctx)` on signal with a 15s grace period.

For `GET /products`: handler parses `?limit=` and `?offset=` (defaults 25 and 0; max 100) → `services.GetAllProducts(ctx, pool, limit, offset)` → `repositories.GetAllProducts(ctx, pool, limit, offset)` → SQL `SELECT id, name, description, price_cents, currency, image_url, stock, created_at FROM products ORDER BY id LIMIT $1 OFFSET $2` → `pgx.CollectRows(rows, pgx.RowToStructByName[models.Product])` → JSON-encoded response with `Content-Type: application/json`.

## Architectural notes / things to know

- **Money is `int64` cents, not `float64` dollars.** `models.Product.PriceCents` maps to the `price_cents INTEGER` column. JSON tag is `price_cents`. Do not introduce float-based money fields.
- **Layering: handlers own HTTP, services own business logic, repositories own SQL.** Services currently pass through to repositories — that's intentional, it's where non-trivial logic (filtering, validation, transformations) goes later. Don't add SQL to services, don't bypass services from handlers.
- **Handlers are factory functions returning `http.HandlerFunc`** (`handlers.GetAllProductsHandler(pool)`), not methods on a struct. New endpoints follow the same pattern.
- **`context.Context` flows through every layer.** Handlers pass `r.Context()` down; never use `context.Background()` inside the request path.
- **Router is `go-chi/chi` v1.5.5**, not the stdlib mux. Use `chi.URLParam` for path params once they appear. `routes.New(pool)` returns `*chi.Mux`.
- **DB driver is `pgx/v5` with `pgxpool`** — use the pool, don't open connections per-request.
- **Model-to-row mapping uses `pgx.RowToStructByName`** with `db:` tags on `models.Product`. Don't add positional `rows.Scan(&product.X, &product.Y, ...)` — that's the old pattern.
- **Errors are logged server-side, never leaked to clients.** The handler returns a generic `"internal server error"` body and logs the real error via `log.Printf`. Don't `http.Error(w, err.Error(), ...)`.
- **CORS, auth, request logging middleware are not set up yet.** Add them in `routes.New` when needed.
- **No tests yet.** The handler is straightforward to test with `httptest.NewRecorder`; the repo can be tested with a real Postgres via `pgxpool` or with `pgxmock`.

## Environment

`.env` (server-local, gitignored) must define `DATABASE_URL`. Optional: `PORT` (defaults to `8080`). The file in this repo uses PostgreSQL on `localhost:5432`, db `shop`, user `vara`. `config.LoadConfig` always calls `godotenv.Load()` (no-op if absent); process env overrides `.env`.

To bootstrap a local DB:

```sh
psql -U vara -d shop -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"  # one-time, schema needs it
psql -U vara -d shop -f server/seed/schema.sql
psql -U vara -d shop -f server/seed/seed.sql
```

The schema uses `gen_random_uuid()` from `pgcrypto`; that extension is not enabled by `schema.sql` itself, so run the `CREATE EXTENSION` command first or the table creation will fail on a fresh DB. Seed is **not idempotent** — re-running `seed.sql` on an already-seeded DB will fail with duplicate-row errors.