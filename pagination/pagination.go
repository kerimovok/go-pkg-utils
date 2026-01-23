package pagination

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kerimovok/go-pkg-utils/httpx"
	"github.com/kerimovok/go-pkg-utils/validator"
	"gorm.io/gorm"
)

// Params represents pagination query parameters
type Params struct {
	Page      int    `query:"page" validate:"min=1"`
	PerPage   int    `query:"per_page" validate:"min=1,max=100"`             // Items per page
	SortBy    string `query:"sort_by"`                                       // Sort field name
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"` // Sort order: asc or desc
}

// Defaults holds default values for pagination
type Defaults struct {
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

// Default returns sensible defaults for pagination
func Default() Defaults {
	return Defaults{
		Page:      1,
		PerPage:   20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// ParseParams parses and validates pagination parameters from Fiber context
func ParseParams(c *fiber.Ctx, defaults Defaults) (*Params, error) {
	var params Params
	if err := c.QueryParser(&params); err != nil {
		return nil, err
	}

	// Apply defaults first (before validation) for missing/zero values
	if params.Page <= 0 {
		params.Page = defaults.Page
	}
	if params.PerPage <= 0 {
		params.PerPage = defaults.PerPage
	}
	if params.SortBy == "" {
		params.SortBy = defaults.SortBy
	}
	if params.SortOrder == "" {
		params.SortOrder = defaults.SortOrder
	}

	// Validate after defaults are applied
	if err := validator.ValidateStruct(&params); err != nil {
		return nil, err
	}

	return &params, nil
}

// Query applies pagination to a GORM query and returns results with metadata
func Query[T any](
	ctx context.Context,
	query *gorm.DB,
	params *Params,
	message string,
) (*httpx.PaginatedResponse, error) {
	// Clone query for counting (before applying per_page/offset)
	// Use Session to create a new query instance with the same conditions
	countQuery := query.Session(&gorm.Session{}).WithContext(ctx)

	// Get total count
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting and pagination
	offset := (params.Page - 1) * params.PerPage
	query = query.WithContext(ctx).Order(params.SortBy + " " + params.SortOrder).
		Offset(offset).
		Limit(params.PerPage)

	// Execute query
	var results []T
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// Build paginated response
	pagination := httpx.NewPagination(params.Page, params.PerPage, total)
	response := httpx.Paginated(message, results, pagination)

	return &response, nil
}

// HandleRequest is a convenience function that handles the full pagination flow
// It parses, validates, applies defaults, executes query, and returns response
func HandleRequest[T any](
	c *fiber.Ctx,
	query *gorm.DB,
	defaults Defaults,
	message string,
) error {
	// Parse and validate pagination params
	params, err := ParseParams(c, defaults)
	if err != nil {
		response := httpx.BadRequest("Invalid query parameters", err)
		return httpx.SendResponse(c, response)
	}

	// Create context with timeout from request context
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Execute paginated query
	response, err := Query[T](ctx, query, params, message)
	if err != nil {
		// Error logging should be handled by the caller or middleware
		response := httpx.InternalServerError("Failed to retrieve data", err)
		return httpx.SendResponse(c, response)
	}

	return httpx.SendPaginatedResponse(c, *response)
}
