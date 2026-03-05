# Go Web Enterprise Demo

A production-ready Go web application following enterprise best practices with clean architecture.

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── handler/              # HTTP handlers (presentation layer)
│   ├── middleware/           # HTTP middleware
│   ├── repository/           # Data access layer
│   ├── router/               # Route definitions
│   └── service/              # Business logic layer
├── pkg/
│   └── logger/               # Structured logging
├── configs/
│   └── config.yaml           # Configuration file
├── data/                     # Database storage (SQLite)
├── Dockerfile                # Docker image
├── docker-compose.yml        # Docker compose configuration
└── go.mod                    # Go module definition
```

## Features

- **Clean Architecture**: Separation of concerns with handler -> service -> repository layers
- **Dependency Injection**: Wiring dependencies through main.go
- **Structured Logging**: JSON logging with zerolog
- **Configuration Management**: YAML-based config with viper
- **Middleware Stack**:
  - Request ID tracking
  - Structured logging
  - Panic recovery
  - CORS support
- **Graceful Shutdown**: Proper handling of SIGINT/SIGTERM
- **Docker Support**: Multi-stage Dockerfile for production builds
- **Health Checks**: Liveness and readiness endpoints

## API Endpoints

### Health Check
- `GET /health/liveness` - Liveness probe
- `GET /health/readiness` - Readiness probe with DB status

### Hello API
- `GET /api/v1/hello` - Hello endpoint

### User Management
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## Getting Started

### Prerequisites
- Go 1.23+
- Docker (optional)

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run cmd/server/main.go
```

3. Test the API:
```bash
curl http://localhost:8080/api/v1/hello
curl http://localhost:8080/health/liveness
```

### Docker Deployment

1. Build and run with Docker Compose:
```bash
docker-compose up -d
```

2. Build the Docker image:
```bash
docker build -t go-web-demo .
```

3. Run the container:
```bash
docker run -p 8080:8080 -v $(pwd)/data:/app/data go-web-demo
```

## Configuration

Edit `configs/config.yaml` to customize:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug, release, test

database:
  driver: "sqlite"
  dbname: "./data/app.db"

log:
  level: "info"
  format: "json"
```

## Environment Variables

- `PORT` - Override server port (default: 8080)

## License

MIT
