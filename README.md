# go-pkg-utils

A comprehensive Go utilities package providing essential functionality for modern Go applications including configuration management, HTTP responses, string manipulation, cryptography, collections, validation, error handling, structured logging, pagination, message queues, Lua scripting, and more.

## 🚀 Features

- **🔧 Configuration**: Environment variable parsing with type safety, YAML configuration loading with environment variable substitution, and comprehensive validation
- **🌐 HTTP Responses**: Standardized API response structures with comprehensive HTTP status code support and Fiber integration
- **📝 Text Processing**: String manipulation, case conversion, validation, and more
- **🔐 Cryptography**: Secure random generation, password hashing, AES/RSA encryption, JWT
- **🔒 HMAC Authentication**: HMAC-authenticated HTTP client for secure service-to-service communication
- **📊 Collections**: Generic utilities for slices and maps with functional programming support
- **📅 Date/Time**: Comprehensive time manipulation and formatting utilities
- **🔄 JSON**: Advanced JSON marshaling, unmarshaling, and manipulation
- **✅ Validation**: Struct validation with tags and detailed error reporting
- **🚨 Error Handling**: Structured error system with metadata and stack traces
- **💬 Messages**: Predefined success and error messages
- **🌐 Network**: IP address extraction and validation utilities
- **🆔 UUID**: UUID utilities and validation
- **📋 Logging**: Structured logging with Zap, file rotation, and Fiber middleware integration
- **📄 Pagination**: GORM-based pagination utilities with Fiber integration and query parameter parsing
- **🔍 Filtering**: Unified query filter system with field operators (eq, gt, gte, lt, lte, like, in, not_in) for reusable filtering across microservices
- **📨 Queue**: RabbitMQ producer/consumer with automatic reconnection, retry logic, and dead letter queues
- **🎯 Lua Scripting**: Configurable sandboxed Lua script execution with worker pools, timeout handling, and result recording

## 📦 Installation

```bash
go get github.com/kerimovok/go-pkg-utils
```

## 🏗️ Package Structure

```
go-pkg-utils/
├── collections/     # Generic slice and map utilities
├── config/         # Environment variable and validation utilities
├── crypto/         # Cryptographic functions and password hashing
├── datetime/       # Time and date manipulation utilities
├── errors/         # Structured error handling system
├── httpx/          # HTTP response utilities for Fiber
├── hmac/           # HMAC-authenticated HTTP client for service-to-service communication
├── jsonx/          # Advanced JSON processing utilities
├── logger/         # Structured logging with Zap and Fiber middleware
├── lua/            # Configurable sandboxed Lua script execution and worker pools
├── messages/       # Predefined message constants
├── filter/        # Unified query filter system with operators for GORM queries
├── net/           # Network utilities (package netx; IP extraction)
├── pagination/    # GORM pagination utilities with Fiber integration
├── queue/         # RabbitMQ producer/consumer with retry and DLQ support
│   ├── events/    # Event producer for direct exchange routing
│   └── tasks/     # Task producer for topic exchange with wildcard routing
├── text/          # String manipulation and processing
├── uuid/          # UUID utilities (package uuidx)
└── validator/     # Configuration and struct validation
```

## 📖 Quick Start

### Configuration Management

```go
import "github.com/kerimovok/go-pkg-utils/config"

// Environment variables with type safety
port := config.GetEnvInt("PORT", 8080)
debug := config.GetEnvBool("DEBUG", false)
timeout := config.GetEnvDuration("TIMEOUT", 30*time.Second)

// Required environment variables
dbURL := config.RequiredEnv("DATABASE_URL") // Panics if not set

// YAML Configuration with Environment Variable Substitution
type AppConfig struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
}

var cfg AppConfig
err := config.LoadYAMLConfig("config.yaml", &cfg)

// Environment variable substitution in YAML files
// config.yaml example:
// database:
//   host: ${DB_HOST:-localhost}  // or ${DB_HOST=localhost}
//   port: ${DB_PORT:-5432}       // or ${DB_PORT=5432}
//   username: ${DB_USERNAME}
//   password: ${DB_PASSWORD}
// server:
//   port: ${SERVER_PORT:-8080}   // or ${SERVER_PORT=8080}
//   host: ${SERVER_HOST:-0.0.0.0} // or ${SERVER_HOST=0.0.0.0}

// Network and URL validation
if !config.IsValidPort("8080") {
    // Handle invalid port
}
if !config.IsValidIP("192.168.1.1") {
    // Handle invalid IP
}
if !config.IsValidEmail("user@example.com") {
    // Handle invalid email
}
if !config.IsValidURL("https://example.com") {
    // Handle invalid URL
}
if !config.IsValidDomain("example.com") {
    // Handle invalid domain
}
if !config.IsValidHost("api.example.com") {
    // Handle invalid host
}

// Number validation
if !config.IsValidInteger("123") {
    // Handle invalid integer
}
if !config.IsValidPositiveInteger("456") {
    // Handle invalid positive integer
}
if !config.IsValidFloat("123.45") {
    // Handle invalid float
}
if !config.IsValidNumber("123.45") {
    // Handle invalid number
}

// String validation
if !config.IsValidNonEmptyString("hello") {
    // Handle empty string
}
if !config.IsValidAlphabetic("HelloWorld") {
    // Handle non-alphabetic string
}
if !config.IsValidAlphanumeric("Hello123") {
    // Handle non-alphanumeric string
}
if !config.IsValidLowercase("hello") {
    // Handle non-lowercase string
}
if !config.IsValidUppercase("HELLO") {
    // Handle non-uppercase string
}

// Format validation
if !config.IsValidUUID("550e8400-e29b-41d4-a716-446655440000") {
    // Handle invalid UUID
}
if !config.IsValidBase64("SGVsbG8gV29ybGQ=") {
    // Handle invalid base64
}
if !config.IsValidHex("1a2b3c4d") {
    // Handle invalid hex
}
if !config.IsValidSlug("hello-world") {
    // Handle invalid slug
}
```

