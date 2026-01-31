package filter

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Operator represents a filter operator
type Operator string

const (
	// OperatorEQ equals
	OperatorEQ Operator = "eq"
	// OperatorNE not equals
	OperatorNE Operator = "ne"
	// OperatorGT greater than
	OperatorGT Operator = "gt"
	// OperatorGTE greater than or equal
	OperatorGTE Operator = "gte"
	// OperatorLT less than
	OperatorLT Operator = "lt"
	// OperatorLTE less than or equal
	OperatorLTE Operator = "lte"
	// OperatorLIKE like (for strings)
	OperatorLIKE Operator = "like"
	// OperatorIN in (for arrays)
	OperatorIN Operator = "in"
	// OperatorNOTIN not in (for arrays)
	OperatorNOTIN Operator = "not_in"
)

// AllowedOperators is a map of all allowed operators
var AllowedOperators = map[Operator]bool{
	OperatorEQ:    true,
	OperatorNE:    true,
	OperatorGT:    true,
	OperatorGTE:   true,
	OperatorLT:    true,
	OperatorLTE:   true,
	OperatorLIKE:  true,
	OperatorIN:    true,
	OperatorNOTIN: true,
}

// Filter represents a single filter condition
type Filter struct {
	Field    string
	Operator Operator
	Value    interface{}
}

// Config holds filter configuration
type Config struct {
	// AllowedFields is a map of allowed field names to their types
	// If empty, all fields are allowed
	AllowedFields map[string]string // "field_name" -> "string", "int", "time", etc.
	// FieldMapping maps query parameter names to database column names
	// If empty, field names are used as-is
	FieldMapping map[string]string
	// CustomValidators allows custom validation for specific fields
	CustomValidators map[string]func(value string) error
}

// ParseFilters parses filters from Fiber query parameters
// Format: field_operator=value (e.g., created_at_gte=2024-01-01, status_eq=active)
func ParseFilters(c *fiber.Ctx, config *Config) ([]Filter, error) {
	var filters []Filter

	// Get all query parameters
	queryParams := c.Queries()

	for key, value := range queryParams {
		// Skip pagination and sorting params
		if isReservedParam(key) {
			continue
		}

		// Parse field and operator from key (format: field_operator)
		field, operator, err := parseFilterKey(key)
		if err != nil {
			// If it doesn't match the pattern, skip it (might be a custom param)
			continue
		}

		// Validate operator
		if !AllowedOperators[operator] {
			return nil, fmt.Errorf("invalid operator '%s' for field '%s'", operator, field)
		}

		// Map field name if mapping is provided
		dbField := field
		if config != nil && config.FieldMapping != nil {
			if mapped, ok := config.FieldMapping[field]; ok {
				dbField = mapped
			}
		}

		// Check if field is allowed
		if config != nil && config.AllowedFields != nil {
			if _, allowed := config.AllowedFields[field]; !allowed {
				return nil, fmt.Errorf("field '%s' is not allowed for filtering", field)
			}
		}

		// Custom validation if provided
		if config != nil && config.CustomValidators != nil {
			if validator, ok := config.CustomValidators[field]; ok {
				if err := validator(value); err != nil {
					return nil, fmt.Errorf("validation failed for field '%s': %w", field, err)
				}
			}
		}

		// Convert value based on field type and operator
		convertedValue, err := convertValue(value, field, operator, config)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value for field '%s': %w", field, err)
		}

		filters = append(filters, Filter{
			Field:    dbField,
			Operator: operator,
			Value:    convertedValue,
		})
	}

	return filters, nil
}

// ApplyFilters applies filters to a GORM query
func ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		query = applyFilter(query, f)
	}
	return query
}

// ApplyFiltersFromContext is a convenience function that parses and applies filters
func ApplyFiltersFromContext(c *fiber.Ctx, query *gorm.DB, config *Config) (*gorm.DB, error) {
	filters, err := ParseFilters(c, config)
	if err != nil {
		return nil, err
	}
	return ApplyFilters(query, filters), nil
}

