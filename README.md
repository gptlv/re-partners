# Pack Calculator

Pack Calculator is a small Go web application that determines how many bulk packs to ship for any order quantity. The application:

- Guarantees we never break packs.
- Minimises overshoot (extra items sent) and, for equal overshoot, minimises the pack count.
- Stores pack sizes in PostgreSQL
- Exposes a RESTful JSON API

---

## Project Layout

```
.
├── cmd/server          # Application entrypoint
├── internal
│   ├── api             # HTTP handlers, request/response models
│   ├── app             # Domain service layer
│   ├── repository      # Data access (pack sizes)
│   └── router          # HTTP mux wiring
├── pkg
│   ├── calculate       # Pure packing algorithm + tests
│   └── db              # PostgreSQL connection + migrations
├── migrations          # SQL migration files
├── deploy              # Docker Compose configuration
├── Taskfile.yml        # Automation helpers
└── Dockerfile          # Multi-stage container build
```

---

## Architecture

The application follows a clean architecture pattern:

- **cmd/server**: Application entry point and dependency wiring.
- **internal/api**: HTTP handlers, JSON serialization, and HTTP-specific logic.
- **internal/app**: Business logic and domain services.
- **internal/repository**: PostgreSQL data access layer.
- **internal/router**: HTTP route configuration.
- **pkg/calculate**: Pure algorithm for pack calculation (no external dependencies).
- **pkg/db**: Database connection and migration utilities.
- **migrations**: SQL migration files embedded and applied with goose.

---

## Requirements

- Go 1.25.3
- PostgreSQL 14+ (or Docker Compose)
- [Task](https://taskfile.dev/) (optional, for task automation)
- Docker & Docker Compose (for containerised deployment)

---

## Configuration

The application is configured via environment variables:

| Variable            | Required | Default   | Description              |
| ------------------- | -------- | --------- | ------------------------ |
| `POSTGRES_HOST`     | Yes      | -         | PostgreSQL host address  |
| `POSTGRES_PORT`     | No       | `5432`    | PostgreSQL port          |
| `POSTGRES_USER`     | Yes      | -         | PostgreSQL username      |
| `POSTGRES_PASSWORD` | No       | -         | PostgreSQL password      |
| `POSTGRES_DB`       | Yes      | -         | PostgreSQL database name |
| `POSTGRES_SSLMODE`  | No       | `disable` | PostgreSQL SSL mode      |

Default pack sizes (250, 500, 1000, 2000, 5000) are seeded automatically during migration.

---

## Getting Started

```bash
# From the repository root
cp .env.example .env

docker compose --env-file .env -f deploy/docker-compose.yml up
```

The API will be available at `http://localhost:8080`.

### Testing the API

```bash
# List configured pack sizes
curl http://localhost:8080/api/packs

# Calculate packs for an order
curl -X POST http://localhost:8080/api/calc \
  -H "Content-Type: application/json" \
  -d '{"amount": 251}'

# Add a new pack size
curl -X POST http://localhost:8080/api/sizes \
  -H "Content-Type: application/json" \
  -d '{"size": 3000}'

# Delete a pack size (replace {id} with the actual ID)
curl -X DELETE http://localhost:8080/api/sizes/{id}
```

---

## API Endpoints

### Pack Calculation

| Method | Path         | Description                        | Request Body        | Response                                                  |
| ------ | ------------ | ---------------------------------- | ------------------- | --------------------------------------------------------- |
| GET    | `/api/packs` | List all configured pack sizes     | -                   | `{ "packs": [{"id": 1, "size": 250}, ...] }`              |
| POST   | `/api/calc`  | Calculate optimal pack combination | `{ "amount": 251 }` | `{ "amount": 251, "packs": [{"size": 500, "count": 1}] }` |

### Pack Size Management

| Method | Path              | Description         | Request Body       | Response                    |
| ------ | ----------------- | ------------------- | ------------------ | --------------------------- |
| POST   | `/api/sizes`      | Add a new pack size | `{ "size": 3000 }` | `{ "id": 6, "size": 3000 }` |
| DELETE | `/api/sizes/{id}` | Remove a pack size  | -                  | `204 No Content`            |

### Example Requests

```bash
curl -X POST http://localhost:8080/api/calc \
  -H "Content-Type: application/json" \
  -d '{"amount": 12001}'

# Response:
{
  "amount": 12001,
  "packs": [
    {"size": 250, "count": 1},
    {"size": 2000, "count": 1},
    {"size": 5000, "count": 2}
  ]
}
```

```bash
curl -X POST http://localhost:8080/api/sizes \
  -H "Content-Type: application/json" \
  -d '{"size": 750}'

# Response:
{
  "id": 6,
  "size": 750
}
```

```bash
curl -X DELETE http://localhost:8080/api/sizes/6

# Response: 204 No Content
```

## Taskfile Commands

| Command             | Description                                    |
| ------------------- | ---------------------------------------------- |
| `task run`          | Run the development server (requires env vars) |
| `task dev`          | Start PostgreSQL via Docker Compose and run it |
| `task test`         | Run `go test ./...`                            |
| `task build`        | Build the binary to `bin/server`               |
| `task fmt`          | Run `go fmt ./...`                             |
| `task tidy`         | Run `go mod tidy` and `go mod download`        |
| `task db:up`        | Start PostgreSQL via Docker Compose            |
| `task db:down`      | Stop PostgreSQL                                |
| `task db:reset`     | Reset PostgreSQL (down, remove volumes, up)    |
| `task docker:build` | Build the Docker image                         |
| `task docker:run`   | Run the full stack via Docker Compose          |
| `task docker:down`  | Stop all containers                            |