### HTTP Responses with Fiber

```go
import "github.com/kerimovok/go-pkg-utils/httpx"

func handler(c *fiber.Ctx) error {
    data := map[string]string{"message": "Hello World"}

    // Success response
    response := httpx.OK("Data retrieved successfully", data)
    return httpx.SendResponse(c, response)

    // Error response
    err := errors.New("something went wrong")
    errorResponse := httpx.BadRequest("Invalid request", err)
    return httpx.SendResponse(c, errorResponse)

    // Paginated response
    pagination := httpx.NewPagination(1, 10, 100)
    paginatedResponse := httpx.Paginated("Users retrieved", users, pagination)
    return httpx.SendPaginatedResponse(c, paginatedResponse)
}
```

### String Manipulation

```go
import "github.com/kerimovok/go-pkg-utils/text"

// Case conversion
snake := text.ToSnakeCase("HelloWorld")     // "hello_world"
camel := text.ToCamelCase("hello_world")    // "helloWorld"
pascal := text.ToPascalCase("hello_world")  // "HelloWorld"
kebab := text.ToKebabCase("HelloWorld")     // "hello-world"

// Text processing
slug := text.ToSlug("Hello World!")        // "hello-world"
truncated := text.TruncateWithEllipsis("Long text", 10) // "Long te..."
reversed := text.Reverse("hello")          // "olleh"

// Extraction
emails := text.ExtractEmails("Contact us at: admin@example.com or support@test.com")
urls := text.ExtractURLs("Visit https://example.com and https://github.com")

// Masking
masked := text.MaskEmail("user@example.com") // "u***@example.com"
```

### Cryptography

```go
import "github.com/kerimovok/go-pkg-utils/crypto"

// Password hashing (bcrypt)
hash, err := crypto.HashPassword("mypassword")
isValid := crypto.CheckPassword("mypassword", hash)

// Secure password hashing (scrypt)
hash, salt, err := crypto.HashPasswordSecure("mypassword")
isValid, err := crypto.VerifyPasswordSecure("mypassword", hash, salt)

// Random generation
token, err := crypto.GenerateToken(32)
apiKey, err := crypto.GenerateAPIKey()

// AES encryption
key, _ := crypto.GenerateSecretKey()
encrypted, err := crypto.EncryptString("sensitive data", key)
decrypted, err := crypto.DecryptString(encrypted, key)

// RSA encryption
privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
encrypted, err := crypto.RSAEncrypt([]byte("data"), publicKey)
decrypted, err := crypto.RSADecrypt(encrypted, privateKey)

// JWT (simple implementation)
jwt := crypto.NewSimpleJWT([]byte("secret"))
claims := crypto.JWTClaims{
    Subject: "user123",
    Custom: map[string]interface{}{"role": "admin"},
}
token, err := jwt.CreateToken(claims)
parsedClaims, err := jwt.VerifyToken(token)
```

### HMAC Authentication

The HMAC package provides a secure HTTP client for service-to-service communication using HMAC-SHA256 signatures. The signature includes the HTTP method, path, query string, timestamp, and request body to prevent request tampering.

```go
import "github.com/kerimovok/go-pkg-utils/hmac"

// Create HMAC client
client := hmac.NewClient(hmac.Config{
    BaseURL:    "https://api.example.com",
    HMACSecret: "your-secret-key",
    Timeout:    10 * time.Second, // Optional, defaults to 10s
})

// Make authenticated request with JSON body
type RequestBody struct {
    UserID string `json:"user_id"`
    Action string `json:"action"`
}

body := RequestBody{
    UserID: "123",
    Action: "verify",
}

resp, err := client.DoRequest("POST", "/api/v1/users/verify", body)
if err != nil {
    log.Fatal("Request failed:", err)
}
defer resp.Body.Close()

// Parse JSON response
var result map[string]interface{}
if err := hmac.ParseJSONResponse(resp, &result); err != nil {
    log.Fatal("Failed to parse response:", err)
}

// Make request with raw body bytes
bodyBytes := []byte(`{"key": "value"}`)
resp, err := client.DoRequestWithBody("POST", "/api/v1/data", bodyBytes)

// Compute signature manually (for server-side validation)
signature := hmac.ComputeSignature(
    "POST",
    "/api/v1/users",
    "status=active",
    "1234567890",
    []byte(`{"id": "123"}`),
    "secret-key",
)

// Validate signature (for server-side middleware)
isValid := hmac.ValidateSignature(
    "POST",
    "/api/v1/users",
    "status=active",
    "1234567890",
    []byte(`{"id": "123"}`),
    "computed-signature",
    "secret-key",
)
```

#### Signature Format

The HMAC signature is computed as:

```
HMAC-SHA256(method + path + query + timestamp + body, secret)
```

**Components:**

