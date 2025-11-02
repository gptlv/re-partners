# Pack Calculator

Pack Calculator is a small Go web application that determines how many bulk packs to ship for any order quantity. The application:

- Guarantees we never break packs.
- Minimises overshoot (extra items sent) and, for equal overshoot, minimises the pack count.
- Stores pack sizes in SQLite so they can be changed without recompiling.
- Exposes both a JSON API and an HTMX-powered HTML UI.

---

## Project Layout

```
.
├── cmd/server          # Application entrypoint
├── internal
│   ├── api             # HTTP handlers, view models
│   ├── app             # Domain service layer
│   ├── db              # SQLite opener + migration
│   ├── repository      # Data access (pack sizes)
│   ├── router          # HTTP mux wiring
│   └── view            # Template renderer abstraction
├── pkg/calculate       # Pure packing algorithm + tests
├── web/templates       # HTMX-friendly templates
├── Taskfile.yml        # Automation helpers
└── Dockerfile          # Multi-stage container build
```

---

## Requirements

- Go 1.25.3
- [Task](https://taskfile.dev/)
- Docker

---

## Getting Started

```bash
# One-shot setup (go mod tidy + download)
task tidy

# Run unit tests
task test

# Start the development server on :8080
task run
```

Open <http://localhost:8080> to use the HTMX UI. Enter an amount and the page will show the minimal pack combination, the requested amount, and the configured pack sizes.

---

## API Endpoints

| Method | Path         | Description                                      |
| ------ | ------------ | ------------------------------------------------ |
| GET    | `/`          | HTMX UI. Lists available pack sizes and form.    |
| POST   | `/ui/calc`   | HTMX form target. Returns HTML fragment.         |
| GET    | `/api/packs` | JSON array of pack sizes.                        |
| POST   | `/api/calc`  | JSON request `{ "amount": <int> }` → combination |

HTMX responses return `422 Unprocessable Entity` with a human-readable message if the order cannot be fulfilled with the current packs. JSON responses mirror that behaviour with the same status code and message.

---

## Running in Docker

```bash
task docker:run
# or manually
docker build -t pack-calculator .
docker run --rm -p 8080:8080 pack-calculator
```

The container image runs migrations on startup and listens on port `8080`.

## Taskfile Commands

| Command             | Description                       |
| ------------------- | --------------------------------- |
| `task run`          | Run the dev server                |
| `task test`         | Run `go test ./...`               |
| `task build`        | Build binary to `bin/server`      |
| `task fmt`          | Run `gofmt` across all packages   |
| `task tidy`         | `go mod tidy` + `go mod download` |
| `task docker:build` | Build the Docker image            |
| `task docker:run`   | Build + run the Docker container  |
