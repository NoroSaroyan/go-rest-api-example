# Todo API

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/NoroSaroyan/go-rest-api-example?style=for-the-badge)](https://goreportcard.com/report/github.com/NoroSaroyan/go-rest-api-example)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/NoroSaroyan/go-rest-api-example/ci.yml?branch=main&style=for-the-badge&logo=github-actions&logoColor=white)](https://github.com/NoroSaroyan/go-rest-api-example/actions)
[![Codecov](https://img.shields.io/codecov/c/github/NoroSaroyan/go-rest-api-example?style=for-the-badge&logo=codecov&logoColor=white)](https://codecov.io/gh/NoroSaroyan/go-rest-api-example)

[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)](https://swagger.io/)

[![Clean Architecture](https://img.shields.io/badge/architecture-clean-blue?style=for-the-badge)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![API](https://img.shields.io/badge/API-RESTful-orange?style=for-the-badge)](http://localhost:8080/swagger/)
[![GitHub stars](https://img.shields.io/github/stars/NoroSaroyan/go-rest-api-example?style=for-the-badge&logo=github)](https://github.com/NoroSaroyan/go-rest-api-example/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/NoroSaroyan/go-rest-api-example?style=for-the-badge&logo=github)](https://github.com/NoroSaroyan/go-rest-api-example/network/members)

My attempt at a clean, idiomatic-go REST API for managing todo items, built with Go and following clean architecture
principles. Perfect for learning modern Go development practices or as a foundation for your next project.

## Status

This project maintains high code quality standards with:

- **Continuous Integration**: Automated testing on every commit
- **Code Coverage**: Comprehensive test coverage tracked via Codecov
- **Code Quality**: Static analysis with golangci-lint and Go Report Card
- **Security**: Vulnerability scanning with gosec and govulncheck
- **Documentation**: Auto-generated API docs with Swagger/OpenAPI

## Quick Start

### Prerequisites

- Go 1.21+ installed
- Docker & Docker Compose (for database)

### Get it running in 3 steps:

1. **Clone and setup**
   ```bash
   git clone https://github.com/NoroSaroyan/go-rest-api-example.git
   cd go-rest-api-example
   ```

2. **Create your local config**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Start everything**
   ```bash
   make docker-up    # Start PostgreSQL
   make migrate-up   # Create database tables
   make run         # Start the API server
   ```

That's it! Your API is now running at `http://localhost:8080`

## Try it out

### Option 1: Interactive Swagger UI (Recommended)

Open your browser and go to: **http://localhost:8080/swagger/**

You'll get a beautiful, interactive interface where you can:

- See all available endpoints
- Try requests with real data
- View response examples
- No Postman needed!

### Option 2: curl commands

```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}'

# Get all todos  
curl http://localhost:8080/api/v1/todos

# Get a specific todo
curl http://localhost:8080/api/v1/todos/1

# Delete a todo
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

## Development

### Available commands

```bash
make run          # Start the server
make build        # Build the binary
make test         # Run all tests
make docs         # Generate API documentation
make lint         # Run code quality checks

# Database operations
make docker-up    # Start PostgreSQL container
make docker-down  # Stop containers
make migrate-up   # Apply database migrations
make migrate-down # Rollback migrations
```

### Project structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── app/            # Application setup & routing
│   ├── config/         # Configuration management
│   ├── domain/         # Business entities & errors
│   ├── service/        # Business logic layer
│   ├── repository/     # Data access layer
│   └── transport/      # HTTP handlers & middleware
├── migrations/         # Database schema changes
├── docs/              # Auto-generated API documentation
└── .env               # Local environment configuration
```

### Why this architecture?

**Clean Architecture** means each layer has a single responsibility:

- **Domain**: Core business rules (what is a todo?)
- **Service**: Business logic (how do we create todos?)
- **Repository**: Data access (where do we store todos?)
- **Transport**: HTTP handling (how do clients talk to us?)

This makes the code:

- Easy to test (mock any layer)
- Easy to change (swap PostgreSQL for MongoDB? No problem!)
- Easy to understand (each part has one job)

## Configuration

Create a `.env` file in the project root:

```env
# Server
APP_PORT=8080

# Database  
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# Logging
LOG_LEVEL=debug
```

### Available environment variables:

| Variable      | Default     | Description                           |
|---------------|-------------|---------------------------------------|
| `APP_PORT`    | `8080`      | Port for the HTTP server              |
| `DB_HOST`     | `localhost` | PostgreSQL host                       |
| `DB_PORT`     | `5432`      | PostgreSQL port                       |
| `DB_USER`     | `todo`      | Database username                     |
| `DB_PASSWORD` | `todo`      | Database password                     |
| `DB_NAME`     | `todo_db`   | Database name                         |
| `LOG_LEVEL`   | `info`      | Logging level (debug/info/warn/error) |

## Testing

We take testing seriously! Run the test suite with:

```bash
make test
```

**What's tested:**

- **Service layer** - Business logic with mock repositories
- **Configuration** - Environment loading and validation
- **Validation** - Request validation and error handling
- **Database operations** - Repository layer functionality

**Testing philosophy**: Each layer is tested in isolation using dependency injection and mocks. This means fast,
reliable tests that don't depend on external services.

## API Documentation

### Interactive Documentation

Visit `http://localhost:8080/swagger/` for the full interactive API documentation.

### Quick Reference

| Method   | Endpoint             | Description         |
|----------|----------------------|---------------------|
| `POST`   | `/api/v1/todos`      | Create a new todo   |
| `GET`    | `/api/v1/todos`      | List all todos      |
| `GET`    | `/api/v1/todos/{id}` | Get a specific todo |
| `DELETE` | `/api/v1/todos/{id}` | Delete a todo       |
| `GET`    | `/health`            | Health check        |

### Example requests/responses

**Create a todo:**

```bash
POST /api/v1/todos
{
  "title": "Buy groceries"
}

# Response: 201 Created
{
  "id": 1
}
```

**Get all todos:**

```bash
GET /api/v1/todos

# Response: 200 OK
[
  {
    "id": 1,
    "title": "Buy groceries", 
    "completed": false,
    "created_at": "2023-01-01T12:00:00Z"
  }
]
```

## Deployment

### Building for production

```bash
# Build optimized binary
make build

# The binary will be created as ./todo-api
./todo-api
```

### Environment considerations

- Set `LOG_LEVEL=info` in production
- Use strong database credentials
- Consider using environment-specific configuration files
- Set up proper database connection pooling (already configured!)

## Contributing

This project follows Go best practices and clean architecture principles. When contributing:

1. **Write tests** for new functionality
2. **Follow the existing patterns** in each layer
3. **Update documentation** if you change APIs
4. **Run the full test suite** before submitting

```bash
# Before submitting changes
make test
make lint
make docs  # If you changed API endpoints
```

## What you'll learn

Using this project, you'll get hands-on experience with:

- **Clean Architecture** in Go
- **Dependency Injection** patterns
- **Interface-driven development**
- **Testing with mocks**
- **Database migrations**
- **Structured logging**
- **API documentation** with Swagger
- **Docker** for development
- **Environment-based configuration**

## Resources and Literature

This project was built following established patterns and best practices from the Go community. Here are key resources
that influenced the design and implementation:

### Architecture and Design Patterns

- **Clean Architecture** by Robert C. Martin - The foundational principles for the layered architecture
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob's
  original blog post
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn - Ports and
  adapters pattern
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle) - Core principle for
  testable code

### Go-Specific Resources

- [Effective Go](https://golang.org/doc/effective_go.html) - Official Go documentation on writing idiomatic Go
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Style guide for Go code
- [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) by Bill Kennedy
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout) - Community conventions for Go
  project structure

### API Design and Documentation

- [REST API Design Best Practices](https://restfulapi.net/) - RESTful API design principles
- [OpenAPI Specification](https://swagger.io/specification/) - API documentation standards
- [API Design Guidelines](https://cloud.google.com/apis/design) by Google Cloud

### Testing and Quality

- [Test Driven Development](https://martinfowler.com/bliki/TestDrivenDevelopment.html) by Martin Fowler
- [Go Testing](https://golang.org/pkg/testing/) - Official Go testing documentation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit for Go (not used but influential)

### Database and Persistence

- [Database Design](https://en.wikipedia.org/wiki/Database_design) - Relational database design principles
- [PostgreSQL Documentation](https://www.postgresql.org/docs/) - Official PostgreSQL documentation
- [Go Database/SQL Tutorial](https://golang.org/doc/database/sql) - Official Go database handling

### Configuration and Deployment

- [12-Factor App](https://12factor.net/) - Methodology for building SaaS applications
- [Environment Variables in Go](https://golang.org/pkg/os/#Getenv) - Configuration management patterns
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/) - Containerization guidelines

### Tools and Libraries Used

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router and URL matcher
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit
- [Zap](https://github.com/uber-go/zap) - Structured logging library
- [Swag](https://github.com/swaggo/swag) - Swagger documentation generator
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader

## License

MIT License - feel free to use this code for learning, personal projects, or as a foundation for commercial
applications.

## Questions?

- Check the **Swagger UI** at `/swagger/` for API details
- Look at the **test files** for usage examples
- Explore the **`internal/`** directory to understand the architecture
- The code is heavily documented - each package has clear documentation

---

**Happy coding!** If this project helped you learn something new, consider giving it a star.