- `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
- `path`: Request path (e.g., `/api/v1/users`)
- `query`: Query string without `?` prefix (e.g., `status=active&per_page=10`)
- `timestamp`: Unix timestamp as string (e.g., `"1234567890"`)
- `body`: Raw request body bytes

**Headers:**

- `X-Signature`: The computed HMAC-SHA256 signature (hex-encoded)
- `X-Timestamp`: Unix timestamp of the request

**Security Features:**

- Constant-time signature comparison to prevent timing attacks
- Timestamp included to prevent replay attacks
- All request components included in signature to prevent tampering
- Automatic JSON marshaling for request bodies

### Collections

```go
import "github.com/kerimovok/go-pkg-utils/collections"

// Slice utilities
numbers := []int{1, 2, 3, 4, 5}
contains := collections.Contains(numbers, 3)     // true
filtered := collections.Filter(numbers, func(n int) bool { return n > 3 }) // [4, 5]
doubled := collections.Map(numbers, func(n int) int { return n * 2 })       // [2, 4, 6, 8, 10]
sum := collections.Reduce(numbers, 0, func(acc, n int) int { return acc + n }) // 15

// Map utilities
data := map[string]int{"a": 1, "b": 2, "c": 3}
keys := collections.Keys(data)                   // ["a", "b", "c"]
values := collections.Values(data)               // [1, 2, 3]
filtered := collections.FilterMap(data, func(k string, v int) bool { return v > 1 })
```

### Date/Time Utilities

```go
import "github.com/kerimovok/go-pkg-utils/datetime"

// Current time helpers
now := datetime.Now()
today := datetime.Today()
tomorrow := datetime.Tomorrow()

// Period boundaries
startOfWeek := datetime.StartOfWeek(now)
endOfMonth := datetime.EndOfMonth(now)
startOfYear := datetime.StartOfYear(now)

// Calculations
age := datetime.Age(birthDate)
daysBetween := datetime.DaysBetween(start, end)
businessDays := datetime.BusinessDaysBetween(start, end)

// Formatting
timeAgo := datetime.TimeAgo(pastTime)           // "2 hours ago"
timeUntil := datetime.TimeUntil(futureTime)     // "in 3 days"
formatted := datetime.FormatDuration(duration)  // "2h30m"

// Parsing
parsed, err := datetime.ParseDate("2023-12-25")
```

#### UTC-normalized API responses

The `datetime` package also provides helpers to **canonicalize all timestamps to UTC**
before JSON serialization, which is especially useful when your database uses
`timestamp without time zone` but your API contract is “always UTC”.

```go
import "github.com/kerimovok/go-pkg-utils/datetime"

// Example response DTO with time.Time fields
type SessionResponse struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"userId"`
    ExpiresAt time.Time `json:"expiresAt"`
    CreatedAt time.Time `json:"createdAt"`
    RevokedAt *time.Time `json:"revokedAt,omitempty"`
}

// NormalizeTimeFieldsToUTC walks the struct (and nested slices/structs)
// and converts all time.Time fields to UTC in-place.
func ListSessions(c *fiber.Ctx) error {
    sessions := make([]SessionResponse, 0)
    // ... load sessions into the slice ...

    datetime.NormalizeTimeFieldsToUTC(&sessions)

    resp := httpx.OK("Sessions retrieved successfully", sessions)
    return httpx.SendResponse(c, resp)
}

// Helpers for parsing / formatting RFC3339 timestamps in UTC:
t, err := datetime.ParseRFC3339ToUTC("2026-02-09T12:12:35.472122Z")
utc := datetime.ToUTC(t)                     // no-op if already UTC
iso := datetime.FormatRFC3339UTC(utc)        // always RFC3339 with Z suffix
excel := datetime.FormatExcelUTC(utc)        // "2006-01-02 15:04:05" in UTC
```

### JSON Utilities

```go
import "github.com/kerimovok/go-pkg-utils/jsonx"

data := map[string]interface{}{
    "user": map[string]interface{}{
        "name": "John",
        "age": 30,
    },
}

// Path-based access
name, err := jsonx.GetValue(data, "user.name")  // "John"
err = jsonx.SetValue(data, "user.email", "john@example.com")
err = jsonx.DeleteValue(data, "user.age")

// Flattening
flat := jsonx.Flatten(data)                     // {"user.name": "John", "user.email": "john@example.com"}
nested := jsonx.Unflatten(flat)                 // Original nested structure

// Type-safe access
userMap, err := jsonx.GetObject(data, "user")
name, err := jsonx.GetString(userMap, "name")
age, err := jsonx.GetInt(userMap, "age")
```

### Validation

```go
import "github.com/kerimovok/go-pkg-utils/validator"
import "github.com/kerimovok/go-pkg-utils/config"

type User struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=18,max=100"`
}

user := User{Name: "John", Email: "invalid-email", Age: 15}
errors := validator.ValidateStruct(user)

if errors.HasErrors() {
    for _, err := range errors {
        fmt.Printf("Field: %s, Error: %s\n", err.Field, err.Message)
    }
}

// Direct validation using config functions
if !config.IsValidEmail(user.Email) {
    fmt.Println("Invalid email format")
}

if !config.IsValidNonEmptyString(user.Name) {
    fmt.Println("Name cannot be empty")
}

// Number validation
ageStr := "25"
if !config.IsValidPositiveInteger(ageStr) {
    fmt.Println("Age must be a positive integer")
}

