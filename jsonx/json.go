package jsonx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Marshal marshals v to JSON with error handling
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent marshals v to JSON with indentation
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// MarshalPretty marshals v to pretty-printed JSON
func MarshalPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// Unmarshal unmarshals JSON data into v
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalFromString unmarshals JSON string into v
func UnmarshalFromString(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// UnmarshalFromReader unmarshals JSON from io.Reader into v
func UnmarshalFromReader(reader io.Reader, v interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(v)
}

// ToJSON converts any value to JSON string
func ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToPrettyJSON converts any value to pretty JSON string
func ToPrettyJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON converts JSON string to specified type
func FromJSON[T any](jsonStr string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// IsValidJSON checks if a string is valid JSON
func IsValidJSON(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// Minify removes unnecessary whitespace from JSON
func Minify(jsonBytes []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, jsonBytes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MinifyString removes unnecessary whitespace from JSON string
func MinifyString(jsonStr string) (string, error) {
	minified, err := Minify([]byte(jsonStr))
	if err != nil {
		return "", err
	}
	return string(minified), nil
}

// DeepCopy performs a deep copy of an object using JSON serialization
func DeepCopy[T any](src T) (T, error) {
	var dst T
	data, err := json.Marshal(src)
	if err != nil {
		return dst, err
	}
	err = json.Unmarshal(data, &dst)
	return dst, err
}

// GetValue gets a value from JSON using dot notation path
func GetValue(jsonData map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	current := jsonData

	for i, key := range keys {
		if i == len(keys)-1 {
			return current[key], nil
		}

		next, ok := current[key]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in path '%s'", key, path)
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("value at key '%s' is not an object", key)
		}

		current = nextMap
	}

	return nil, fmt.Errorf("empty path")
}

// SetValue sets a value in JSON using dot notation path
func SetValue(jsonData map[string]interface{}, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	current := jsonData

	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = value
			return nil
		}

		next, exists := current[key]
		if !exists {
			next = make(map[string]interface{})
			current[key] = next
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("value at key '%s' is not an object", key)
		}

		current = nextMap
	}

	return fmt.Errorf("empty path")
}

// DeleteValue deletes a value from JSON using dot notation path
func DeleteValue(jsonData map[string]interface{}, path string) error {
	keys := strings.Split(path, ".")
	if len(keys) == 0 {
		return fmt.Errorf("empty path")
	}

	if len(keys) == 1 {
		delete(jsonData, keys[0])
		return nil
	}

	current := jsonData
	for _, key := range keys[:len(keys)-1] {
		next, ok := current[key]
		if !ok {
			return fmt.Errorf("key '%s' not found in path '%s'", key, path)
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("value at key '%s' is not an object", key)
		}

		current = nextMap
	}

	delete(current, keys[len(keys)-1])
	return nil
}

// HasPath checks if a path exists in JSON data
func HasPath(jsonData map[string]interface{}, path string) bool {
	_, err := GetValue(jsonData, path)
	return err == nil
}

// Flatten flattens nested JSON into a flat map with dot notation keys
func Flatten(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive(data, "", result)
	return result
}

func flattenRecursive(data map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]interface{}); ok {
			flattenRecursive(nested, newKey, result)
		} else {
			result[newKey] = value
		}
	}
}

// Unflatten converts a flat map with dot notation keys back to nested JSON
func Unflatten(flat map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range flat {
		_ = SetValue(result, key, value)
	}

	return result
}

// Merge merges multiple JSON objects into one
func Merge(objects ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, obj := range objects {
		for key, value := range obj {
			result[key] = value
		}
	}

	return result
}

// DeepMerge deeply merges multiple JSON objects
func DeepMerge(objects ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, obj := range objects {
		deepMergeInto(result, obj)
	}

	return result
}

func deepMergeInto(dst, src map[string]interface{}) {
	for key, srcValue := range src {
		if dstValue, exists := dst[key]; exists {
			if dstMap, ok := dstValue.(map[string]interface{}); ok {
				if srcMap, ok := srcValue.(map[string]interface{}); ok {
					deepMergeInto(dstMap, srcMap)
					continue
				}
			}
		}
		dst[key] = srcValue
	}
}

// Equal compares two JSON values for equality
func Equal(a, b interface{}) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(aJSON, bJSON)
}

// ConvertType converts a value to the specified type using JSON round-trip
func ConvertType[T any](value interface{}) (T, error) {
	var result T

	// If the value is already the correct type, return it directly
	if v, ok := value.(T); ok {
		return v, nil
	}

	// Use JSON round-trip for conversion
	data, err := json.Marshal(value)
	if err != nil {
		return result, fmt.Errorf("failed to marshal value: %w", err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
	}

	return result, nil
}

// GetString safely gets a string value from JSON data
func GetString(data map[string]interface{}, key string) (string, error) {
	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt safely gets an int value from JSON data
func GetInt(data map[string]interface{}, key string) (int, error) {
	value, exists := data[key]
	if !exists {
		return 0, fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// GetFloat safely gets a float64 value from JSON data
func GetFloat(data map[string]interface{}, key string) (float64, error) {
	value, exists := data[key]
	if !exists {
		return 0, fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// GetBool safely gets a bool value from JSON data
func GetBool(data map[string]interface{}, key string) (bool, error) {
	value, exists := data[key]
	if !exists {
		return false, fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// GetArray safely gets an array value from JSON data
func GetArray(data map[string]interface{}, key string) ([]interface{}, error) {
	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case []interface{}:
		return v, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("value at key '%s' is not an array", key)
	}
}

// GetObject safely gets an object value from JSON data
func GetObject(data map[string]interface{}, key string) (map[string]interface{}, error) {
	value, exists := data[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("value at key '%s' is not an object", key)
	}
}

// Size returns the size of JSON data (number of keys for objects, length for arrays)
func Size(data interface{}) int {
	switch v := data.(type) {
	case map[string]interface{}:
		return len(v)
	case []interface{}:
		return len(v)
	case string:
		return len(v)
	default:
		return 0
	}
}

// Keys returns all keys from a JSON object
func Keys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values from a JSON object
func Values(data map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(data))
	for _, value := range data {
		values = append(values, value)
	}
	return values
}

// ToMap converts any struct to map[string]interface{} using JSON tags
func ToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// FromMap converts map[string]interface{} to any struct using JSON tags
func FromMap[T any](data map[string]interface{}) (T, error) {
	var result T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(jsonData, &result)
	return result, err
}
