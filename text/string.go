package text

import (
	"regexp"
	"strings"
	"unicode"
)

// Normalize converts string to lowercase and trims whitespace
func Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

// ToSnakeCase converts string to snake_case
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// ToCamelCase converts string to camelCase
func ToCamelCase(str string) string {
	if str == "" {
		return ""
	}

	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(strings.ToLower(words[0]))

	for _, word := range words[1:] {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])) + strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

// ToPascalCase converts string to PascalCase
func ToPascalCase(str string) string {
	camel := ToCamelCase(str)
	if len(camel) == 0 {
		return ""
	}
	return strings.ToUpper(string(camel[0])) + camel[1:]
}

// ToKebabCase converts string to kebab-case
func ToKebabCase(str string) string {
	return strings.ReplaceAll(ToSnakeCase(str), "_", "-")
}

// Truncate truncates a string to the specified length
func Truncate(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length]
}

// TruncateWithEllipsis truncates a string and adds ellipsis
func TruncateWithEllipsis(str string, length int) string {
	if len(str) <= length {
		return str
	}
	if length <= 3 {
		return str[:length]
	}
	return str[:length-3] + "..."
}

// Reverse reverses a string
func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsEmpty checks if string is empty or contains only whitespace
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// ContainsAny checks if string contains any of the substrings
func ContainsAny(str string, substrings ...string) bool {
	for _, substring := range substrings {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}

// StartsWithAny checks if string starts with any of the prefixes
func StartsWithAny(str string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

// EndsWithAny checks if string ends with any of the suffixes
func EndsWithAny(str string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

// RemoveAccents removes accents from characters
func RemoveAccents(str string) string {
	// Simple ASCII transliteration
	replacements := map[rune]string{
		'á': "a", 'à': "a", 'ä': "a", 'â': "a", 'ā': "a", 'ă': "a", 'ą': "a",
		'é': "e", 'è': "e", 'ë': "e", 'ê': "e", 'ē': "e", 'ĕ': "e", 'ė': "e", 'ę': "e",
		'í': "i", 'ì': "i", 'ï': "i", 'î': "i", 'ī': "i", 'ĭ': "i", 'į': "i",
		'ó': "o", 'ò': "o", 'ö': "o", 'ô': "o", 'ō': "o", 'ŏ': "o", 'ő': "o",
		'ú': "u", 'ù': "u", 'ü': "u", 'û': "u", 'ū': "u", 'ŭ': "u", 'ű': "u", 'ų': "u",
		'ñ': "n", 'ń': "n", 'ň': "n", 'ņ': "n",
		'ç': "c", 'ć': "c", 'č': "c", 'ĉ': "c", 'ċ': "c",
		'ý': "y", 'ÿ': "y", 'ŷ': "y",
		'ž': "z", 'ź': "z", 'ż': "z",
		'š': "s", 'ś': "s", 'ş': "s", 'ŝ': "s",
		'đ': "d", 'ď': "d",
		'ř': "r", 'ŕ': "r", 'ŗ': "r",
		'ł': "l", 'ľ': "l", 'ŀ': "l", 'ļ': "l",
		'ť': "t", 'ţ': "t",
		'ğ': "g", 'ĝ': "g", 'ġ': "g", 'ģ': "g",
		'ĥ': "h", 'ħ': "h",
		'ĵ': "j",
		'ķ': "k",
		'ĺ': "l",
		'ŵ': "w",
	}

	var result strings.Builder
	for _, r := range str {
		if replacement, exists := replacements[unicode.ToLower(r)]; exists {
			if unicode.IsUpper(r) {
				result.WriteString(strings.ToUpper(replacement))
			} else {
				result.WriteString(replacement)
			}
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToSlug converts string to URL-friendly slug
func ToSlug(str string) string {
	// Remove accents
	slug := RemoveAccents(str)

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// WordCount counts words in a string
func WordCount(str string) int {
	fields := strings.Fields(str)
	return len(fields)
}

// CharCount counts characters in a string (excluding spaces)
func CharCount(str string) int {
	count := 0
	for _, r := range str {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}

// Pad pads string to specified length with given character
func Pad(str string, length int, padChar rune, leftPad bool) string {
	if len(str) >= length {
		return str
	}

	padding := strings.Repeat(string(padChar), length-len(str))
	if leftPad {
		return padding + str
	}
	return str + padding
}

// PadLeft pads string on the left to specified length
func PadLeft(str string, length int, padChar rune) string {
	return Pad(str, length, padChar, true)
}

// PadRight pads string on the right to specified length
func PadRight(str string, length int, padChar rune) string {
	return Pad(str, length, padChar, false)
}

// ExtractNumbers extracts all numbers from a string
func ExtractNumbers(str string) []string {
	re := regexp.MustCompile(`\d+`)
	return re.FindAllString(str, -1)
}

// ExtractEmails extracts all email addresses from a string
func ExtractEmails(str string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return re.FindAllString(str, -1)
}

// ExtractURLs extracts all URLs from a string
func ExtractURLs(str string) []string {
	re := regexp.MustCompile(`https?://[^\s]+`)
	return re.FindAllString(str, -1)
}

// Mask masks part of a string with asterisks
func Mask(str string, start, end int) string {
	if start < 0 || end > len(str) || start >= end {
		return str
	}

	masked := str[:start] + strings.Repeat("*", end-start) + str[end:]
	return masked
}

// MaskEmail masks email address (keeps first and last char of username)
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return "*@" + domain
	}

	maskedUsername := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
	return maskedUsername + "@" + domain
}

// LevenshteinDistance calculates the Levenshtein distance between two strings
func LevenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}

	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