// Network validation
if !config.IsValidIP("192.168.1.1") {
    fmt.Println("Invalid IP address")
}

if !config.IsValidPort("8080") {
    fmt.Println("Invalid port number")
}

// Format validation
if !config.IsValidUUID("550e8400-e29b-41d4-a716-446655440000") {
    fmt.Println("Invalid UUID format")
}

if !config.IsValidSlug("user-profile") {
    fmt.Println("Invalid slug format")
}
```

### Error Handling

```go
import "github.com/kerimovok/go-pkg-utils/errors"

// Create structured errors
err := errors.ValidationError("INVALID_EMAIL", "Email format is invalid").
    WithMetadata("field", "email").
    WithUserID("user123").
    WithRequestID("req456")

// Error chain for multiple errors
chain := errors.NewErrorChain()
chain.Add(errors.ValidationError("REQUIRED", "Name is required"))
chain.Add(errors.ValidationError("INVALID", "Email is invalid"))

// Error handler with panic recovery
handler := errors.NewErrorHandler("user-service", func(err error) {
    log.Printf("Error: %s", err.Error())
})

err := handler.SafeExecute(func() error {
    // Your code that might panic
    return nil
})
```

### Logging

```go
import "github.com/kerimovok/go-pkg-utils/logger"
import "go.uber.org/zap"

// Create logger from configuration
config := &logger.Config{
    Enabled:    func() *bool { b := true; return &b }(),
    FilePath:   "/var/log/app.log",
    MaxSize:    104857600, // 100MB
    MaxBackups: 3,
    MaxAge:     28, // days
    Level:      "info",
}

log, err := logger.NewLogger(config)
if err != nil {
    panic(err)
}
defer log.Sync()

// Use logger
log.Info("Application started", zap.String("version", "1.0.0"))
log.Error("Something went wrong", zap.Error(err))

// Setup Fiber middleware
app := fiber.New()
loggerInstance, err := logger.SetupFiberLogger(app, config)
if err != nil {
    log.Fatal("Failed to setup logger", zap.Error(err))
}
defer loggerInstance.Sync()

// Development logger (console output, colored)
devLogger, err := logger.NewDevelopmentLogger()

// Production logger (JSON output, file rotation)
prodLogger, err := logger.NewProductionLogger("/var/log/app.log", 100, 3, 28)
```

### Network and UUID Utilities

```go
import (
    netx "github.com/kerimovok/go-pkg-utils/net"
    uuidx "github.com/kerimovok/go-pkg-utils/uuid"
)

// Get client IP from Fiber context (CF, X-Forwarded-For, X-Real-IP)
ip := netx.GetUserIP(c)

// Parse UUID (format validation is available via config.IsValidUUID)
id, err := uuidx.Parse("550e8400-e29b-41d4-a716-446655440000")
```

### Pagination

```go
import "github.com/kerimovok/go-pkg-utils/pagination"
import "gorm.io/gorm"

// Define your model
type User struct {
    ID        uint
    Name      string
    Email     string
    CreatedAt time.Time
}

// Simple pagination handler
func GetUsers(c *fiber.Ctx, db *gorm.DB) error {
    defaults := pagination.Default() // Page: 1, PerPage: 20, SortBy: "created_at", SortOrder: "desc"

    // Customize defaults if needed
    defaults.SortBy = "name"
    defaults.PerPage = 10

    return pagination.HandleRequest[User](c, db.Model(&User{}), defaults, "Users retrieved successfully")
}

// Manual pagination with custom query
func GetFilteredUsers(c *fiber.Ctx, db *gorm.DB) error {
    defaults := pagination.Default()
    params, err := pagination.ParseParams(c, defaults)
    if err != nil {
        response := httpx.BadRequest("Invalid pagination parameters", err)
        return httpx.SendResponse(c, response)
    }

    // Build custom query
    query := db.Model(&User{}).Where("active = ?", true)

    // Execute paginated query
    ctx := c.Context()
    response, err := pagination.Query[User](ctx, query, params, "Active users retrieved")
    if err != nil {
        response := httpx.InternalServerError("Failed to retrieve users", err)
        return httpx.SendResponse(c, response)
    }

    return httpx.SendPaginatedResponse(c, *response)
}
```

### Filtering

The filter package provides a unified query filtering system that can be reused across microservices. It supports various operators and automatically handles type conversion.

```go
import "github.com/kerimovok/go-pkg-utils/filter"
import "gorm.io/gorm"

// Define your model
type QRCode struct {
    ID        uuid.UUID
    Data      string
    Status    string
    Size      int
    CreatedAt time.Time
}

