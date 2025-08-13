package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kerimovok/go-pkg-utils/config"
	"github.com/kerimovok/go-pkg-utils/jsonx"
)

// ValidationRule represents a validation rule for environment variables
type ValidationRule struct {
	Variable string
	Default  string
	Rule     func(value string) bool
	Message  string
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

// Error implements the error interface
func (fe FieldError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", fe.Field, fe.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []FieldError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors checks if there are any validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// ValidateConfig validates environment variables using rules
func ValidateConfig(rules []ValidationRule) error {
	var errors []string
	for _, rule := range rules {
		value := config.GetEnvOrDefault(rule.Variable, rule.Default)
		if !rule.Rule(value) {
			errors = append(errors, rule.Message)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ValidateStruct validates a struct using reflection and tags
func ValidateStruct(s interface{}) ValidationErrors {
	var errors ValidationErrors

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ValidationErrors{FieldError{
				Field:   "struct",
				Message: "struct cannot be nil",
			}}
		}
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ValidationErrors{FieldError{
			Field:   "input",
			Message: "input must be a struct",
		}}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldName := getFieldName(fieldType)
		fieldValue := field.Interface()

		if fieldErrors := validateField(fieldName, fieldValue, tag); len(fieldErrors) > 0 {
			errors = append(errors, fieldErrors...)
		}
	}

	return errors
}

// getFieldName returns the field name for validation (uses json tag if available)
func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		// Extract field name from json tag (before comma)
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}
	return field.Name
}

// validateField validates a single field value against validation tags
func validateField(fieldName string, value interface{}, tag string) ValidationErrors {
	var errors ValidationErrors

	// Parse validation tags
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Parse rule and parameters
		parts := strings.SplitN(rule, "=", 2)
		ruleName := parts[0]
		var param string
		if len(parts) > 1 {
			param = parts[1]
		}

		if err := applyValidationRule(fieldName, value, ruleName, param); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// applyValidationRule applies a specific validation rule
func applyValidationRule(fieldName string, value interface{}, ruleName, param string) *FieldError {
	switch ruleName {
	case "required":
		if isEmpty(value) {
			return &FieldError{
				Field:   fieldName,
				Message: "field is required",
				Tag:     "required",
			}
		}
	case "min":
		if param == "" {
			return &FieldError{
				Field:   fieldName,
				Message: "min rule requires a parameter",
				Tag:     "min",
			}
		}
		return validateMin(fieldName, value, param)
	case "max":
		if param == "" {
			return &FieldError{
				Field:   fieldName,
				Message: "max rule requires a parameter",
				Tag:     "max",
			}
		}
		return validateMax(fieldName, value, param)
	case "email":
		return validateEmail(fieldName, value)
	case "url":
		return validateURL(fieldName, value)
	case "regex":
		if param == "" {
			return &FieldError{
				Field:   fieldName,
				Message: "regex rule requires a pattern parameter",
				Tag:     "regex",
			}
		}
		return validateRegex(fieldName, value, param)
	case "numeric":
		return validateNumeric(fieldName, value)
	case "alpha":
		return validateAlpha(fieldName, value)
	case "alphanum":
		return validateAlphaNumeric(fieldName, value)
	case "uuid":
		return validateUUID(fieldName, value)
	case "json":
		return validateJSON(fieldName, value)
	case "ip":
		return validateIP(fieldName, value)
	case "ipv4":
		return validateIPv4(fieldName, value)
	case "ipv6":
		return validateIPv6(fieldName, value)
	case "date":
		return validateDate(fieldName, value)
	case "datetime":
		return validateDateTime(fieldName, value)
	}

	return nil
}

// isEmpty checks if a value is empty
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if t, ok := value.(time.Time); ok {
			return t.IsZero()
		}
	}

	return false
}

// validateMin validates minimum value/length
func validateMin(fieldName string, value interface{}, param string) *FieldError {
	minVal, err := strconv.Atoi(param)
	if err != nil {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid min parameter",
			Tag:     "min",
		}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if v.Len() < minVal {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at least %d", minVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "min",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < int64(minVal) {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at least %d", minVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "min",
			}
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() < float64(minVal) {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at least %d", minVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "min",
			}
		}
	}

	return nil
}

