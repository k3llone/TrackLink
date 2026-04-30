# TrackLink

Backend foundation for TrackLink.

## Local run (one command)

Run API + PostgreSQL + Redis:

`docker compose up --build`

After startup:

- API health: `http://localhost:8080/health`
- API base group: `http://localhost:8080/api/v1`

## Environment

For local `go run`, copy and adjust:

`backend/.env.example` -> `backend/.env`

Required variables:

- `HTTP_ADDR`
- `PUBLIC_URL`
- `DATABASE_URL` (legacy fallback: `POSTGRES_DSN`)
- `REDIS_ADDR` (legacy fallback: `REDIS_DSN`)
- `SESSION_SECRET`

## Run without Docker API container

You can run only infra:

`docker compose up -d postgres redis`

Then run backend app from `backend` directory:

`go run ./cmd/api`

## Migrations

Migrations directory: `backend/migrations`

Recommended CLI: [golang-migrate](https://github.com/golang-migrate/migrate)

Example command from `backend` directory:

`migrate -path ./migrations -database "$DATABASE_URL" up`