// Simple filtering with defaults
func ListQRCodes(c *fiber.Ctx, db *gorm.DB) error {
    // Configure allowed fields and their types
    filterConfig := &filter.Config{
        AllowedFields: map[string]string{
            "status":     "string",
            "created_at": "time",
            "size":       "int",
            "data":       "string",
        },
        CustomValidators: map[string]func(value string) error{
            "status": func(value string) error {
                allowed := map[string]bool{"active": true, "inactive": true, "archived": true}
                if !allowed[value] {
                    return fmt.Errorf("invalid status: %s", value)
                }
                return nil
            },
        },
    }

    // Build query
    query := db.Model(&QRCode{})

    // Apply filters from query parameters
    // Format: field_operator=value
    // Examples:
    //   ?status_eq=active
    //   ?created_at_gte=2024-01-01
    //   ?size_gt=300
    //   ?data_like=example
    //   ?status_in=active,inactive
    var err error
    query, err = filter.ApplyFiltersFromContext(c, query, filterConfig)
    if err != nil {
        response := httpx.BadRequest("Invalid filter parameters", err)
        return httpx.SendResponse(c, response)
    }

    // Use with pagination
    defaults := pagination.Default()
    return pagination.HandleRequest[QRCode](c, query, defaults, "QR codes retrieved successfully")
}
```

#### Supported Operators

- `eq` - Equals: `?status_eq=active`
- `ne` - Not equals: `?status_ne=deleted`
- `gt` - Greater than: `?size_gt=300`
- `gte` - Greater than or equal: `?created_at_gte=2024-01-01`
- `lt` - Less than: `?size_lt=500`
- `lte` - Less than or equal: `?created_at_lte=2024-12-31`
- `like` - Like (for strings): `?data_like=example` (adds `%` wildcards automatically)
- `in` - In (for arrays): `?status_in=active,inactive,archived`
- `not_in` - Not in (for arrays): `?status_not_in=deleted,archived`

#### Query Parameter Format

Filters use the format: `field_operator=value`

**Examples:**

```
GET /api/v1/qrcodes?status_eq=active&created_at_gte=2024-01-01&size_gt=300
GET /api/v1/qrcodes?status_in=active,inactive&data_like=example
GET /api/v1/qrcodes?created_at_lte=2024-12-31T23:59:59Z
```

#### Field Type Conversion

The filter automatically converts values based on the field type specified in `AllowedFields`:

- `string` - Default, no conversion
- `int` or `integer` - Converts to integer
- `float` or `float64` - Converts to float
- `bool` or `boolean` - Converts to boolean (`true`, `1` = true)
- `time`, `datetime`, or `date` - Parses time (supports RFC3339, `2006-01-02`, `2006-01-02T15:04:05`)

#### Advanced Usage

```go
// Field mapping (map query field names to database column names)
filterConfig := &filter.Config{
    FieldMapping: map[string]string{
        "uploadedAfter":  "created_at",  // Query: uploadedAfter_gte, DB: created_at
        "uploadedBefore": "created_at", // Query: uploadedBefore_lte, DB: created_at
    },
    AllowedFields: map[string]string{
        "uploadedAfter":  "time",
        "uploadedBefore": "time",
    },
}

// Manual filter parsing and application
filters, err := filter.ParseFilters(c, filterConfig)
if err != nil {
    return httpx.BadRequest("Invalid filters", err)
}

query := filter.ApplyFilters(db.Model(&Model{}), filters)

// Combine with custom query conditions
query = query.Where("deleted_at IS NULL")
query = filter.ApplyFilters(query, filters)
```

#### Integration with Pagination

Filters work seamlessly with the pagination package:

```go
func ListResources(c *fiber.Ctx, db *gorm.DB) error {
    // Configure filters
    filterConfig := &filter.Config{
        AllowedFields: map[string]string{
            "status":     "string",
            "created_at": "time",
        },
    }

    // Build query with filters
    query := db.Model(&Resource{})
    query, err := filter.ApplyFiltersFromContext(c, query, filterConfig)
    if err != nil {
        return httpx.BadRequest("Invalid filters", err)
    }

    // Apply pagination
    defaults := pagination.Default()
    return pagination.HandleRequest[Resource](c, query, defaults, "Resources retrieved")
}
```

**Note**: The filter package automatically skips reserved pagination parameters (`page`, `per_page`, `sort_by`, `sort_order`).

### Queue (RabbitMQ)

```go
import "github.com/kerimovok/go-pkg-utils/queue"
import "github.com/kerimovok/go-pkg-utils/queue/events"
import "github.com/kerimovok/go-pkg-utils/queue/tasks"

// Producer setup
connConfig := queue.ConnectionConfig{
    Host:     "localhost",
    Port:     "5672",
    Username: "guest",
    Password: "guest",
    VHost:    "/",
}

queueConfig := &queue.Config{
    ExchangeName:    "my_exchange",
    QueueName:       "my_queue",
    RoutingKey:      "my.routing.key",
    DLXExchangeName: "my_exchange.dlx",
    DLQName:         "my_dlq",
    DLQRoutingKey:   "my.failed",
}

producer, err := queue.NewProducer(connConfig, queueConfig)
if err != nil {
    log.Fatal("Failed to create producer:", err)
}
defer producer.Close()

// Publish message
headers := amqp.Table{
    "x-custom-header": "value",
}
err = producer.Publish(ctx, []byte(`{"message": "hello"}`), headers)

// Consumer setup with retry
retryConfig := queue.RetryConfig{
    MaxRetries:     3,
    RetryDelayBase: 1, // seconds
    MaxRetryDelay:  60, // seconds
}

handler := func(msg amqp.Delivery) error {
    // Process message
    var data map[string]interface{}
    if err := json.Unmarshal(msg.Body, &data); err != nil {
        return err // Will trigger retry
    }

    // Process data...
    return nil // Success - message will be acknowledged
}

consumer, err := queue.NewConsumer(connConfig, queueConfig, retryConfig, handler)
if err != nil {
    log.Fatal("Failed to create consumer:", err)
}
defer consumer.Close()

// Start consuming
if err := consumer.StartConsuming(); err != nil {
    log.Fatal("Failed to start consuming:", err)
}