// validateMax validates maximum value/length
func validateMax(fieldName string, value interface{}, param string) *FieldError {
	maxVal, err := strconv.Atoi(param)
	if err != nil {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid max parameter",
			Tag:     "max",
		}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if v.Len() > maxVal {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at most %d", maxVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "max",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > int64(maxVal) {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at most %d", maxVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "max",
			}
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() > float64(maxVal) {
			return &FieldError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at most %d", maxVal),
				Value:   fmt.Sprintf("%v", value),
				Tag:     "max",
			}
		}
	}

	return nil
}

// validateEmail validates email format
func validateEmail(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "email validation requires string value",
			Tag:     "email",
		}
	}

	if !config.IsValidEmail(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid email format",
			Value:   str,
			Tag:     "email",
		}
	}

	return nil
}

// validateURL validates URL format
func validateURL(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "url validation requires string value",
			Tag:     "url",
		}
	}

	if !config.IsValidURL(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid URL format",
			Value:   str,
			Tag:     "url",
		}
	}

	return nil
}

// validateRegex validates against a regex pattern
func validateRegex(fieldName string, value interface{}, pattern string) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "regex validation requires string value",
			Tag:     "regex",
		}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid regex pattern",
			Tag:     "regex",
		}
	}

	if !regex.MatchString(str) {
		return &FieldError{
			Field:   fieldName,
			Message: fmt.Sprintf("value does not match pattern %s", pattern),
			Value:   str,
			Tag:     "regex",
		}
	}

	return nil
}

// validateNumeric validates numeric values
func validateNumeric(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "numeric validation requires string value",
			Tag:     "numeric",
		}
	}

	regex := regexp.MustCompile(`^[0-9]+$`)
	if !regex.MatchString(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "value must be numeric",
			Value:   str,
			Tag:     "numeric",
		}
	}

	return nil
}

// validateAlpha validates alphabetic values
func validateAlpha(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "alpha validation requires string value",
			Tag:     "alpha",
		}
	}

	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !regex.MatchString(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "value must contain only letters",
			Value:   str,
			Tag:     "alpha",
		}
	}

	return nil
}

// validateAlphaNumeric validates alphanumeric values
func validateAlphaNumeric(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "alphanum validation requires string value",
			Tag:     "alphanum",
		}
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !regex.MatchString(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "value must contain only letters and numbers",
			Value:   str,
			Tag:     "alphanum",
		}
	}

	return nil
}

// validateUUID validates UUID format
func validateUUID(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "uuid validation requires string value",
			Tag:     "uuid",
		}
	}

	if !config.IsValidUUID(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid UUID format",
			Value:   str,
			Tag:     "uuid",
		}
	}

	return nil
}

// validateJSON validates JSON format
func validateJSON(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "json validation requires string value",
			Tag:     "json",
		}
	}

	if !jsonx.IsValidJSON(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid JSON format",
			Value:   str,
			Tag:     "json",
		}
	}

	return nil
}

// validateIP validates IP address format
func validateIP(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "ip validation requires string value",
			Tag:     "ip",
		}
	}

	if !config.IsValidIP(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid IP address format",
			Value:   str,
			Tag:     "ip",
		}
	}

	return nil
}

// validateIPv4 validates IPv4 address format
func validateIPv4(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "ipv4 validation requires string value",
			Tag:     "ipv4",
		}
	}

	if !config.IsValidIPv4(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid IPv4 address format",
			Value:   str,
			Tag:     "ipv4",
		}
	}

	return nil
}

// validateIPv6 validates IPv6 address format
func validateIPv6(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "ipv6 validation requires string value",
			Tag:     "ipv6",
		}
	}

	if !config.IsValidIPv6(str) {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid IPv6 address format",
			Value:   str,
			Tag:     "ipv6",
		}
	}

	return nil
}

// validateDate validates date format (YYYY-MM-DD)
func validateDate(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "date validation requires string value",
			Tag:     "date",
		}
	}

	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return &FieldError{
			Field:   fieldName,
			Message: "invalid date format (expected YYYY-MM-DD)",
			Value:   str,
			Tag:     "date",
		}
	}

	return nil
}

// validateDateTime validates datetime format
func validateDateTime(fieldName string, value interface{}) *FieldError {
	str, ok := value.(string)
	if !ok {
		return &FieldError{
			Field:   fieldName,
			Message: "datetime validation requires string value",
			Tag:     "datetime",
		}
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, str); err == nil {
			return nil
		}
	}

	return &FieldError{
		Field:   fieldName,
		Message: "invalid datetime format",
		Value:   str,
		Tag:     "datetime",
	}
}
