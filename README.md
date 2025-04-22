# fiberpoc
 POC for a Go Fiber REST API

## Prerequisites

- Go 1.23.4 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the `app` directory with the following variables:

```env
# Copy from app/.env.example
DATABASE_URL=postgresql://user:password@localhost:5432/dbname
PORT=3000
ENV=development
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

3. Set up the environment:
```bash
cp app/.env.example app/.env
# Edit app/.env with your configuration
```

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
- `tests/`: Integration and end-to-end tests

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
    cache   interfaces.CacheClient
}

// Constructor with dependency injection
func NewUserService(
    repo interfaces.UserRepository,
    logger interfaces.Logger,
    cache interfaces.CacheClient,
) interfaces.UserService {
    return &userService{
        repo: repo,
        logger: logger,
        cache: cache,
    }
}
```

### Testing Benefits
- Mock implementations are provided in the `common/mocks` package
- Easy to swap real implementations with mocks for unit testing
- Enables isolated testing of individual components

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
cd app/tests
go test ./...

cd ../../common
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
3. Add mock implementation in `common/mocks`
4. Wire up in the dependency injection container
5. Inject into required services

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[Add your license information here]

## Contact

[Add your contact information here]
