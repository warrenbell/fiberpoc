# fiberpoc

POC for a Go Fiber REST API. This POC demonstrates the useage of env vars, logging, unit testing, integration testing, deployment to a cloud provider, database migration and database seeding. Integration testing is done with a live postgres insatnce spun up in a container. I am going to try a few things to make the database setup fast for integration tests.

## Prerequisites

- Go 1.23.4 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the `app` directory with the following variables:

```env
# Possible values "debug", "info", "warn", "error", "dpanic", "panic", "fatal"
LOG_LEVEL=debug

# Possible values true or false
LOG_TO_FILE=true

# A valid PostgreSQL url
POSTGRESQL_URL=postgresql://user:password@localhost:5432/dbname
```

## Getting Started

1. Clone the repository:

```bash
git clone [repository-url]
cd fiberpoc
```

2. Install dependencies:

```bash
cd app
go mod download
cd ../common
go mod download
```

3. Set up the environment.

4. Run the application:

```bash
cd app/cmd
go run main.go
```

## Project Components

### App Package

- `handlers/`: HTTP request handlers and routing logic
- `migrations/`: Database schema migrations
- `seeds/`: Initial data for database seeding
- `tests/`: Integration tests

### Common Package

- `models/`: Data structures and domain models
- `repos/`: Data access layer and repository implementations
- `services/`: Business logic and service layer
- `interfaces/`: Interface definitions for dependency injection
- `clients/`: External service client implementations

## Dependency Injection

The project implements a clean dependency injection pattern that promotes:

### Interface-Driven Design

- All dependencies are defined as interfaces in the `common/interfaces` package
- This allows for loose coupling between components and easier testing
- Implementation details are hidden behind interface contracts

### Service Layer Pattern

```go
// Example interface definition
type UserService interface {
    GetUser(ctx context.Context, id string) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
}

// Example service implementation
type userService struct {
    repo    interfaces.UserRepository
    logger  interfaces.Logger
}

// Constructor with dependency injection
func NewUserService(
    repo interfaces.UserRepository,
    logger interfaces.Logger,
) interfaces.UserService {
    return &userService{
        repo: repo,
        logger: logger,
    }
}
```

### Testing Benefits

- Mock implementations are provided in the `common/mocks` package
- Easy to swap real implementations with mocks for unit testing
- Enables isolated testing of individual components
- Mocks are created from an interface using mockgen. See: [GoMock](https://github.com/uber-go/mock)

### Dependency Graph

The application follows a clear dependency hierarchy:

1. Handlers depend on Services
2. Services depend on Repositories and Clients
3. Repositories depend on Database connections
4. All dependencies are injected at startup

This structure ensures:

- Clear separation of concerns
- Testable code
- Maintainable architecture
- Easy to extend functionality

## Key Features

- Modular architecture with clear separation of concerns
- PostgreSQL integration using pgx driver
- Structured logging with Uber's Zap logger
- Environment-based configuration
- Database migrations and seeding support
- Mock implementations for testing
- Comprehensive dependency injection system

## Development

### Running Tests

```bash
cd app/
go test ./...

cd common/
go test ./...
```

### Adding New Features

1. Define interfaces in `common/interfaces`
2. Implement business logic in `common/services`
3. Add data access in `common/repos` if needed
4. Create HTTP handlers in `app/handlers`
5. Update routes in `app/cmd/main.go`

### Adding New Dependencies

1. Define the interface in `common/interfaces`
2. Create the implementation in appropriate package
3. Create mock using mockgen and add implementation in `common/mocks`
4. Wire up the dependencies in main.go
5. Inject into required services
