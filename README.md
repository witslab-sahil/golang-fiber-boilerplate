# Fiber Microservice Boilerplate with YugabyteDB

A production-ready Golang microservice boilerplate built with Fiber framework, featuring clean architecture, Docker support, comprehensive testing, workflow orchestration, and policy-based authorization, using YugabyteDB as the database.

## 📋 Table of Contents

- [Documentation Summary](#documentation-summary)
- [Core Features](#🚀-core-features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Start](#🏃-quick-start)
- [Available Endpoints](#📡-available-endpoints)
- [Testing](#testing)
- [Development](#development)
- [Configuration](#configuration)
- [API Examples](#api-examples)
- [YugabyteDB Features](#yugabytedb-features)
- [Architecture](#architecture)
- [Observability](#observability)
- [Temporal Workflow Engine](#🔄-temporal-workflow-engine)
- [OPA Authorization](#🛡️-opa-authorization)
- [Troubleshooting](#🔍-troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Documentation Summary

This boilerplate provides a complete microservice foundation with:

- **Fiber Framework** - Fast, Express-inspired web framework written in Go
- **Clean Architecture** with separated layers (handlers, services, repositories)
- **YugabyteDB** distributed SQL database (PostgreSQL-compatible)
- **Temporal** workflow engine for orchestrating long-running processes
- **OPA** (Open Policy Agent) for fine-grained authorization
- **OpenTelemetry** for distributed tracing and observability
- **Docker** containerization with hot reload support
- **Comprehensive testing** with unit tests and mocks

## 🚀 Core Features

### Architecture & Development

- **Fiber Framework**: High-performance web framework with Express-like API
- **Clean Architecture**: Separation of concerns with distinct layers (handlers, services, repositories)
- **RESTful API**: Built with Fiber's powerful routing
- **Hot Reload**: Development with hot reload support using Air
- **Testing**: Comprehensive unit tests with mocking

### Data & Storage

- **Database**: YugabyteDB (PostgreSQL-compatible distributed SQL database) with GORM ORM
- **Migrations**: Database migration support

### Security & Authorization

- **JWT Authentication**: Secure token-based authentication
- **OPA Integration**: Policy-based authorization with role-based access control (RBAC)
- **Middleware**: Request ID, CORS, logging, recovery, and tracing middleware

### Workflow & Orchestration

- **Temporal Integration**: Workflow engine for long-running, reliable processes
- **Sample Workflows**: User onboarding workflow included
- **Worker Service**: Dedicated service for processing workflow tasks

### Observability & Monitoring

- **OpenTelemetry**: Distributed tracing and metrics collection
- **Jaeger Integration**: Trace visualization and analysis
- **Structured Logging**: JSON logging with Logrus

### Infrastructure

- **Docker**: Full Docker and Docker Compose support
- **Configuration**: Environment-based configuration
- **Service Discovery**: Multiple services orchestrated via Docker Compose

## Project Structure

```
fiber-boilerplate/
├── cmd/
│   ├── api/
│   │   └── main.go           # Application entry point
│   └── worker/
│       └── main.go           # Temporal worker entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── handlers/             # HTTP handlers (controllers)
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Domain models and DTOs
│   ├── repository/           # Data access layer
│   ├── service/              # Business logic layer
│   ├── temporal/             # Temporal workflows and activities
│   │   ├── activities/
│   │   ├── worker/
│   │   └── workflows/
│   └── opa/                  # OPA policies and middleware
│       ├── middleware/
│       └── policies/
├── pkg/
│   ├── database/             # Database connection
│   ├── logger/               # Logger interface and implementation
│   ├── telemetry/            # OpenTelemetry setup
│   ├── temporal/             # Temporal client
│   └── utils/                # Utility functions
├── migrations/               # Database migrations
├── scripts/                  # Utility scripts
├── temporal/                 # Temporal configuration
├── docker-compose.yml        # Docker compose configuration
├── docker-compose.dev.yml    # Development override
├── Dockerfile                # Docker image definition
├── Dockerfile.dev            # Development Docker image
├── Makefile                  # Common tasks
├── .air.toml                 # Air configuration for hot reload
├── .env.example              # Example environment variables
└── go.mod                    # Go module file
```

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- YugabyteDB (if running locally)
- Make (optional, for using Makefile commands)

## 🏃 Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:

```bash
git clone <repository-url>
cd fiber-boilerplate
```

2. Copy the environment file:

```bash
cp .env.example .env
```

3. Start the services:

**For Production-like environment:**

```bash
docker compose up -d
# or
make docker-up
```

**For Development with hot reload:**

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
# or
make dev
```

4. The API will be available at `http://localhost:8080`

5. Access the web interfaces:
   - **API**: `http://localhost:8080`
   - **YB-Master UI**: `http://localhost:7001` (YugabyteDB management)
   - **YB-TServer UI**: `http://localhost:9001` (YugabyteDB monitoring)

### Docker Compose Configuration

This boilerplate includes two docker-compose files:

#### 1. `docker-compose.yml` (Core Services)

The primary configuration file containing essential services:

- **Application**: Main API service built with Fiber
- **Database**: YugabyteDB (PostgreSQL-compatible distributed SQL)

Note: Optional services (Temporal, Jaeger, OPA) have been removed from the main docker-compose.yml to keep it minimal. The application is configured with these features disabled by default.

#### 2. `docker-compose.dev.yml` (Development Override)

Development-specific overrides for hot reloading:

- Uses `Dockerfile.dev` with Air for live code reloading
- Mounts source code as volume
- Enables automatic restart on code changes

### Service Architecture

```
┌─────────────────┐     ┌──────────────┐
│   Application   │────▶│  YugabyteDB  │
│   (Port 8080)   │     │  (Port 5433) │
└─────────────────┘     └──────────────┘
```

The minimal setup includes only the essential services. Optional services (Temporal, Jaeger, OPA) can be enabled by modifying the environment variables in docker-compose.yml and deploying the required services separately.

### Docker Commands

```bash
# Start all services
docker compose up -d

# Start with development hot reload
docker compose -f docker-compose.yml -f docker-compose.dev.yml up

# View logs
docker compose logs -f [service-name]

# Stop all services
docker compose down

# Stop and remove volumes (clean slate)
docker compose down -v

# Run database migrations
docker compose run --rm migrate

# Run tests in Docker
docker compose --profile test run --rm test

# Access service shells
docker compose exec app sh
docker compose exec db bash
```

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Set up YugabyteDB and update the `.env` file with your database credentials (YugabyteDB uses port 5433 for YSQL)

3. Run migrations:

```bash
make migrate-up
```

4. Run the application:

```bash
make run
# or
go run cmd/api/main.go
```

## 📡 Available Endpoints

### Health Check

- `GET /health` - Service health check

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token

### User Management

- `GET /api/v1/users` - Get all users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Workflow Management (Temporal)

- `POST /api/v1/workflows/user-onboarding` - Start user onboarding workflow
- `GET /api/v1/workflows/:id/status` - Check workflow status
- `GET /api/v1/workflows` - List all workflows

## Testing

### Run all tests

```bash
make test
```

### Run tests with coverage

```bash
make test-coverage
```

### Run tests in Docker

```bash
make docker-test
```

## Development

### Common Make commands

```bash
make help           # Show all available commands
make run            # Run the application
make build          # Build the application
make test           # Run tests
make lint           # Run linter
make fmt            # Format code
make docker-build   # Build Docker image
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make migrate-up     # Run migrations
make clean          # Clean build artifacts
```

## Configuration

The application uses environment variables for configuration. See `.env.example` for available options:

- `ENVIRONMENT` - Application environment (development/production)
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - YugabyteDB connection string (PostgreSQL-compatible, default port 5433)
- `LOG_LEVEL` - Logging level (debug/info/warn/error)
- `JWT_SECRET` - Secret key for JWT tokens

### OpenTelemetry Configuration

- `OTEL_ENABLED` - Enable/disable OpenTelemetry (true/false)
- `OTEL_SERVICE_NAME` - Service name for tracing (default: fiber-boilerplate)
- `OTEL_EXPORTER_TYPE` - Exporter type (jaeger/otlp)
- `OTEL_ENDPOINT` - Trace collector endpoint (default: http://localhost:14268/api/traces)

## API Examples

### Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "secure123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Get All Users

```bash
curl http://localhost:8080/api/v1/users?page=1&page_size=10
```

### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Why Fiber?

This boilerplate uses Fiber instead of Gin for several reasons:

### Performance
- **Fiber** is built on top of Fasthttp, which is known for its exceptional performance
- Benchmarks show Fiber can handle more requests per second with lower latency
- Memory efficient with zero memory allocations in hot paths

### Express-like API
- Familiar API for developers coming from Node.js/Express background
- Simple and intuitive routing syntax
- Built-in middleware support

### Features
- Built-in WebSocket support
- Built-in rate limiter
- Built-in compression
- Built-in session management
- Built-in template engine support

### Ecosystem
- Rich middleware ecosystem
- Good documentation
- Active community
- Regular updates and maintenance

## YugabyteDB Features

This boilerplate leverages YugabyteDB, a PostgreSQL-compatible distributed SQL database that provides:

- **PostgreSQL Compatibility**: Works seamlessly with existing PostgreSQL drivers and ORMs
- **Horizontal Scalability**: Scales out across multiple nodes
- **High Availability**: Built-in replication and fault tolerance
- **Geo-Distribution**: Deploy data across regions
- **ACID Compliance**: Full ACID transactions across distributed nodes

### YugabyteDB Web UI Access

When running with Docker Compose, you can access:

- YB-Master UI: http://localhost:7001
- YB-TServer UI: http://localhost:9001

## Architecture

This boilerplate follows clean architecture principles:

1. **Handlers Layer**: HTTP request handling and response formatting
2. **Service Layer**: Business logic and orchestration
3. **Repository Layer**: Data access and persistence
4. **Models**: Domain entities and data transfer objects

Each layer has clear responsibilities and dependencies flow inward, making the code testable and maintainable.

## Observability

This boilerplate includes comprehensive observability through OpenTelemetry:

### Distributed Tracing

- **Automatic HTTP tracing** for all API endpoints via Fiber middleware
- **Database query tracing** with GORM OpenTelemetry plugin
- **Service layer tracing** with custom spans for business logic
- **Error recording** and span attributes for debugging

### Metrics Collection

- **Request duration metrics** for service operations
- **Operation counters** with success/failure status
- **Custom business metrics** for user operations

### Jaeger Integration

Access the Jaeger UI at `http://localhost:16686` to:

- View distributed traces across service boundaries
- Analyze request flows and performance bottlenecks
- Debug errors with detailed span information
- Monitor service dependencies

### Adding Custom Tracing

```go
// In your service methods
ctx, span := otel.Tracer("service-name").Start(ctx, "operation-name")
defer span.End()

// Add attributes
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("items.count", count),
)

// Record errors
if err != nil {
    span.RecordError(err)
}
```

## 🔄 Temporal Workflow Engine

### Overview

Temporal provides durable, reliable workflow execution for complex business processes.

### Features

- **Durable Execution**: Workflows survive failures and restarts
- **Built-in Retries**: Automatic retry with exponential backoff
- **Long-Running Processes**: Support for workflows that run for days/months
- **Visibility**: Full workflow history and state inspection

### Sample Workflow: User Onboarding

The boilerplate includes a user onboarding workflow that:

1. Sends welcome email
2. Creates user profile
3. Sends notifications
4. Handles failures gracefully

### Usage Example

```bash
# Start workflow
curl -X POST http://localhost:8080/api/v1/workflows/user-onboarding \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "workflow_id": "user-onboarding-123",
    "input": {
      "user_id": 1,
      "email": "user@example.com"
    }
  }'
```

## 🛡️ OPA Authorization

### Overview

Open Policy Agent provides fine-grained, policy-based authorization.

### Authorization Rules

- **Public endpoints**: Health check, login, register
- **Role-based access**:
  - `admin`: Full access to all endpoints
  - `user`: Can only access their own profile
  - `workflow_executor`: Can trigger workflows
  - `premium`: Higher rate limits

### Policy Testing

```bash
# Test policy directly
curl -X POST http://localhost:8181/v1/data/authz/allow \
  -d '{
    "input": {
      "method": "GET",
      "path": "/api/v1/users",
      "user": {"id": "1", "roles": ["admin"]}
    }
  }'
```

## 🔍 Troubleshooting

### Common Issues

#### Temporal Connection

```bash
# Check Temporal health
docker logs fiber-boilerplate-temporal
docker exec fiber-boilerplate-temporal-admin tctl cluster health
```

#### OPA Policy Issues

```bash
# View OPA decision logs
docker logs fiber-boilerplate-opa
```

#### Database Migration

```bash
# Manual migration
docker exec -e PGPASSWORD=yugabyte fiber-boilerplate-db \
  bin/ysqlsh -h 172.19.0.3 -p 5433 -U yugabyte -d yugabyte \
  -f /migrations/000001_create_users_table.up.sql
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.