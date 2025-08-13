package httpx

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// 2xx Success Status Codes

// OK creates a 200 OK response
func OK(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Status:    fiber.StatusOK,
		Timestamp: time.Now(),
	}
}

// Created creates a 201 Created response
func Created(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Status:    fiber.StatusCreated,
		Timestamp: time.Now(),
	}
}

// Accepted creates a 202 Accepted response
func Accepted(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Status:    fiber.StatusAccepted,
		Timestamp: time.Now(),
	}
}

// NoContent creates a 204 No Content response
func NoContent(message string) Response {
	return Response{
		Success:   true,
		Message:   message,
		Status:    fiber.StatusNoContent,
		Timestamp: time.Now(),
	}
}

// PartialContent creates a 206 Partial Content response
func PartialContent(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Status:    fiber.StatusPartialContent,
		Timestamp: time.Now(),
	}
}

// 3xx Redirection Status Codes

// NotModified creates a 304 Not Modified response
func NotModified(message string) Response {
	return Response{
		Success:   true,
		Message:   message,
		Status:    fiber.StatusNotModified,
		Timestamp: time.Now(),
	}
}

// 4xx Client Error Status Codes

// BadRequest creates a 400 Bad Request response
func BadRequest(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success:   false,
		Message:   message,
		Error:     errMsg,
		Status:    fiber.StatusBadRequest,
		Timestamp: time.Now(),
	}
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusUnauthorized,
		Timestamp: time.Now(),
	}
}

// PaymentRequired creates a 402 Payment Required response
func PaymentRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusPaymentRequired,
		Timestamp: time.Now(),
	}
}

// Forbidden creates a 403 Forbidden response
func Forbidden(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusForbidden,
		Timestamp: time.Now(),
	}
}

// NotFound creates a 404 Not Found response
func NotFound(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusNotFound,
		Timestamp: time.Now(),
	}
}

// MethodNotAllowed creates a 405 Method Not Allowed response
func MethodNotAllowed(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusMethodNotAllowed,
		Timestamp: time.Now(),
	}
}

// NotAcceptable creates a 406 Not Acceptable response
func NotAcceptable(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusNotAcceptable,
		Timestamp: time.Now(),
	}
}

// ProxyAuthenticationRequired creates a 407 Proxy Authentication Required response
func ProxyAuthenticationRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusProxyAuthRequired,
		Timestamp: time.Now(),
	}
}

// RequestTimeout creates a 408 Request Timeout response
func RequestTimeout(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusRequestTimeout,
		Timestamp: time.Now(),
	}
}

// Conflict creates a 409 Conflict response
func Conflict(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success:   false,
		Message:   message,
		Error:     errMsg,
		Status:    fiber.StatusConflict,
		Timestamp: time.Now(),
	}
}

// Gone creates a 410 Gone response
func Gone(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusGone,
		Timestamp: time.Now(),
	}
}

// LengthRequired creates a 411 Length Required response
func LengthRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusLengthRequired,
		Timestamp: time.Now(),
	}
}

// PreconditionFailed creates a 412 Precondition Failed response
func PreconditionFailed(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusPreconditionFailed,
		Timestamp: time.Now(),
	}
}

// PayloadTooLarge creates a 413 Payload Too Large response
func PayloadTooLarge(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusRequestEntityTooLarge,
		Timestamp: time.Now(),
	}
}

// URITooLong creates a 414 URI Too Long response
func URITooLong(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusRequestURITooLong,
		Timestamp: time.Now(),
	}
}

// UnsupportedMediaType creates a 415 Unsupported Media Type response
func UnsupportedMediaType(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusUnsupportedMediaType,
		Timestamp: time.Now(),
	}
}

// RangeNotSatisfiable creates a 416 Range Not Satisfiable response
func RangeNotSatisfiable(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusRequestedRangeNotSatisfiable,
		Timestamp: time.Now(),
	}
}

// ExpectationFailed creates a 417 Expectation Failed response
func ExpectationFailed(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusExpectationFailed,
		Timestamp: time.Now(),
	}
}

// Teapot creates a 418 I'm a teapot response (RFC 2324)
func Teapot(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusTeapot,
		Timestamp: time.Now(),
	}
}

