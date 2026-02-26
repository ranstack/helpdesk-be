# Copilot Instructions - Helpdesk API

## Project Context
- Go REST API using Echo v5 and PostgreSQL.
- Clean architecture: handler -> service -> repository.
- Standardized JSON responses and AppError codes.
- Tests are currently removed; manual testing is preferred.

## Coding Style
- Go idioms and standard library first.
- Keep functions small and focused.
- Prefer explicit error handling; avoid panic.
- No comments unless absolutely necessary.
- Use ASCII only in source files unless a file already uses Unicode.

## Architecture Rules
- Handlers: HTTP request/response only; no business logic.
- Services: validation, business rules, and error translation.
- Repositories: database access only; no validation.
- Do not bypass service layer from handlers.

## Error Handling
- Use AppError helpers from internal/utils/errors.
- Error codes must be uppercase with underscores.
- Convert input parsing errors to BadRequest with a clear message.

## Responses
- Use internal/utils/response helpers for JSON responses.
- Success response: { success, message, data }.
- Error response: { success: false, error: { code, message, details } }.

## Validation
- Use internal/utils/validator for request validation.
- Return Validation errors with details map when invalid.

## Database
- Use sqlx with PostgreSQL.
- Keep queries in repositories; no SQL in services or handlers.

## Dependencies
- Avoid adding new libraries unless required.
- Use go mod tidy after adding/removing dependencies.

## Files and Structure
- New features live in internal/features/<feature>.
- Each feature should have dto.go, handler.go, models.go, repository.go, routes.go, service.go.

## Testing
- No automated tests expected unless explicitly requested.
- Prefer manual verification via curl or Postman.

## README Updates
- Keep README aligned with actual project behavior.
- Update endpoints and configuration when changes are made.
