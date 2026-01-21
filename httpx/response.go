package httpx

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Response
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page         int   `json:"page"`
	PerPage      int   `json:"perPage"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"totalPages"`
	HasNext      bool  `json:"hasNext"`
	HasPrevious  bool  `json:"hasPrevious"`
	NextPage     *int  `json:"nextPage,omitempty"`
	PreviousPage *int  `json:"previousPage,omitempty"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResponse represents a validation error response
type ValidationResponse struct {
	Response
	Errors []ValidationError `json:"validation_errors,omitempty"`
}

// Paginated creates a paginated success response
func Paginated(message string, data interface{}, pagination *Pagination) PaginatedResponse {
	return PaginatedResponse{
		Response: Response{
			Success:   true,
			Message:   message,
			Data:      data,
			Status:    fiber.StatusOK,
			Timestamp: time.Now(),
		},
		Pagination: pagination,
	}
}

// NewPagination creates pagination metadata
func NewPagination(page, perPage int, total int64) *Pagination {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &Pagination{
		Page:        page,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	if pagination.HasNext {
		nextPage := page + 1
		pagination.NextPage = &nextPage
	}

	if pagination.HasPrevious {
		prevPage := page - 1
		pagination.PreviousPage = &prevPage
	}

	return pagination
}

// SendResponse sends a response using Fiber context
func SendResponse(c *fiber.Ctx, response Response) error {
	return c.Status(response.Status).JSON(response)
}

// SendPaginatedResponse sends a paginated response using Fiber context
func SendPaginatedResponse(c *fiber.Ctx, response PaginatedResponse) error {
	return c.Status(response.Status).JSON(response)
}

// SendValidationResponse sends a validation error response using Fiber context
func SendValidationResponse(c *fiber.Ctx, response ValidationResponse) error {
	return c.Status(response.Status).JSON(response)
}
