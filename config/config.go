package config

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func IsValidPort(port string) bool {
	if port == "" {
		return false
	}
	portNum, err := strconv.Atoi(port)
	return err == nil && portNum > 0 && portNum <= 65535
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := GetEnv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnvInt returns an environment variable as int with a default value
func GetEnvInt(key string, defaultValue int) int {
	if value := GetEnv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvBool returns an environment variable as bool with a default value
func GetEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(GetEnv(key))
	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// GetEnvDuration returns an environment variable as time.Duration with a default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := GetEnv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetEnvFloat returns an environment variable as float64 with a default value
func GetEnvFloat(key string, defaultValue float64) float64 {
	if value := GetEnv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

// RequiredEnv returns an environment variable or panics if not set
func RequiredEnv(key string) string {
	value := GetEnv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// IsValidIP checks if a string is a valid IP address
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidIPv4 checks if a string is a valid IPv4 address
func IsValidIPv4(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.To4() != nil
}

// IsValidIPv6 checks if a string is a valid IPv6 address
func IsValidIPv6(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.To4() == nil
}

// IsValidHost checks if a string is a valid hostname or IP
func IsValidHost(host string) bool {
	if IsValidIP(host) {
		return true
	}
	// Check hostname format
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return hostnameRegex.MatchString(host) && len(host) <= 253
}

// IsValidDomain checks if a string is a valid domain name
func IsValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// IsValidBase64 checks if a string is valid base64
func IsValidBase64(base64Str string) bool {
	_, err := base64.StdEncoding.DecodeString(base64Str)
	return err == nil
}

// IsValidHex checks if a string contains only hexadecimal characters
func IsValidHex(hex string) bool {
	hexRegex := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	return hexRegex.MatchString(hex)
}

// IsValidSlug checks if a string is a valid URL slug
func IsValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 100 {
		return false
	}
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return slugRegex.MatchString(slug)
}

// IsValidNumber checks if a string represents a valid number
func IsValidNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// IsValidInteger checks if a string represents a valid integer
func IsValidInteger(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// IsValidPositiveInteger checks if a string represents a valid positive integer
func IsValidPositiveInteger(str string) bool {
	num, err := strconv.Atoi(str)
	return err == nil && num > 0
}

// IsValidNonNegativeInteger checks if a string represents a valid non-negative integer
func IsValidNonNegativeInteger(str string) bool {
	num, err := strconv.Atoi(str)
	return err == nil && num >= 0
}

// IsValidFloat checks if a string represents a valid float
func IsValidFloat(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// IsValidPositiveFloat checks if a string represents a valid positive float
func IsValidPositiveFloat(str string) bool {
	num, err := strconv.ParseFloat(str, 64)
	return err == nil && num > 0
}

// IsValidNonNegativeFloat checks if a string represents a valid non-negative float
func IsValidNonNegativeFloat(str string) bool {
	num, err := strconv.ParseFloat(str, 64)
	return err == nil && num >= 0
}

// IsValidNonEmptyString checks if a string is not empty and not just whitespace
func IsValidNonEmptyString(str string) bool {
	return strings.TrimSpace(str) != ""
}

// IsValidAlphabetic checks if a string contains only alphabetic characters
func IsValidAlphabetic(str string) bool {
	if str == "" {
		return false
	}
	alphabeticRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return alphabeticRegex.MatchString(str)
}

// IsValidAlphanumeric checks if a string contains only alphanumeric characters
func IsValidAlphanumeric(str string) bool {
	if str == "" {
		return false
	}
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphanumericRegex.MatchString(str)
}

// IsValidNumeric checks if a string contains only numeric characters
func IsValidNumeric(str string) bool {
	if str == "" {
		return false
	}
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return numericRegex.MatchString(str)
}

// IsValidLowercase checks if a string contains only lowercase characters
func IsValidLowercase(str string) bool {
	if str == "" {
		return false
	}
	return str == strings.ToLower(str)
}

// IsValidUppercase checks if a string contains only uppercase characters
func IsValidUppercase(str string) bool {
	if str == "" {
		return false
	}
	return str == strings.ToUpper(str)
}

// SubstituteEnvVars replaces ${VARIABLE} patterns with environment variable values
// Supports ${VARIABLE}, ${VARIABLE:-default}, and ${VARIABLE=default} syntax
// Both :- and = use the default value if the variable is unset or empty
func SubstituteEnvVars(content []byte) []byte {
	// Regex to match ${VARIABLE} or ${VARIABLE:-default} or ${VARIABLE=default}
	envRegex := regexp.MustCompile(`\$\{([^}]+)\}`)

	return envRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		// Extract the variable name and optional default value
		inner := string(match[2 : len(match)-1]) // Remove ${ and }

		var varName, defaultValue string
		if strings.Contains(inner, ":-") {
			parts := strings.SplitN(inner, ":-", 2)
			varName = parts[0]
			defaultValue = parts[1]
		} else if strings.Contains(inner, "=") {
			parts := strings.SplitN(inner, "=", 2)
			varName = parts[0]
			defaultValue = parts[1]
		} else {
			varName = inner
		}

		// Get environment variable value
		value := GetEnv(varName)
		if value == "" && defaultValue != "" {
			value = defaultValue
		}

		return []byte(value)
	})
}

// LoadYAMLConfig loads and parses a YAML config file with environment variable substitution
func LoadYAMLConfig(filename string, target interface{}) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}

	content := SubstituteEnvVars(file)
	if err := yaml.Unmarshal(content, target); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return nil
}
