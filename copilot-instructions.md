# Copilot Instructions for Vania

## Project Overview
Vania is a Go project following standard Go conventions and best practices. The project emphasizes clean architecture, testability, and idiomatic Go code. The architecture follows the standard Go project layout with clear separation of concerns between packages.

## Code Organization

### Directory Structure
Follow the standard Go project layout:
- **cmd/**: Main applications for this project. The directory name should match the executable name.
- **internal/**: Private application and library code. Code that you don't want others importing.
- **pkg/**: Library code that's ok to use by external applications.
- **api/**: API definitions (OpenAPI/Swagger specs, protocol buffers, etc.)
- **web/**: Web application specific components (static files, templates, etc.)
- **configs/**: Configuration file templates or default configs.
- **scripts/**: Scripts for builds, installs, analysis, etc.
- **test/**: Additional external test apps and test data.
- **docs/**: Design and user documents.
- **examples/**: Examples for your applications and/or public libraries.

### Package Naming
- Use short, lowercase package names (e.g., `user`, `http`, `config`)
- Package names should be singular, not plural (e.g., `user` not `users`)
- Avoid generic names like `util`, `common`, `base` - be specific about purpose
- Package name should match the directory name
- Use descriptive names that clearly indicate the package's purpose

**Examples:**
```go
// Good
package user
package storage
package httpserver

// Bad
package users          // Plural form
package my_package     // Underscores
package HTTPServer     // Mixed case
```

### File Naming
- Use lowercase with underscores for multi-word names (e.g., `user_service.go`)
- Test files must end with `_test.go`
- Main package files typically named `main.go`
- Group related functionality in the same file
- Keep files focused on a single responsibility

## Coding Standards

### Error Handling

**Primary Pattern**: Always handle errors explicitly. Never ignore errors unless there's a compelling reason (and add a comment explaining why).

```go
// Bad
result, _ := DoSomething()

// Good
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

**Error Wrapping**: Use `fmt.Errorf` with `%w` verb to wrap errors, preserving the error chain for `errors.Is` and `errors.As`.

```go
// Wrap errors with context
func (s *Service) GetUser(id string) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("service: failed to get user %s: %w", id, err)
    }
    return user, nil
}
```

**Custom Error Types**: Define custom error types for domain-specific errors in their respective packages.

```go
// Define sentinel errors for expected error conditions
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidInput     = errors.New("invalid input")
    ErrUnauthorized     = errors.New("unauthorized access")
)

// Custom error types for errors requiring additional context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// Usage
if user == nil {
    return nil, ErrUserNotFound
}

// Check specific error types
if errors.Is(err, ErrUserNotFound) {
    // Handle not found case
}

var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error with details
}
```

### Naming Conventions

#### Variables
- Use camelCase for local variables and function parameters
- Use short names for short-lived variables (e.g., `i` for loop index, `err` for errors)
- Use descriptive names for variables with larger scope
- Avoid single-letter names except for: `i`, `j`, `k` (indexes), `n` (count), `err` (errors)

```go
// Good
var userCount int
var maxRetries = 3
for i := 0; i < n; i++ {
    // ...
}

// Bad
var uc int              // Too cryptic
var maximumNumberOfRetries = 3  // Too verbose
```

#### Functions
- Use MixedCaps (PascalCase) for exported functions
- Use mixedCaps (camelCase) for unexported functions
- Function names should be verbs or verb phrases
- Getter functions should not have `Get` prefix (e.g., `User()` not `GetUser()`)
- Setter functions should have `Set` prefix (e.g., `SetUser()`)

```go
// Good
func CreateUser(name string) *User { }
func (u *User) Name() string { }
func (u *User) SetName(name string) { }

// Bad
func create_user(name string) *User { }  // Snake case
func (u *User) GetName() string { }      // Unnecessary Get prefix
```

#### Constants
- Use MixedCaps (PascalCase) for exported constants
- Use mixedCaps (camelCase) for unexported constants
- Group related constants using `const` blocks with `iota` for enumerations

```go
// Good
const MaxConnections = 100
const defaultTimeout = 30 * time.Second

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
    StatusDeleted
)

// Bad
const MAX_CONNECTIONS = 100  // Snake case
const Default_Timeout = 30   // Mixed styles
```

#### Interfaces
- Interface names should end with `-er` suffix when possible (e.g., `Reader`, `Writer`, `Formatter`)
- Single-method interfaces should be named with method name + `-er`
- Keep interfaces small and focused (prefer many small interfaces over large ones)
- Define interfaces in the package that uses them, not the package that implements them

```go
// Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

type UserRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
    Delete(id string) error
}

// Bad
type IUserRepository interface { }  // Don't prefix with I
type UserRepositoryInterface interface { }  // Don't suffix with Interface
```

### Testing

#### Test File Naming
- Test files must end with `_test.go`
- Place tests in the same package as the code being tested
- Use `_test` package suffix for black-box testing (testing public API only)

```go
// user.go - production code
package user

// user_test.go - white-box testing
package user

// user_test.go - black-box testing
package user_test
```

#### Table-Driven Tests
Use table-driven tests for testing multiple scenarios of the same functionality.

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        want    bool
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@example.com",
            want:    true,
            wantErr: false,
        },
        {
            name:    "invalid email - missing @",
            email:   "userexample.com",
            want:    false,
            wantErr: true,
        },
        {
            name:    "empty email",
            email:   "",
            want:    false,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### Mocking Approach
- Use interfaces for dependencies that need to be mocked
- Create mock implementations manually or use tools like `gomock` or `testify/mock`
- Store mock implementations in the same package as tests

```go
// Mock implementation for testing
type mockUserRepository struct {
    users map[string]*User
}

func (m *mockUserRepository) FindByID(id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
}

func (m *mockUserRepository) Save(user *User) error {
    m.users[user.ID] = user
    return nil
}

// Usage in tests
func TestUserService_GetUser(t *testing.T) {
    repo := &mockUserRepository{
        users: map[string]*User{
            "123": {ID: "123", Name: "John Doe"},
        },
    }
    service := NewUserService(repo)
    
    user, err := service.GetUser("123")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "John Doe" {
        t.Errorf("expected name %q, got %q", "John Doe", user.Name)
    }
}
```

#### Coverage Expectations
- Aim for 80%+ code coverage for critical business logic
- 100% coverage is not always necessary; focus on meaningful tests
- Use `go test -cover` to check coverage
- Use `go test -coverprofile=coverage.out` and `go tool cover -html=coverage.out` for detailed coverage reports

### Concurrency

#### When to Use Goroutines
- Use goroutines for truly concurrent operations (I/O, network calls, independent computations)
- Don't use goroutines for simple sequential operations
- Always consider whether concurrency adds value vs. complexity

```go
// Good - concurrent I/O operations
func FetchMultipleUsers(ids []string) ([]*User, error) {
    results := make(chan *User, len(ids))
    errors := make(chan error, len(ids))
    
    for _, id := range ids {
        go func(id string) {
            user, err := FetchUser(id)
            if err != nil {
                errors <- err
                return
            }
            results <- user
        }(id)
    }
    
    // Collect results...
}

// Bad - unnecessary goroutine for simple operation
func Add(a, b int) int {
    result := make(chan int)
    go func() {
        result <- a + b
    }()
    return <-result  // Unnecessary overhead
}
```

#### Channel Patterns
- Buffered channels for producer-consumer patterns
- Unbuffered channels for synchronization
- Always close channels from the sender side
- Use select for multiplexing channel operations

```go
// Fan-out pattern
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output
        go func(out chan<- int) {
            defer close(out)
            for val := range input {
                out <- process(val)
            }
        }(output)
    }
    return outputs
}

// Select with timeout
func DoWithTimeout(ctx context.Context) error {
    result := make(chan error, 1)
    
    go func() {
        result <- doWork()
    }()
    
    select {
    case err := <-result:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

#### Context Usage Patterns
- Always pass context as the first parameter
- Never store context in a struct; pass it explicitly
- Use context for cancellation, deadlines, and request-scoped values
- Create child contexts with `context.WithTimeout`, `context.WithCancel`, etc.

```go
// Good
func (s *Service) ProcessRequest(ctx context.Context, req *Request) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return s.repo.Save(ctx, req)
}

// Bad
type Service struct {
    ctx context.Context  // Don't store context in structs
}
```

#### Synchronization Approach
- Prefer channels for communication between goroutines
- Use `sync.Mutex` for protecting shared state
- Use `sync.RWMutex` when reads are more frequent than writes
- Use `sync.WaitGroup` for waiting on multiple goroutines
- Use `sync.Once` for one-time initialization

```go
// Mutex for shared state
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// WaitGroup for coordinating goroutines
func ProcessItems(items []Item) {
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(item Item) {
            defer wg.Done()
            process(item)
        }(item)
    }
    
    wg.Wait()
}

// sync.Once for initialization
var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        instance = &Database{
            // Initialize...
        }
    })
    return instance
}
```

### Dependencies

#### How to Add New Dependencies
1. Use `go get` to add a new dependency: `go get github.com/pkg/errors`
2. Import the package in your code
3. Run `go mod tidy` to clean up the go.mod and go.sum files
4. Commit both go.mod and go.sum files

```bash
go get github.com/pkg/errors
go mod tidy
```

#### Preferred Libraries for Common Tasks
- **HTTP Router**: `net/http` (standard library) or `github.com/gorilla/mux` for more features
- **Logging**: `log/slog` (standard library, Go 1.21+) or `github.com/sirupsen/logrus`
- **Configuration**: `github.com/spf13/viper` or environment variables with `os.Getenv`
- **Database**: `database/sql` (standard library) with appropriate driver
- **ORM**: `gorm.io/gorm` or `github.com/jmoiron/sqlx`
- **Testing**: `testing` (standard library) with `github.com/stretchr/testify` for assertions
- **Mocking**: `github.com/golang/mock` or `github.com/stretchr/testify/mock`
- **HTTP Client**: `net/http` (standard library)
- **JSON**: `encoding/json` (standard library)
- **UUID**: `github.com/google/uuid`
- **Time**: `time` (standard library)

#### Dependency Injection Pattern
Use constructor functions to inject dependencies, promoting testability and loose coupling.

```go
// Define interfaces for dependencies
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

type NotificationService interface {
    SendEmail(ctx context.Context, email string, message string) error
}

// Service with dependencies
type UserService struct {
    repo         UserRepository
    notification NotificationService
    logger       *slog.Logger
}

// Constructor function for dependency injection
func NewUserService(
    repo UserRepository,
    notification NotificationService,
    logger *slog.Logger,
) *UserService {
    return &UserService{
        repo:         repo,
        notification: notification,
        logger:       logger,
    }
}

// Usage
func main() {
    logger := slog.Default()
    repo := NewPostgresUserRepository(db)
    notification := NewEmailService(config)
    
    service := NewUserService(repo, notification, logger)
    // Use service...
}
```

## Architecture Patterns

### Repository Pattern
Use the repository pattern to abstract data access, making the application database-agnostic and easier to test.

**When to use**: When you need to decouple business logic from data access logic.

```go
// Domain model (in internal/domain or internal/model)
package domain

type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Repository interface (in internal/domain or internal/repository)
package repository

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Save(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// Implementation (in internal/repository/postgres or internal/repository/mysql)
package postgres

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`
    
    var user domain.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, repository.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return &user, nil
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID,
        user.Email,
        user.Name,
        user.CreatedAt,
        user.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}
```

### Service Layer Pattern
Use the service layer to encapsulate business logic, keeping it separate from HTTP handlers and data access.

**When to use**: When you have complex business logic that doesn't belong in handlers or repositories.

```go
// Service interface (in internal/service)
package service

type UserService interface {
    CreateUser(ctx context.Context, email, name string) (*domain.User, error)
    GetUser(ctx context.Context, id string) (*domain.User, error)
    UpdateUser(ctx context.Context, id string, name string) (*domain.User, error)
    DeleteUser(ctx context.Context, id string) error
    ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// Implementation
type userService struct {
    repo   repository.UserRepository
    logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) UserService {
    return &userService{
        repo:   repo,
        logger: logger,
    }
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    // Validation
    if email == "" {
        return nil, errors.New("email is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }
    
    // Check if user already exists
    existing, err := s.repo.FindByEmail(ctx, email)
    if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    if existing != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // Create user
    user := &domain.User{
        ID:        uuid.New().String(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    s.logger.Info("user created", "user_id", user.ID, "email", email)
    
    return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### HTTP Handler Pattern
Keep HTTP handlers thin, delegating business logic to services.

```go
// Handler (in internal/handler or internal/api)
package handler

type UserHandler struct {
    service service.UserService
    logger  *slog.Logger
}

func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    user, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
    if err != nil {
        h.logger.Error("failed to create user", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to create user")
        return
    }
    
    h.respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        h.respondError(w, http.StatusBadRequest, "id parameter is required")
        return
    }
    
    user, err := h.service.GetUser(r.Context(), id)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound) {
            h.respondError(w, http.StatusNotFound, "user not found")
            return
        }
        h.logger.Error("failed to get user", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to get user")
        return
    }
    
    h.respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) respondError(w http.ResponseWriter, status int, message string) {
    h.respondJSON(w, status, map[string]string{"error": message})
}
```

### Configuration Pattern
Use a centralized configuration struct loaded from environment variables or config files.

```go
// Configuration (in internal/config)
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Port            int
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type LoggerConfig struct {
    Level  string
    Format string
}

func Load() (*Config, error) {
    return &Config{
        Server: ServerConfig{
            Port:            getEnvAsInt("SERVER_PORT", 8080),
            ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
            WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
            ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "myapp"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Logger: LoggerConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
        },
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

## Common Pitfalls to Avoid

### Ignoring Errors
**Problem**: Ignoring errors can lead to silent failures and hard-to-debug issues.

```go
// Bad
data, _ := ioutil.ReadFile("config.json")

// Good
data, err := ioutil.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}
```

### Returning Nil Pointers
**Problem**: Returning nil pointers without error can cause nil pointer panics.

```go
// Bad
func GetUser(id string) *User {
    // If not found, returns nil
    return nil
}

// Good
func GetUser(id string) (*User, error) {
    // Explicit error handling
    return nil, ErrUserNotFound
}
```

### Not Closing Resources
**Problem**: Not closing files, database connections, or HTTP response bodies leads to resource leaks.

```go
// Bad
func ReadFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    return ioutil.ReadAll(file)  // File never closed
}

// Good
func ReadFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    return ioutil.ReadAll(file)
}
```

### Goroutine Leaks
**Problem**: Starting goroutines without a way to stop them leads to resource leaks.

```go
// Bad
func StartWorker() {
    go func() {
        for {
            // Infinite loop with no way to stop
            doWork()
            time.Sleep(time.Second)
        }
    }()
}

// Good
func StartWorker(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return  // Proper cleanup
            case <-ticker.C:
                doWork()
            }
        }
    }()
}
```

### Mutating Shared State Without Synchronization
**Problem**: Concurrent access to shared state without proper synchronization causes race conditions.

```go
// Bad
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++  // Race condition
}

// Good
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

### Using Panic for Control Flow
**Problem**: Using panic for normal error conditions makes code unpredictable.

```go
// Bad
func GetUser(id string) *User {
    user, err := findUser(id)
    if err != nil {
        panic(err)  // Don't panic for expected errors
    }
    return user
}

// Good
func GetUser(id string) (*User, error) {
    user, err := findUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### Not Using Context for Cancellation
**Problem**: Long-running operations without context can't be cancelled or timed out.

```go
// Bad
func ProcessRequest(req *Request) error {
    // No way to cancel or timeout
    result := longRunningOperation(req)
    return result
}

// Good
func ProcessRequest(ctx context.Context, req *Request) error {
    // Can be cancelled or timed out
    return longRunningOperation(ctx, req)
}
```

### Pointer Receivers on Methods That Don't Modify State
**Problem**: Using pointer receivers unnecessarily can prevent certain compiler optimizations.

```go
// Bad - pointer receiver when value receiver would work
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// Good - value receiver for methods that don't modify state
func (u User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// Good - pointer receiver when modifying state
func (u *User) SetName(first, last string) {
    u.FirstName = first
    u.LastName = last
}
```

## Documentation Requirements

### Public APIs
All exported (capitalized) functions, types, methods, and constants must have GoDoc comments.

```go
// User represents a user in the system.
// It contains the basic information needed to identify and contact a user.
type User struct {
    // ID is the unique identifier for the user
    ID string
    // Email is the user's email address
    Email string
    // Name is the user's full name
    Name string
}

// NewUser creates a new User with the given email and name.
// It generates a unique ID for the user automatically.
//
// Example:
//   user := NewUser("john@example.com", "John Doe")
func NewUser(email, name string) *User {
    return &User{
        ID:    uuid.New().String(),
        Email: email,
        Name:  name,
    }
}

// Save persists the user to the database.
// It returns an error if the user cannot be saved.
func (u *User) Save(ctx context.Context) error {
    // Implementation...
}
```

### Complex Logic
Add inline comments for complex algorithms or non-obvious code.

```go
func calculateScore(values []int) int {
    // Use a sliding window approach to find the maximum sum
    // of any contiguous subarray of size 3
    if len(values) < 3 {
        return 0
    }
    
    maxSum := 0
    currentSum := values[0] + values[1] + values[2]
    maxSum = currentSum
    
    // Slide the window through the array
    for i := 3; i < len(values); i++ {
        currentSum = currentSum - values[i-3] + values[i]
        if currentSum > maxSum {
            maxSum = currentSum
        }
    }
    
    return maxSum
}
```

### Package Documentation
Use `doc.go` files for package-level documentation.

```go
// Package user provides functionality for managing user accounts.
//
// This package implements user creation, authentication, and profile management.
// It uses bcrypt for password hashing and supports email verification.
//
// Basic usage:
//
//   service := user.NewService(repo)
//   user, err := service.CreateUser(ctx, "john@example.com", "password123")
//   if err != nil {
//       // Handle error
//   }
//
package user
```

## Before Submitting Code

Before submitting your code, ensure you've completed these steps:

- [ ] Run `go fmt ./...` to format all code
- [ ] Run `go vet ./...` to check for common mistakes
- [ ] Run `go test ./...` to ensure all tests pass
- [ ] Run `go test -race ./...` to check for race conditions
- [ ] Run `golangci-lint run` if available for comprehensive linting
- [ ] Run `go mod tidy` to clean up dependencies
- [ ] Check test coverage with `go test -cover ./...`
- [ ] Update relevant documentation (README, GoDoc comments)
- [ ] Ensure commit messages are clear and descriptive
- [ ] Review your own code changes before requesting review

### Additional Checks for Production Code
- [ ] Check for security vulnerabilities with `gosec` or similar tools
- [ ] Verify no sensitive data (passwords, keys) is hardcoded
- [ ] Ensure proper error handling and logging
- [ ] Verify graceful shutdown for services
- [ ] Test with realistic data volumes
- [ ] Review performance implications of changes
- [ ] Ensure backward compatibility for public APIs

## Additional Guidelines

### Struct Tags
Use struct tags consistently for JSON, database, and validation tags.

```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Email     string    `json:"email" db:"email" validate:"required,email"`
    Name      string    `json:"name" db:"name" validate:"required"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### Zero Values
Design types so that their zero value is useful.

```go
// Good - zero value is useful
type Buffer struct {
    buf []byte
}

// Can be used without initialization
var b Buffer
b.Write([]byte("hello"))

// Bad - requires initialization
type BadBuffer struct {
    buf []byte
    initialized bool  // Unnecessary flag
}
```

### Prefer Composition Over Inheritance
Go doesn't have inheritance; use composition and interfaces instead.

```go
// Good - composition
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}

// Good - embedding
type BufferedReader struct {
    reader Reader
    buffer []byte
}
```

### Use Make for Slices and Maps
Always use `make` for slices and maps when you know the size.

```go
// Good
users := make([]*User, 0, 100)  // Pre-allocate capacity
cache := make(map[string]*User, 100)

// Less efficient
users := []*User{}  // No pre-allocation
cache := map[string]*User{}
```

### String Building
Use `strings.Builder` for efficient string concatenation in loops.

```go
// Bad
var result string
for _, s := range strings {
    result += s  // Creates new string each iteration
}

// Good
var builder strings.Builder
for _, s := range strings {
    builder.WriteString(s)
}
result := builder.String()
```

This guide should be updated as the project evolves and new patterns emerge. Always prioritize code clarity, maintainability, and idiomatic Go over clever solutions.
