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
- In dto.go files: organize with constants first, then all type definitions, then methods, then functions.

## Architecture Rules
- Handlers: HTTP request/response only; no business logic.
- Services: validation, business rules, and error translation.
- Repositories: database access only; no validation.
- Do not bypass service layer from handlers.

## REST API Conventions
- Use PATCH for partial updates (only changed fields sent).
- Use PUT only for full resource replacement (all fields required).
- JSON tags must use camelCase (e.g., `json:"createdAt"` not `json:"created_at"`).

## Error Handling
- Use AppError helpers from internal/utils/errors.
- Error codes must be uppercase with underscores.
- Convert input parsing errors to BadRequest with a clear message.

## Responses
- Use internal/utils/response helpers for JSON responses.
- Success response: { success, message, data }.
- Error response: { success: false, error: { code, message, details } }.
- Prefer reusing a shared response DTO per feature (for example, one `CategoryResponse`) to reduce boilerplate.
- Add endpoint-specific response DTOs only when the response contract is intentionally different.
- For paginated lists: use `response.ListResponse[T]` generic type from internal/utils/response for consistent structure and reduced duplication.
  - Example: Service returns `*response.ListResponse[CategoryResponse]` with items array and pagination metadata.
  - Eliminates feature-specific ListResponse types (e.g., CategoryListResponse, DivisionListResponse).
- List endpoints should return { items: [], pagination: { page, limit, totalItems, totalPages } }.

## Validation
- Use internal/utils/validator for request validation.
- Return Validation errors with details map when invalid.
- For enums, define constants and ValidValues map in models.go, then validate against it in dto.go.
  - Example: `const (RoleAdmin = "ADMIN"; RoleIT = "IT"; RoleStaff = "STAFF")`
  - Validate with: `if !ValidRoles[role] { v.AddError("role", "Must be one of: ADMIN, IT, STAFF") }`
- For email validation, use `validator.ValidateEmail(email)` helper:
  - Example: `if email != "" && !validator.ValidateEmail(email) { v.AddError("email", "Must be a valid email address") }`

## Pagination
- Use query parameters: `page`, `limit` for pagination controls.
- Default: page=1, limit=10; enforce max limit (e.g., 100).
- **Shared Pagination:** Embed `response.PaginationQuery` in feature query DTOs to reuse pagination logic.
  - Example: `type GetCategoriesQuery struct { response.PaginationQuery; Name string; IsActive *bool; }`
  - Call `query.NormalizePagination()` to get normalized page, limit, offset values.
  - Constants available: `response.DefaultPage=1`, `response.DefaultLimit=10`, `response.MaxLimit=100`
  - Use `response.ParseDate(dateStr)` helper for parsing date filters (returns *time.Time or error).
  - Use `response.CalculateTotalPages(totalItems, limit)` helper to calculate total pages (handles zero case).
- Repository: run COUNT query for total, then paginated SELECT with LIMIT/OFFSET.
- Return response with items array + pagination metadata (page, limit, totalItems, totalPages).

## Database
- Use sqlx with PostgreSQL.
- Keep queries in repositories; no SQL in services or handlers.
- Keep `SELECT`/`RETURNING` columns aligned with model `db` tags for fields exposed in responses.
- For text fields requiring case-insensitive uniqueness (e.g., name, email), add unique index on LOWER(column).
- Use `ILIKE` for case-insensitive searches in WHERE clauses.

## Migrations
- Use Goose for database migrations.
- Project is not production yet: update existing migrations directly, then reset/re-apply rather than creating new migrations for schema changes.

## Dependencies
- Avoid adding new libraries unless required.
- Use go mod tidy after adding/removing dependencies.

## Files and Structure
- New features live in internal/features/<feature>.
- Each feature should have dto.go, handler.go, models.go, repository.go, routes.go, service.go.
- Use standard CRUD method order in interfaces and implementations:
  1. All GET methods (GetAll, GetByID, GetByName, Exists, etc.)
  2. All CREATE methods (Create)
  3. All UPDATE methods (Update)
  4. All DELETE methods (Delete)

## Code Review & DRY Principle
- After implementing a feature, scan other features for code duplication.
- Extract common patterns into shared utilities (e.g., response, validator, errors packages).
- Examples of duplication to watch for:
  - Pagination logic (normalize params, calculate offsets, totalPages)
  - Date parsing and validation
  - Filter/query building patterns
  - Repetitive validation rules
  - Common response/error mapping
- Keep utilities in `internal/utils/` organized by concern (response, errors, validator, etc.).

## Testing
- No automated tests expected unless explicitly requested.
- Prefer manual verification via curl or Postman.

## README Updates
- Keep README aligned with actual project behavior.
- Update endpoints and configuration when changes are made.