// parseFilterKey parses a query key in the format "field_operator"
func parseFilterKey(key string) (field string, operator Operator, err error) {
	parts := strings.Split(key, "_")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid filter format")
	}

	// Get operator (last part)
	opStr := strings.ToLower(parts[len(parts)-1])
	operator = Operator(opStr)

	// Get field (everything before the last part)
	field = strings.Join(parts[:len(parts)-1], "_")

	return field, operator, nil
}

// convertValue converts a string value to the appropriate type based on field configuration
func convertValue(value, field string, operator Operator, config *Config) (interface{}, error) {
	// Handle IN and NOT_IN operators - they expect arrays
	if operator == OperatorIN || operator == OperatorNOTIN {
		values := strings.Split(value, ",")
		var result []interface{}

		// Determine type from config
		fieldType := "string"
		if config != nil && config.AllowedFields != nil {
			if ft, ok := config.AllowedFields[field]; ok {
				fieldType = ft
			}
		}

		for _, v := range values {
			v = strings.TrimSpace(v)
			converted, err := convertSingleValue(v, fieldType)
			if err != nil {
				return nil, err
			}
			result = append(result, converted)
		}
		return result, nil
	}

	// Single value conversion
	fieldType := "string"
	if config != nil && config.AllowedFields != nil {
		if ft, ok := config.AllowedFields[field]; ok {
			fieldType = ft
		}
	}

	return convertSingleValue(value, fieldType)
}

// convertSingleValue converts a single string value to the specified type
func convertSingleValue(value, fieldType string) (interface{}, error) {
	switch strings.ToLower(fieldType) {
	case "int", "integer":
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
			return nil, fmt.Errorf("invalid integer value: %s", value)
		}
		return result, nil
	case "float", "float64":
		var result float64
		if _, err := fmt.Sscanf(value, "%f", &result); err != nil {
			return nil, fmt.Errorf("invalid float value: %s", value)
		}
		return result, nil
	case "bool", "boolean":
		return strings.ToLower(value) == "true" || value == "1", nil
	case "time", "datetime", "date":
		// Try ISO 8601 format first
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			return t, nil
		}
		// Try date only
		if t, err := time.Parse("2006-01-02", value); err == nil {
			return t, nil
		}
		// Try date with time
		if t, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
			return t, nil
		}
		return nil, fmt.Errorf("invalid time format: %s (expected RFC3339, 2006-01-02, or 2006-01-02T15:04:05)", value)
	default:
		// Default to string
		return value, nil
	}
}

// applyFilter applies a single filter to a GORM query
func applyFilter(query *gorm.DB, f Filter) *gorm.DB {
	switch f.Operator {
	case OperatorEQ:
		return query.Where(f.Field+" = ?", f.Value)
	case OperatorNE:
		return query.Where(f.Field+" != ?", f.Value)
	case OperatorGT:
		return query.Where(f.Field+" > ?", f.Value)
	case OperatorGTE:
		return query.Where(f.Field+" >= ?", f.Value)
	case OperatorLT:
		return query.Where(f.Field+" < ?", f.Value)
	case OperatorLTE:
		return query.Where(f.Field+" <= ?", f.Value)
	case OperatorLIKE:
		// Add % wildcards for LIKE queries
		likeValue := fmt.Sprintf("%%%s%%", f.Value)
		return query.Where(f.Field+" LIKE ?", likeValue)
	case OperatorIN:
		// Value should be a comma-separated string or slice
		return query.Where(f.Field+" IN ?", f.Value)
	case OperatorNOTIN:
		// Value should be a comma-separated string or slice
		return query.Where(f.Field+" NOT IN ?", f.Value)
	default:
		return query
	}
}

// isReservedParam checks if a parameter is reserved for pagination/sorting
// Uses case-insensitive matching to support both snake_case and camelCase
func isReservedParam(key string) bool {
	reserved := []string{"page", "per_page", "sort_by", "sort_order"}
	keyLower := strings.ToLower(key)
	for _, r := range reserved {
		if keyLower == strings.ToLower(r) {
			return true
		}
	}
	return false
}
