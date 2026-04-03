# Plan Balance Service - Backend

This project follows a structured architecture for better maintainability and scalability.

## Directory Structure

- `cmd/api/`: Application entry point.
- `internal/`: Private application and business logic.
    - `config/`: Configuration and environment variable loading.
    - `handler/`: HTTP handlers (Controllers).
    - `service/`: Core business logic.
    - `repository/`: Database access layer.
    - `model/`: Domain entities and database structs.
    - `dto/`: Data Transfer Objects for request/response.
    - `db/`: Database related code (connection).
    - `middleware/`: Gin middleware (auth, logging, etc.).
- `pkg/`: Public utility packages.
    - `logger/`: Zap logger integration.
    - `utils/`: Common helpers (security, validator, etc.).
- `db/`:
    - `migrations/`: SQL migration files.
    - `seed/`: Optional dummy data scripts.
- `scripts/`: Useful automation scripts (deploy, migrate).

## How to Run

1. Copy `.env.example` to `.env`.
2. Update the configuration values.
3. Run the application:
   ```bash
   make run
   ```

## Development

- **Build**: `make build`
- **Lint**: `go vet ./...`
- **Test**: `make test`
