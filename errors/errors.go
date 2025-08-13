package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ErrorTypeValidation         ErrorType = "validation"
	ErrorTypeNotFound           ErrorType = "not_found"
	ErrorTypeUnauthorized       ErrorType = "unauthorized"
	ErrorTypeForbidden          ErrorType = "forbidden"
	ErrorTypeConflict           ErrorType = "conflict"
	ErrorTypeInternal           ErrorType = "internal"
	ErrorTypeExternal           ErrorType = "external"
	ErrorTypeTimeout            ErrorType = "timeout"
	ErrorTypeRateLimit          ErrorType = "rate_limit"
	ErrorTypeBadRequest         ErrorType = "bad_request"
	ErrorTypeServiceUnavailable ErrorType = "service_unavailable"
)

// Error represents a structured error with additional context
type Error struct {
	Type       ErrorType              `json:"type"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Cause      error                  `json:"-"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	StackTrace []StackFrame           `json:"stack_trace,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	Operation  string                 `json:"operation,omitempty"`
	Component  string                 `json:"component,omitempty"`
	Retryable  bool                   `json:"retryable"`
	HTTPStatus int                    `json:"http_status,omitempty"`
}

// StackFrame represents a stack frame
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s:%s] %s: %s", e.Type, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Type == t.Type && e.Code == t.Code
	}
	return false
}

// JSON returns the error as JSON
func (e *Error) JSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// WithCause adds a cause to the error
func (e *Error) WithCause(cause error) *Error {
	e.Cause = cause
	return e
}

// WithDetails adds details to the error
func (e *Error) WithDetails(details string) *Error {
	e.Details = details
	return e
}

// WithMetadata adds metadata to the error
func (e *Error) WithMetadata(key string, value interface{}) *Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithRequestID adds a request ID to the error
func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

// WithUserID adds a user ID to the error
func (e *Error) WithUserID(userID string) *Error {
	e.UserID = userID
	return e
}

// WithOperation adds an operation name to the error
func (e *Error) WithOperation(operation string) *Error {
	e.Operation = operation
	return e
}

// WithComponent adds a component name to the error
func (e *Error) WithComponent(component string) *Error {
	e.Component = component
	return e
}

// WithHTTPStatus adds an HTTP status code to the error
func (e *Error) WithHTTPStatus(status int) *Error {
	e.HTTPStatus = status
	return e
}

// MarkRetryable marks the error as retryable
func (e *Error) MarkRetryable() *Error {
	e.Retryable = true
	return e
}

// NewError creates a new structured error
func NewError(errorType ErrorType, code, message string) *Error {
	return &Error{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Timestamp:  time.Now(),
		StackTrace: captureStackTrace(),
		Metadata:   make(map[string]interface{}),
	}
}

// captureStackTrace captures the current stack trace
func captureStackTrace() []StackFrame {
	var frames []StackFrame

	// Skip the first 3 frames (runtime.Callers, captureStackTrace, NewError)
	pcs := make([]uintptr, 10)
	n := runtime.Callers(3, pcs)

	for i := 0; i < n; i++ {
		pc := pcs[i]
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		file, line := fn.FileLine(pc)

		frames = append(frames, StackFrame{
			Function: fn.Name(),
			File:     file,
			Line:     line,
		})
	}

	return frames
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, code, message string) *Error {
	if err == nil {
		return nil
	}

	// If it's already our Error type, add context
	if e, ok := err.(*Error); ok {
		return &Error{
			Type:       errorType,
			Code:       code,
			Message:    message,
			Cause:      e,
			Timestamp:  time.Now(),
			StackTrace: captureStackTrace(),
			Metadata:   make(map[string]interface{}),
		}
	}

	return &Error{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Cause:      err,
		Timestamp:  time.Now(),
		StackTrace: captureStackTrace(),
		Metadata:   make(map[string]interface{}),
	}
}

// Common error constructors

// ValidationError creates a validation error
func ValidationError(code, message string) *Error {
	return NewError(ErrorTypeValidation, code, message).WithHTTPStatus(400)
}

// NotFoundError creates a not found error
func NotFoundError(code, message string) *Error {
	return NewError(ErrorTypeNotFound, code, message).WithHTTPStatus(404)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(code, message string) *Error {
	return NewError(ErrorTypeUnauthorized, code, message).WithHTTPStatus(401)
}

// ForbiddenError creates a forbidden error
func ForbiddenError(code, message string) *Error {
	return NewError(ErrorTypeForbidden, code, message).WithHTTPStatus(403)
}

// ConflictError creates a conflict error
func ConflictError(code, message string) *Error {
	return NewError(ErrorTypeConflict, code, message).WithHTTPStatus(409)
}

// InternalError creates an internal error
func InternalError(code, message string) *Error {
	return NewError(ErrorTypeInternal, code, message).WithHTTPStatus(500)
}

// ExternalError creates an external service error
func ExternalError(code, message string) *Error {
	return NewError(ErrorTypeExternal, code, message).WithHTTPStatus(502).MarkRetryable()
}

// TimeoutError creates a timeout error
func TimeoutError(code, message string) *Error {
	return NewError(ErrorTypeTimeout, code, message).WithHTTPStatus(408).MarkRetryable()
}

// RateLimitError creates a rate limit error
func RateLimitError(code, message string) *Error {
	return NewError(ErrorTypeRateLimit, code, message).WithHTTPStatus(429).MarkRetryable()
}

// BadRequestError creates a bad request error
func BadRequestError(code, message string) *Error {
	return NewError(ErrorTypeBadRequest, code, message).WithHTTPStatus(400)
}

// ServiceUnavailableError creates a service unavailable error
func ServiceUnavailableError(code, message string) *Error {
	return NewError(ErrorTypeServiceUnavailable, code, message).WithHTTPStatus(503).MarkRetryable()
}

// ErrorChain represents a chain of errors for aggregation
type ErrorChain struct {
	Errors []*Error `json:"errors"`
}

// Add adds an error to the chain
func (ec *ErrorChain) Add(err *Error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
	}
}

// AddError adds a regular error to the chain by wrapping it
func (ec *ErrorChain) AddError(err error, errorType ErrorType, code, message string) {
	if err != nil {
		ec.Add(Wrap(err, errorType, code, message))
	}
}

// HasErrors returns true if there are any errors in the chain
func (ec *ErrorChain) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Error implements the error interface
func (ec *ErrorChain) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}

	if len(ec.Errors) == 1 {
		return ec.Errors[0].Error()
	}

	var messages []string
	for _, err := range ec.Errors {
		messages = append(messages, err.Error())
	}

	return fmt.Sprintf("multiple errors: %s", strings.Join(messages, "; "))
}

// First returns the first error in the chain
func (ec *ErrorChain) First() *Error {
	if len(ec.Errors) > 0 {
		return ec.Errors[0]
	}
	return nil
}

// Last returns the last error in the chain
func (ec *ErrorChain) Last() *Error {
	if len(ec.Errors) > 0 {
		return ec.Errors[len(ec.Errors)-1]
	}
	return nil
}

// Count returns the number of errors in the chain
func (ec *ErrorChain) Count() int {
	return len(ec.Errors)
}

// Filter returns errors of a specific type
func (ec *ErrorChain) Filter(errorType ErrorType) []*Error {
	var filtered []*Error
	for _, err := range ec.Errors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// JSON returns the error chain as JSON
func (ec *ErrorChain) JSON() string {
	data, _ := json.Marshal(ec)
	return string(data)
}

// NewErrorChain creates a new error chain
func NewErrorChain() *ErrorChain {
	return &ErrorChain{
		Errors: make([]*Error, 0),
	}
}

// ErrorHandler provides utilities for error handling
type ErrorHandler struct {
	DefaultComponent string
	Logger           func(error)
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(component string, logger func(error)) *ErrorHandler {
	return &ErrorHandler{
		DefaultComponent: component,
		Logger:           logger,
	}
}

// Handle handles an error by logging it and optionally returning a sanitized version
func (eh *ErrorHandler) Handle(err error, sanitize bool) error {
	if err == nil {
		return nil
	}

	// Log the error
	if eh.Logger != nil {
		eh.Logger(err)
	}

	// If sanitization is requested and it's an internal error, return a generic error
	if sanitize {
		if e, ok := err.(*Error); ok {
			if e.Type == ErrorTypeInternal {
				return InternalError("INTERNAL_ERROR", "An internal error occurred")
			}
		}
	}

	return err
}

// Recover recovers from panics and converts them to errors
func (eh *ErrorHandler) Recover() error {
	if r := recover(); r != nil {
		err := InternalError("PANIC", fmt.Sprintf("Panic recovered: %v", r))
		if eh.DefaultComponent != "" {
			err.WithComponent(eh.DefaultComponent)
		}

		if eh.Logger != nil {
			eh.Logger(err)
		}

		return err
	}
	return nil
}

// SafeExecute executes a function and handles any panics
func (eh *ErrorHandler) SafeExecute(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = InternalError("PANIC", fmt.Sprintf("Panic recovered: %v", r))
			if eh.DefaultComponent != "" {
				err.(*Error).WithComponent(eh.DefaultComponent)
			}

			if eh.Logger != nil {
				eh.Logger(err)
			}
		}
	}()

	return fn()
}

// IsType checks if an error is of a specific type
func IsType(err error, errorType ErrorType) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == errorType
	}
	return false
}

// IsCode checks if an error has a specific code
func IsCode(err error, code string) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Retryable
	}
	return false
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if e, ok := err.(*Error); ok && e.HTTPStatus > 0 {
		return e.HTTPStatus
	}
	return 500 // Default to internal server error
}

// GetErrorType returns the error type
func GetErrorType(err error) ErrorType {
	if e, ok := err.(*Error); ok {
		return e.Type
	}
	return ErrorTypeInternal
}

// GetErrorCode returns the error code
func GetErrorCode(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return "UNKNOWN"
}

// CombineErrors combines multiple errors into an error chain
func CombineErrors(errors ...error) error {
	chain := NewErrorChain()

	for _, err := range errors {
		if err != nil {
			if e, ok := err.(*Error); ok {
				chain.Add(e)
			} else {
				chain.Add(Wrap(err, ErrorTypeInternal, "WRAPPED_ERROR", err.Error()))
			}
		}
	}

	if !chain.HasErrors() {
		return nil
	}

	if chain.Count() == 1 {
		return chain.First()
	}

	return chain
}
