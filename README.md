# Chirpy

A Twitter-style microblogging REST API written in Go, built as part of the Boot.dev "Learn HTTP Servers in Go" course. Users can sign up, log in, post short messages called chirps, and manage them through a JSON API backed by PostgreSQL.

## Features

- User registration and login with Argon2id password hashing
- JWT access tokens (1 hour) and revocable refresh tokens (60 days)
- Create, list, fetch, and delete chirps (users can only delete their own)
- Filter chirps by author and sort them by creation date
- 140-character limit and profanity filtering on chirps
- Update your own email and password
- Webhook endpoint for upgrading users to Chirpy Red, protected by an API key
- Request metrics and a dev-only reset endpoint
- Type-safe database access with sqlc

## Prerequisites

- [Go](https://go.dev/dl/) 1.26 or later (see `go.mod`)
- [PostgreSQL](https://www.postgresql.org/download/)
- [goose](https://github.com/pressly/goose) for database migrations:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/Ravinder2102/chirpy.git
cd chirpy
```

### 2. Create the database

```bash
psql -U postgres -c "CREATE DATABASE chirpy;"
```

### 3. Run the migrations

```bash
goose -dir sql/schema postgres "postgres://postgres:YOUR_PASSWORD@localhost:5432/chirpy?sslmode=disable" up
```

### 4. Create a `.env` file

In the project root:

```env
DB_URL="postgres://postgres:YOUR_PASSWORD@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="your-secret"
POLKA_KEY="your-polka-api-key"
```

- `DB_URL`: PostgreSQL connection string
- `PLATFORM`: Set to `dev` to allow `POST /admin/reset`
- `JWT_SECRET`: Secret used to sign access tokens. Generate one with `openssl rand -base64 64`
- `POLKA_KEY`: API key that the payments webhook must send

All four variables are required. The server exits on startup if any is missing.

### 5. Run the server

```bash
go build -o chirpy && ./chirpy
```

The server listens on `http://localhost:8080`.

## API Reference

Endpoints that need an access token expect the header `Authorization: Bearer <access_token>`.

### Health and admin

| Method | Endpoint         | Description                                                                                         |
| ------ | ---------------- | --------------------------------------------------------------------------------------------------- |
| GET    | `/api/healthz`   | Readiness check, returns `OK`                                                                       |
| GET    | `/admin/metrics` | HTML page showing how many times `/app/` has been visited                                           |
| POST   | `/admin/reset`   | Reset the hit counter and delete all users. Only works when `PLATFORM=dev`, otherwise returns `403` |
| GET    | `/app/`          | Serves the static files in the project root                                                         |

### Users and authentication

| Method | Endpoint       | Auth          | Description                                       |
| ------ | -------------- | ------------- | ------------------------------------------------- |
| POST   | `/api/users`   | None          | Create a user                                     |
| PUT    | `/api/users`   | Access token  | Update your email and password                    |
| POST   | `/api/login`   | None          | Log in, returns an access token and refresh token |
| POST   | `/api/refresh` | Refresh token | Get a new access token                            |
| POST   | `/api/revoke`  | Refresh token | Revoke a refresh token (`204`)                    |

For `/api/refresh` and `/api/revoke`, send the refresh token in the `Authorization: Bearer` header.

### Chirps

| Method | Endpoint                | Auth         | Description                         |
| ------ | ----------------------- | ------------ | ----------------------------------- |
| POST   | `/api/chirps`           | Access token | Create a chirp (max 140 characters) |
| GET    | `/api/chirps`           | None         | List chirps                         |
| GET    | `/api/chirps/{chirpID}` | None         | Get one chirp                       |
| DELETE | `/api/chirps/{chirpID}` | Access token | Delete your own chirp (`204`)       |

`GET /api/chirps` accepts optional query parameters:

- `author_id=<uuid>` returns only that user's chirps
- `sort=desc` returns newest first (the default is oldest first)

Chirps containing certain banned words have them replaced with `****`.

### Webhooks

| Method | Endpoint              | Auth    | Description                                             |
| ------ | --------------------- | ------- | ------------------------------------------------------- |
| POST   | `/api/polka/webhooks` | API key | Upgrade a user to Chirpy Red on a `user.upgraded` event |

Send the key as `Authorization: ApiKey <POLKA_KEY>`. Events other than `user.upgraded` are ignored with a `204`.

## Examples

Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

Log in:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

Response:

```json
{
  "id": "…",
  "created_at": "…",
  "updated_at": "…",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```

Post a chirp:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello, Chirpy!"}'
```

List a user's chirps, newest first:

```bash
curl "http://localhost:8080/api/chirps?author_id=<user_id>&sort=desc"
```

## Running Tests

```bash
go test ./...
```

## Built With

- [Go](https://go.dev) standard library (`net/http`)
- [PostgreSQL](https://www.postgresql.org) with [lib/pq](https://github.com/lib/pq)
- [sqlc](https://sqlc.dev) and [goose](https://github.com/pressly/goose)
- [argon2id](https://github.com/alexedwards/argon2id), [golang-jwt](https://github.com/golang-jwt/jwt), [google/uuid](https://github.com/google/uuid), [godotenv](https://github.com/joho/godotenv)