// Event Producer (simplified event publishing)
// Uses direct exchange with exact routing key match
eventProducer, err := events.NewProducer(connConfig, events.ProducerConfig{
    ServiceName: "user-service",
})
if err != nil {
    log.Fatal("Failed to create event producer:", err)
}
defer eventProducer.Close()

// Publish event (routing key: "event")
err = eventProducer.Publish(ctx, "user.created", map[string]any{
    "user_id": "123",
    "email":   "user@example.com",
})

// Publish asynchronously (fire and forget)
eventProducer.PublishAsync("user.updated", map[string]any{
    "user_id": "123",
})

// Task Producer (simplified task publishing)
// Uses topic exchange with wildcard routing (tasks.<taskType>)
taskProducer, err := tasks.NewProducer(connConfig, tasks.ProducerConfig{
    ServiceName: "auth-service",
})
if err != nil {
    log.Fatal("Failed to create task producer:", err)
}
defer taskProducer.Close()

// Publish task (routing key automatically: "tasks.email.verify")
err = taskProducer.Publish(ctx, "email.verify", map[string]any{
    "to":       "user@example.com",
    "template": "verify-email",
    "data": map[string]any{
        "token": "abc123",
    },
})

// Publish asynchronously (fire and forget)
taskProducer.PublishAsync("email.send", map[string]any{
    "to":      "user@example.com",
    "subject": "Welcome",
})

// Publish with custom routing key (override default pattern)
err = taskProducer.PublishWithCustomRoutingKey(ctx, "email.verify", payload, "tasks.custom.route")
```

#### Event vs Task Producers

**Events Producer** (`queue/events`):

- **Exchange Type**: `direct` (exact routing key match)
- **Routing Key**: `"event"` (fixed)
- **Use Case**: Event-driven architecture, event sourcing, audit logs
- **Message Structure**: `{service, type, payload}`

**Tasks Producer** (`queue/tasks`):

- **Exchange Type**: `topic` (wildcard routing support)
- **Routing Key**: `"tasks.<taskType>"` (auto-constructed from task type)
- **Use Case**: Task queues, job processing, async operations
- **Message Structure**: `{service, type, payload}`
- **Example**: Task type `"email.verify"` → routing key `"tasks.email.verify"`

Both producers automatically:

- Add timestamps to payloads
- Handle connection management and reconnection
- Support async publishing (fire and forget)
- Use consistent message structure with service, type, and payload

### Lua Scripting

```go
import "github.com/kerimovok/go-pkg-utils/lua"

// Define a script
type MyScript struct {
    ID      string
    Name    string
    Version string
    Code    string
}

func (s *MyScript) GetID() string   { return s.ID }
func (s *MyScript) GetName() string { return s.Name }
func (s *MyScript) GetVersion() string { return s.Version }
func (s *MyScript) GetCode() string { return s.Code }

// Create executor with default strict sandboxing
executor := lua.NewExecutor(lua.ExecutorConfig{
    Timeout: 5 * time.Second,
    Logger:  zapLogger, // optional
    HostFunctions: customFunctionRegistry, // optional
    Recorder: executionRecorder, // optional
    // Sandbox: nil // Uses DefaultSandboxConfig() - strict security
})

// Create executor with custom sandbox configuration
customSandbox := lua.SandboxConfig{
    EnableBase:   true,
    EnableTable:  true,
    EnableString: true,
    EnableMath:   true,
    EnableOS:     false, // Keep disabled for security
    EnableIO:     false, // Keep disabled for security
    EnableDebug:  false, // Keep disabled for security
    DisableDofile:     true,
    DisableLoadfile:   true,
    DisableLoad:       true,
    DisableLoadstring: true,
}

executorWithCustomSandbox := lua.NewExecutor(lua.ExecutorConfig{
    Timeout: 5 * time.Second,
    Logger:  zapLogger,
    Sandbox: &customSandbox, // Custom sandbox configuration
})

// Or use the default strict configuration explicitly
strictSandbox := lua.DefaultSandboxConfig()
executorStrict := lua.NewExecutor(lua.ExecutorConfig{
    Timeout: 5 * time.Second,
    Sandbox: &strictSandbox,
})

// Execute script
script := &MyScript{
    ID:      "script-1",
    Name:    "process_data",
    Version: "1.0.0",
    Code: `
        function handle(payload)
            local result = payload.value * 2
            print("Result: " .. result)
        end
    `,
}

payload := map[string]interface{}{
    "value": 42,
}

result := executor.Execute(ctx, script, payload)
if result.Status == lua.ExecutionStatusFailure {
    log.Printf("Script failed: %s", *result.ErrorMessage)
}

// Worker pool for concurrent execution
pool := lua.NewWorkerPool(10) // Max 10 concurrent executions

go func() {
    pool.Acquire()
    defer pool.Release()

    result := executor.Execute(ctx, script, payload)
    // Process result...
}()

// Create VM directly with custom sandbox configuration
vm := lua.NewVM(lua.SandboxConfig{
    EnableBase:   true,
    EnableTable:  true,
    EnableString: true,
    EnableMath:   true,
    // ... other options
})
defer vm.Close()

// Or use the backward-compatible function (uses strict defaults)
vm := lua.NewSandboxedVM()
defer vm.Close()
```

#### Sandbox Configuration

The sandbox configuration allows you to control which Lua libraries and functions are available to scripts:

```go
// Default strict configuration (recommended)
config := lua.DefaultSandboxConfig()
// Only base, table, string, math libraries enabled
// Dangerous functions (dofile, loadfile, load, loadstring) disabled