// UnprocessableEntity creates a 422 Unprocessable Entity response
func UnprocessableEntity(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success:   false,
		Message:   message,
		Error:     errMsg,
		Status:    fiber.StatusUnprocessableEntity,
		Timestamp: time.Now(),
	}
}

// Locked creates a 423 Locked response
func Locked(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusLocked,
		Timestamp: time.Now(),
	}
}

// FailedDependency creates a 424 Failed Dependency response
func FailedDependency(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusFailedDependency,
		Timestamp: time.Now(),
	}
}

// TooEarly creates a 425 Too Early response
func TooEarly(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusTooEarly,
		Timestamp: time.Now(),
	}
}

// UpgradeRequired creates a 426 Upgrade Required response
func UpgradeRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusUpgradeRequired,
		Timestamp: time.Now(),
	}
}

// PreconditionRequired creates a 428 Precondition Required response
func PreconditionRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusPreconditionRequired,
		Timestamp: time.Now(),
	}
}

// TooManyRequests creates a 429 Too Many Requests response
func TooManyRequests(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusTooManyRequests,
		Timestamp: time.Now(),
	}
}

// RequestHeaderFieldsTooLarge creates a 431 Request Header Fields Too Large response
func RequestHeaderFieldsTooLarge(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusRequestHeaderFieldsTooLarge,
		Timestamp: time.Now(),
	}
}

// UnavailableForLegalReasons creates a 451 Unavailable For Legal Reasons response
func UnavailableForLegalReasons(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusUnavailableForLegalReasons,
		Timestamp: time.Now(),
	}
}

// 5xx Server Error Status Codes

// InternalServerError creates a 500 Internal Server Error response
func InternalServerError(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success:   false,
		Message:   message,
		Error:     errMsg,
		Status:    fiber.StatusInternalServerError,
		Timestamp: time.Now(),
	}
}

// CustomStatus creates a response with custom status code
func CustomStatus(message string, err error, status int) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success:   false,
		Message:   message,
		Error:     errMsg,
		Status:    status,
		Timestamp: time.Now(),
	}
}

// NotImplemented creates a 501 Not Implemented response
func NotImplemented(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusNotImplemented,
		Timestamp: time.Now(),
	}
}

// BadGateway creates a 502 Bad Gateway response
func BadGateway(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusBadGateway,
		Timestamp: time.Now(),
	}
}

// ServiceUnavailable creates a 503 Service Unavailable response
func ServiceUnavailable(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusServiceUnavailable,
		Timestamp: time.Now(),
	}
}

// GatewayTimeout creates a 504 Gateway Timeout response
func GatewayTimeout(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusGatewayTimeout,
		Timestamp: time.Now(),
	}
}

// HTTPVersionNotSupported creates a 505 HTTP Version Not Supported response
func HTTPVersionNotSupported(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusHTTPVersionNotSupported,
		Timestamp: time.Now(),
	}
}

// VariantAlsoNegotiates creates a 506 Variant Also Negotiates response
func VariantAlsoNegotiates(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusVariantAlsoNegotiates,
		Timestamp: time.Now(),
	}
}

// InsufficientStorage creates a 507 Insufficient Storage response
func InsufficientStorage(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusInsufficientStorage,
		Timestamp: time.Now(),
	}
}

// LoopDetected creates a 508 Loop Detected response
func LoopDetected(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusLoopDetected,
		Timestamp: time.Now(),
	}
}

// NotExtended creates a 510 Not Extended response
func NotExtended(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusNotExtended,
		Timestamp: time.Now(),
	}
}

// NetworkAuthenticationRequired creates a 511 Network Authentication Required response
func NetworkAuthenticationRequired(message string) Response {
	return Response{
		Success:   false,
		Message:   message,
		Status:    fiber.StatusNetworkAuthenticationRequired,
		Timestamp: time.Now(),
	}
}

// UnprocessableEntityWithValidation creates a 422 Unprocessable Entity response with validation errors
func UnprocessableEntityWithValidation(message string, errors []ValidationError) ValidationResponse {
	return ValidationResponse{
		Response: Response{
			Success:   false,
			Message:   message,
			Status:    fiber.StatusUnprocessableEntity,
			Timestamp: time.Now(),
		},
		Errors: errors,
	}
}
