# go-pkg-utils

A comprehensive Go utilities package providing essential functionality for modern Go applications including configuration management, HTTP responses, string manipulation, cryptography, collections, validation, error handling, and more.

## 🚀 Features

-   **🔧 Configuration**: Environment variable parsing with type safety and validation
-   **🌐 HTTP Responses**: Standardized API response structures with comprehensive HTTP status code support and Fiber integration
-   **📝 Text Processing**: String manipulation, case conversion, validation, and more
-   **🔐 Cryptography**: Secure random generation, password hashing, AES/RSA encryption, JWT
-   **📊 Collections**: Generic utilities for slices and maps with functional programming support
-   **📅 Date/Time**: Comprehensive time manipulation and formatting utilities
-   **🔄 JSON**: Advanced JSON marshaling, unmarshaling, and manipulation
-   **✅ Validation**: Struct validation with tags and detailed error reporting
-   **🚨 Error Handling**: Structured error system with metadata and stack traces
-   **💬 Messages**: Predefined success and error messages
-   **🌐 Network**: IP address extraction and validation utilities
-   **🆔 UUID**: UUID utilities and validation

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
├── jsonx/          # Advanced JSON processing utilities
├── messages/       # Predefined message constants
├── net/           # Network utilities (IP extraction)
├── text/          # String manipulation and processing
├── uuid/          # UUID utilities
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

-   **bcrypt**: Standard bcrypt hashing for passwords
-   **scrypt**: Enhanced security with configurable parameters
-   **Secure Random**: Cryptographically secure random generation

### Encryption

-   **AES-GCM**: Authenticated encryption for data
-   **RSA**: Asymmetric encryption and digital signatures
-   **HMAC**: Message authentication codes

### Key Management

-   **Key Generation**: Secure key generation for various algorithms
-   **PEM Encoding**: Standard key serialization format

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
pagination := httpx.NewPagination(page, limit, total)
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

-   Built with ❤️ for the Go community
-   Inspired by modern utility libraries and best practices
-   Uses battle-tested cryptographic libraries and algorithms

---

**Note**: This package requires Go 1.22 or later for generic support and modern language features.