// Custom configuration
config := lua.SandboxConfig{
    // Libraries
    EnableBase:   true,  // Basic functions (print, type, etc.)
    EnableTable:  true,  // Table manipulation
    EnableString: true,  // String manipulation
    EnableMath:   true,  // Math functions
    EnableOS:     false, // OS functions (security risk)
    EnableIO:     false, // IO functions (security risk)
    EnableDebug:  false, // Debug functions (security risk)

    // Dangerous functions to disable
    DisableDofile:     true, // Disable dofile()
    DisableLoadfile:   true, // Disable loadfile()
    DisableLoad:       true, // Disable load()
    DisableLoadstring: true, // Disable loadstring()
}
```

## 🔧 Configuration

### Comprehensive Validation Functions

The config package provides extensive validation capabilities:

```go
import "github.com/kerimovok/go-pkg-utils/config"

// Network and URL validation
config.IsValidIP("192.168.1.1")           // true
config.IsValidIPv4("192.168.1.1")         // true
config.IsValidIPv6("::1")                 // true
config.IsValidPort("8080")                // true
config.IsValidHost("api.example.com")     // true
config.IsValidDomain("example.com")       // true
config.IsValidURL("https://example.com")  // true
config.IsValidEmail("user@example.com")   // true

// Number validation
config.IsValidNumber("123.45")            // true
config.IsValidInteger("123")              // true
config.IsValidPositiveInteger("456")      // true
config.IsValidNonNegativeInteger("0")     // true
config.IsValidFloat("123.45")             // true
config.IsValidPositiveFloat("123.45")     // true
config.IsValidNonNegativeFloat("0.0")     // true

// String validation
config.IsValidNonEmptyString("hello")     // true
config.IsValidAlphabetic("HelloWorld")   // true
config.IsValidAlphanumeric("Hello123")   // true
config.IsValidNumeric("12345")           // true
config.IsValidLowercase("hello")         // true
config.IsValidUppercase("HELLO")         // true

// Format validation
config.IsValidUUID("550e8400-e29b-41d4-a716-446655440000") // true
config.IsValidBase64("SGVsbG8gV29ybGQ=")                   // true
config.IsValidHex("1a2b3c4d")                              // true
config.IsValidSlug("hello-world")                          // true
```

### Environment Variables

The config package supports various data types:

```go
// String values
dbHost := config.GetEnvOrDefault("DB_HOST", "localhost")

// Type-safe parsing
dbPort := config.GetEnvInt("DB_PORT", 5432)
enableSSL := config.GetEnvBool("DB_SSL", false)
timeout := config.GetEnvDuration("DB_TIMEOUT", 30*time.Second)
maxSize := config.GetEnvFloat("DB_MAX_SIZE", 100.5)

// Required variables (panics if missing)
secretKey := config.RequiredEnv("SECRET_KEY")
```

### YAML Configuration with Environment Variable Substitution

The config package provides powerful YAML configuration loading with automatic environment variable substitution:

```go
// Define your configuration structure
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    SSL      bool   `yaml:"ssl"`
}

