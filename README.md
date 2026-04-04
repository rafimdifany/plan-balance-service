# Plan Balance Service

Plan Balance Service is a robust Backend API built with Golang using the [Gin Web Framework](https://github.com/gin-gonic/gin). It strictly follows Domain-Driven Design (DDD) principles to ensure high maintainability, scalability, and code cleanliness.

## 🚀 Key Features

*   **Clean Architecture (DDD)**: Logical separation of concerns across Handler, Service, and Repository layers.
*   **Complete Authentication Module**:
    *   Standard Email/Password Sign-up & Login capabilities.
    *   Google OAuth identity linking and authentication.
    *   Atomic database transactions when creating credentials to avoid partial data states.
*   **Secure Session Management**:
    *   Issues short-lived JWT Access Tokens.
    *   Implements secure, rotating refresh tokens (hashed using SHA256 before database insertion).
*   **High Performance Database**: Utilizes **PostgreSQL** with `jackc/pgx/v5` connection pooling.
*   **Database Migrations**: Structured schema management via `golang-migrate`.
*   **Interactive Documentation**: Integrated Swagger UI mapping out endpoint schemas and allowing live testing interactively.
*   **Structured Logging**: Built-in high-performance logging utilizing Uber's `zap`.
*   **Strict Validations**: Robust request validation layer utilizing `go-playground/validator`.

## 📂 Directory Map

*   `cmd/api/` - The main application entry point. Handles setup and dependency injection wiring.
*   `internal/` - Private application logic.
    *   `config/`: Environment configuration management.
    *   `db/`: PostgreSQL instance and pool configuration.
    *   `handler/`: HTTP REST controllers.
    *   `service/`: Core business logic processing and orchestration.
    *   `repository/`: Database adapters (includes native `pgx.Tx` transaction support).
    *   `model/`: Core domain entities mirroring database schemas.
    *   `dto/`: Data Transfer Objects (Validation bindings & Response mappers).
    *   `middleware/`: Gin middleware (CORS, Request Tracing).
*   `pkg/` - Shared and abstracted public utilities.
    *   `logger/`: Zap logger integration.
    *   `utils/`: Assorted helpers (JWT, Hashing routines, Validation parsers).
*   `db/migrations/` - SQL directives (.up.sql / .down.sql) defining the DB architecture explicitly.
*   `docs/` - Auto-generated Swagger documentation bindings.

## ⚙️ Prerequisites

You need the following tools available in your environment:
*   [Go](https://go.dev/doc/install) 1.20+
*   [PostgreSQL](https://www.postgresql.org/download/) 14+
*   [Golang Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
*   [Swag CLI](https://github.com/swaggo/swag) (For generating API documentation)

## 🛠️ Running Locally

1.  **Clone the repository and fetch dependencies**:
    ```bash
    go mod tidy
    ```
2.  **Environment Setup**:
    Copy the sample environment values to define your local configurations.
    ```bash
    cp .env.example .env
    ```
    *(Remember to fill in your `DATABASE_URL`, `JWT_SECRET`, and `GOOGLE_CLIENT_ID` in `.env`)*
3.  **Perform Database Migration**:
    Create a database instance locally, and run the schema migration up command:
    ```bash
    make migrate-up
    # OR manually: migrate -path db/migrations -database "$DATABASE_URL" up
    ```
4.  **Boot the Server**:
    ```bash
    make run
    # OR manually: go run cmd/api/main.go
    ```
    The platform will spin up locally mapping to the configured port (default is `8080`).

## 📚 API Documentation (Swagger)

A beautiful, interactive representation of the endpoints is available after spinning up the application locally.
**Navigate your browser to:**
```
http://localhost:8080/swagger/index.html
```

💡 **Note for Contributors:** If you edit Handler structures or endpoints, quickly regenerate the swagger files before committing:
```bash
swag init -g cmd/api/main.go --output docs
```

## 🧪 Development Workflow

*   Build Binary: `make build`
*   Code Linting: `go vet ./...`
*   Run Unit Tests: `make test`

---
*Built with clean coding standards and modern web practices.*
