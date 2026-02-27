# Helpdesk REST API

A production-ready REST API for IT Helpdesk management system built with Go, Echo framework, and PostgreSQL. The API follows clean architecture principles with repository, service, and handler layers.

**Project Structure:** [ERD Diagram](https://app.eraser.io/workspace/MCKUzCCls92JCU5rpuew?origin=share)

## Stack

- **Language:** Go 1.22+
- **Web Framework:** Echo v5
- **Database:** PostgreSQL
- **Database Migration:** Goose
- **Logging:** Structured logging (slog)
- **Containerization:** Docker

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Git
- Docker (optional, for containerization)

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd helpdesk/server
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Setup

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Configure your environment variables:

```env
APP_NAME=Helpdesk API
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=helpdesk
DB_SSLMODE=disable
```

### 4. Database Setup

Create the PostgreSQL database:

```bash
createdb helpdesk
```

Run migrations:

```bash
goose -dir migrations postgres "user=$DB_USER password=$DB_PASSWORD host=$DB_HOST port=$DB_PORT dbname=$DB_NAME sslmode=$DB_SSLMODE" up
```

### 5. Run the Server

```bash
go run ./cmd/api/main.go
```

Server will start on `http://localhost:8080`

## Project Structure

```
cmd/
  └── api/
      └── main.go              # Application entry point

internal/
  ├── config/
  │   └── config.go            # Configuration management
  ├── database/
  │   └── postgres.go          # Database initialization
  ├── features/
  │   └── category/            # Category feature (CRUD)
  │       ├── dto.go           # Request/Response DTOs
  │       ├── handler.go       # HTTP handlers
  │       ├── models.go        # Domain models
  │       ├── repository.go    # Data access layer
  │       ├── routes.go        # Route definitions
  │       └── service.go       # Business logic layer
  ├── middleware/
  │   ├── cors.go              # CORS middleware
  │   ├── logger.go            # Request logging
  │   └── recovery.go          # Panic recovery
  └── utils/
      ├── errors/
      │   └── errors.go        # Error types and helpers
      ├── response/
      │   └── response.go      # Response formatting
      └── validator/
          └── validator.go     # Input validation

migrations/                     # Database migrations (Goose)
```

## Architecture

The project follows **Clean Architecture** principles:

1. **Handler Layer** - HTTP request/response handling with Echo
2. **Service Layer** - Business logic, validation, and domain rules
3. **Repository Layer** - Data access abstraction with interface-based design
4. **Database Layer** - PostgreSQL connection and query execution

### Design Patterns

- **Interface-based Repository** - Enables easy mocking and testing
- **Dependency Injection** - Service and handler dependencies injected at initialization
- **Error Handling** - Centralized AppError type with proper HTTP status codes
- **Request/Response DTOs** - Separation of API contracts from domain models
- **Middleware Stack** - Logger, Recovery, and CORS middleware

## API Endpoints

### Category Management

All endpoints are prefixed with `/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/categories` | Create a new category |
| GET | `/categories` | Get all categories |
| GET | `/categories/:id` | Get category by ID |
| PATCH | `/categories/:id` | Update category |
| DELETE | `/categories/:id` | Delete category |

`GET /categories` supports query parameters:

| Query | Type | Description |
|-------|------|-------------|
| `page` | number | Page number (default `1`) |
| `limit` | number | Items per page (default `10`, max `100`) |
| `name` | string | Case-insensitive partial search by category name |
| `isActive` | boolean | Filter active/inactive categories |
| `createdAt` | string | Filter by creation date in `YYYY-MM-DD` |

### Health Check

```
GET /api/v1/health
```

Returns:
```json
{
  "status": "ok",
  "app": "Helpdesk API"
}
```

## Error Handling

The API uses standardized error responses with specific error codes:

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Category not found",
    "details": {}
  }
}
```

**Error Codes:**
- `NOT_FOUND` (404) - Resource not found
- `ALREADY_EXISTS` (409) - Resource already exists
- `VALIDATION_ERROR` (400) - Input validation failed
- `BAD_REQUEST` (400) - Invalid request
- `INTERNAL_SERVER_ERROR` (500) - Server error

## Response Format

All successful responses follow this format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

## Running Tests

Currently no automated tests included. Manual testing recommended using:

- **Postman** - API testing collection
- **curl** - Command-line testing
- **Thunder Client** - VS Code extension

Example test:

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Hardware"}'
```

## Development

### Code Style

- Follow Go conventions and idioms
- Use meaningful variable names
- Keep functions focused and single-purpose

### Running with Hot Reload

Install and use `air` for hot reloading:

```bash
go install github.com/cosmtrek/air@latest
air
```

Configuration is in `.air.toml`

### Database Migrations

Create a new migration:

```bash
goose create add_users_table sql
```

Run migrations:

```bash
goose up
```

Rollback:

```bash
goose down
```

## Docker

Build Docker image:

```bash
docker build -f Dockerfile -t helpdesk-api .
```

Run container:

```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=helpdesk \
  helpdesk-api
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | Helpdesk API | Application name |
| `APP_PORT` | 8080 | Server port |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | PostgreSQL user |
| `DB_PASSWORD` | postgres | PostgreSQL password |
| `DB_NAME` | helpdesk | Database name |
| `DB_SSLMODE` | disable | SSL mode for connection |

## Future Features

- [ ] User management
- [ ] Ticket management
- [ ] Division/Department management
- [ ] Authentication & Authorization
- [ ] Ticket attachments
- [ ] Ticket resolutions
- [ ] User roles and permissions
- [ ] API documentation (Swagger)

## Troubleshooting

### Database Connection Error

```
Db connection error: connection refused
```

**Solution:** Ensure PostgreSQL is running and credentials are correct in `.env`

### Port Already in Use

```
listen tcp :8080: bind: An attempt was made to reuse a socket address
```

**Solution:** Change `APP_PORT` in `.env` or kill the process using port 8080

### Migration Errors

Ensure migrations are run in correct order. Check migration status:

```bash
goose -dir migrations postgres "..." status
```

## License

This project is proprietary and closed source.

## Support

For issues or questions, contact the development team.