type ServerConfig struct {
    Port    int    `yaml:"port"`
    Host    string `yaml:"host"`
    Timeout string `yaml:"timeout"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Debug    bool           `yaml:"debug"`
}

// Load configuration from YAML file
var cfg AppConfig
err := config.LoadYAMLConfig("config.yaml", &cfg)
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Use your configuration
fmt.Printf("Database: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
```

#### YAML File Example (config.yaml)

```yaml
database:
    host: ${DB_HOST:-localhost} # or ${DB_HOST=localhost}
    port: ${DB_PORT:-5432} # or ${DB_PORT=5432}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
    ssl: ${DB_SSL:-false} # or ${DB_SSL=false}

server:
    port: ${SERVER_PORT:-8080} # or ${SERVER_PORT=8080}
    host: ${SERVER_HOST:-0.0.0.0} # or ${SERVER_HOST=0.0.0.0}
    timeout: ${SERVER_TIMEOUT:-30s} # or ${SERVER_TIMEOUT=30s}

debug: ${DEBUG:-false} # or ${DEBUG=false}
```

#### Environment Variable Substitution Syntax

- `${VARIABLE}` - Required variable (will be empty if not set)
- `${VARIABLE:-default}` or `${VARIABLE=default}` - Optional variable with default value (both syntaxes are equivalent)
- `${VARIABLE:-}` or `${VARIABLE=}` - Optional variable with empty string default

#### Advanced Usage

```go
// Load multiple configuration files
configs := []struct {
    filename string
    target   interface{}
}{
    {"config/database.yaml", &dbConfig},
    {"config/server.yaml", &serverConfig},
    {"config/features.yaml", &featureConfig},
}

for _, cfg := range configs {
    if err := config.LoadYAMLConfig(cfg.filename, cfg.target); err != nil {
        log.Fatalf("Failed to load %s: %v", cfg.filename, err)
    }
}

// Manual environment variable substitution
yamlContent := []byte(`
database:
  host: ${DB_HOST:-localhost}  # or ${DB_HOST=localhost}
  port: ${DB_PORT:-5432}       # or ${DB_PORT=5432}
`)

substituted := config.SubstituteEnvVars(yamlContent)
// Process substituted content...
```

### Validation Rules

```go
rules := []validator.ValidationRule{
    {Variable: "PORT", Default: "8080", Rule: config.IsValidPort, Message: "Invalid port"},
    {Variable: "EMAIL", Rule: config.IsValidEmail, Message: "Invalid email"},
    {Variable: "DB_HOST", Rule: config.IsValidHost, Message: "Invalid database host"},
    {Variable: "API_URL", Rule: config.IsValidURL, Message: "Invalid API URL"},
    {Variable: "MAX_CONNECTIONS", Rule: config.IsValidPositiveInteger, Message: "Max connections must be positive"},
    {Variable: "TIMEOUT", Rule: config.IsValidPositiveFloat, Message: "Timeout must be positive"},
    {Variable: "SECRET_KEY", Rule: config.IsValidNonEmptyString, Message: "Secret key cannot be empty"},
    {Variable: "JWT_SECRET", Rule: config.IsValidBase64, Message: "JWT secret must be valid base64"},
}

err := validator.ValidateConfig(rules)
```

## 🌐 HTTP Response Standards

### Standard Responses

```go
// Success (200)
httpx.OK("Operation successful", data)

// Created (201)
httpx.Created("Resource created", resource)

// Bad Request (400)
httpx.BadRequest("Invalid input", validationError)

// Unauthorized (401)
httpx.Unauthorized("Authentication required")

// Not Found (404)
httpx.NotFound("Resource not found")

// Validation Error (422)
httpx.UnprocessableEntityWithValidation("Validation failed", validationErrors)

// And many more HTTP status codes...
// 2xx: OK, Created, Accepted, NoContent, PartialContent
// 3xx: NotModified
// 4xx: BadRequest, Unauthorized, Forbidden, NotFound, Conflict, TooManyRequests, etc.
// 5xx: InternalServerError, BadGateway, ServiceUnavailable, GatewayTimeout, etc.
```

### Response Structure

All responses follow a consistent structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...},
  "status": 200,
  "timestamp": "2023-12-25T10:30:00Z"
}
```

## 🔐 Security Features

### Password Security

- **bcrypt**: Standard bcrypt hashing for passwords
- **scrypt**: Enhanced security with configurable parameters
- **Secure Random**: Cryptographically secure random generation

### Encryption

- **AES-GCM**: Authenticated encryption for data
- **RSA**: Asymmetric encryption and digital signatures
- **HMAC**: Message authentication codes

### Key Management

- **Key Generation**: Secure key generation for various algorithms
- **PEM Encoding**: Standard key serialization format

## 📊 Collections Features

### Functional Programming

The collections package provides functional programming utilities:

```go
// Functional composition
result := collections.Map(
    collections.Filter(numbers, isEven),
    double,
)

// Aggregation
total := collections.Reduce(numbers, 0, add)
grouped := collections.GroupBy(users, func(u User) string { return u.Department })
```

### Type Safety

All collection functions are generic and type-safe:

```go
// Type-safe operations
strings := []string{"hello", "world"}
lengths := collections.Map(strings, func(s string) int { return len(s) })
```

## 🌟 Best Practices

### Error Handling

```go
// Use structured errors for better debugging
err := errors.ValidationError("INVALID_INPUT", "User input validation failed").
    WithMetadata("field", "email").
    WithOperation("CreateUser").
    WithComponent("user-service")

// Check error types
if errors.IsType(err, errors.ErrorTypeValidation) {
    // Handle validation error
}

// Use error chains for multiple errors
chain := errors.NewErrorChain()
// Add multiple errors...
if chain.HasErrors() {
    return chain // Implements error interface
}
```

### Configuration Management

```go
// Use YAML configuration with environment variable substitution for complex setups
type Config struct {
    Database struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"database"`
}

var cfg Config
err := config.LoadYAMLConfig("config.yaml", &cfg)

// Use required variables for critical settings
secretKey := config.RequiredEnv("SECRET_KEY")

// Provide sensible defaults
port := config.GetEnvInt("PORT", 8080)
debug := config.GetEnvBool("DEBUG", false)

// Validate configuration with comprehensive checks
if !config.IsValidURL(apiURL) {
    log.Fatal("Invalid API URL")
}

if !config.IsValidPort(portStr) {
    log.Fatal("Invalid port number")
}

if !config.IsValidIP(dbHost) {
    log.Fatal("Invalid database host IP")
}

if !config.IsValidEmail(adminEmail) {
    log.Fatal("Invalid admin email")
}

if !config.IsValidPositiveInteger(maxConnections) {
    log.Fatal("Max connections must be positive")
}

if !config.IsValidNonEmptyString(secretKey) {
    log.Fatal("Secret key cannot be empty")
}

// Validate environment variables before use
if !config.IsValidUUID(jwtSecret) {
    log.Fatal("Invalid JWT secret format")
}

if !config.IsValidBase64(encryptionKey) {
    log.Fatal("Encryption key must be valid base64")
}
```

### HTTP Responses

```go
// Use consistent response format
response := httpx.OK("User created", user)
return httpx.SendResponse(c, response)

// Include pagination for lists
pagination := httpx.NewPagination(page, perPage, total)
response := httpx.Paginated("Users retrieved", users, pagination)
return httpx.SendPaginatedResponse(c, response)
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ for the Go community
- Inspired by modern utility libraries and best practices
- Uses battle-tested cryptographic libraries and algorithms

---

**Note**: This package requires Go 1.22 or later for generic support and modern language features.
