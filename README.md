# Bookstore API Assessment

Build a REST API for managing books and authors.

## Project Structure
- `openapi.yaml`: API specification – refer to this for request/response formats
- `internal/api/types.go`: Request/response types and validation
- `internal/api/handlers.go`: HTTP handlers using standard `net/http`
- `internal/db/`: SQLC-generated database code
- `internal/db/migrations/`: Database schema and seed data

## Implementation Notes
- Most SQL queries are already implemented in `internal/db/queries/`
- Look at `handleGetBook` and `handleCreateBook` in `handlers.go` for reference implementations

## Tasks

1. Implement the remaining handlers:
    - `GET /api/v1/authors`: List all authors
    - `GET /api/v1/books`: List all books
    - `POST /api/v1/books`: Create a new book
    - `PUT /api/v1/books/{id}`: Update a book
    - `DELETE /api/v1/books/{id}`: Delete a book

2. Implement the author stats query and handler:
    - `GET /api/v1/authors/{id}/stats`: Get author statistics
    - Write the SQL query to calculate:
        - Total number of books
        - Average book price
        - Publication date range
        - Total revenue
        - Books published per year
    - See `AuthorStats` type in `types.go` and the existing  queries in `internal/db/queries`.

3. Improve the code organization:
    - Split handlers into logical files if needed
    - Extract common functionality
    - Add middleware if useful
    - Improve error handling
    - Any other improvements you think would help maintainability

4. Describe your testing strategy (no need to implement the tests - we're interested in your approach).

## Prerequisites
Before attempting this assessment, ensure the following tools are installed on your system:
- [**Go** (1.22.0 or higher)](https://go.dev/dl/)
- [**Docker**](https://www.docker.com/products/docker-desktop/)
- [**Make**](https://chatgpt.com/share/673b0945-a79c-8000-926c-8c34bfac6b43)

## Getting Started
```bash
# Install dependencies and tools
make setup

# Start PostgreSQL
make docker-up

# Run migrations
make migrate-up

# Run the API
make run
# API will be available at http://localhost:8080
```

## Development Notes
- Use `make run` to restart the API after changes
- Check logs for request handling and errors
- See `Makefile` for all available commands
- PostgreSQL is accessible at localhost:5432
- Database credentials: bookstore/bookstore
- Each time you run `make migrate-up`, the database is seeded with:
   - A set of authors and books
   - 10,000 randomized book sales records
   - **Note:** Sales data will be different after each migration

## Evaluation Criteria
- Clean, idiomatic Go code
- Proper error handling
- SQL query implementation
- Understanding of HTTP status codes
- Testing approach
