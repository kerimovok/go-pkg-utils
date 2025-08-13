package collections

import (
	"fmt"
	"reflect"
)

// Keys returns all keys from the map
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from the map
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// HasKey checks if map contains the specified key
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

// HasValue checks if map contains the specified value
func HasValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

// GetOrDefault returns the value for key, or defaultValue if key doesn't exist
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// GetOrCreate returns the value for key, or creates and returns defaultValue if key doesn't exist
func GetOrCreate[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, exists := m[key]; exists {
		return value
	}
	m[key] = defaultValue
	return defaultValue
}

// IsEmpty checks if map is empty
func IsEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// IsNotEmpty checks if map is not empty
func IsNotEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) > 0
}

// Clone creates a shallow copy of the map
func CloneMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Merge merges multiple maps into a new map (later maps override earlier ones)
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	totalSize := 0
	for _, m := range maps {
		totalSize += len(m)
	}

	result := make(map[K]V, totalSize)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Filter returns a new map containing only key-value pairs that match the predicate
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms all values in the map using the provided function
func MapValues[K comparable, V, U any](m map[K]V, transform func(V) U) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = transform(v)
	}
	return result
}

// MapKeys transforms all keys in the map using the provided function
func MapKeys[K comparable, V any, L comparable](m map[K]V, transform func(K) L) map[L]V {
	result := make(map[L]V, len(m))
	for k, v := range m {
		result[transform(k)] = v
	}
	return result
}

// Invert swaps keys and values (values must be comparable)
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// MapsEqual checks if two maps are equal
func MapsEqual[K comparable, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, exists := m2[k]; !exists || v1 != v2 {
			return false
		}
	}

	return true
}

// Pick returns a new map with only the specified keys
func Pick[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	result := make(map[K]V)
	for _, key := range keys {
		if value, exists := m[key]; exists {
			result[key] = value
		}
	}
	return result
}

// Omit returns a new map without the specified keys
func Omit[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	omitSet := make(map[K]bool, len(keys))
	for _, key := range keys {
		omitSet[key] = true
	}

	result := make(map[K]V)
	for k, v := range m {
		if !omitSet[k] {
			result[k] = v
		}
	}
	return result
}

// GroupBy groups slice elements by the result of the key function
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// CountBy counts slice elements by the result of the key function
func CountBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	for _, item := range slice {
		key := keyFunc(item)
		result[key]++
	}
	return result
}

// ToSlice converts map to slice of key-value pairs
func ToSlice[K comparable, V any](m map[K]V) []KeyValue[K, V] {
	result := make([]KeyValue[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, KeyValue[K, V]{Key: k, Value: v})
	}
	return result
}

// FromSlice creates map from slice of key-value pairs
func FromSlice[K comparable, V any](slice []KeyValue[K, V]) map[K]V {
	result := make(map[K]V, len(slice))
	for _, kv := range slice {
		result[kv.Key] = kv.Value
	}
	return result
}

// KeyValue represents a key-value pair
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// Reduce reduces map to a single value using the provided function
func ReduceMap[K comparable, V, U any](m map[K]V, initial U, reducer func(U, K, V) U) U {
	result := initial
	for k, v := range m {
		result = reducer(result, k, v)
	}
	return result
}

// ForEach executes a function for each key-value pair
func ForEach[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Any checks if any key-value pair matches the predicate
func AnyMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// All checks if all key-value pairs match the predicate
func AllMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if !predicate(k, v) {
			return false
		}
	}
	return true
}

// Find returns the first key-value pair that matches the predicate
func FindInMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) (K, V, bool) {
	for k, v := range m {
		if predicate(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// MapToInterface converts a map to map[string]interface{}
func MapToInterface[K comparable, V any](m map[K]V) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[fmt.Sprintf("%v", k)] = v
	}
	return result
}

// InterfaceToMap converts map[string]interface{} to a typed map (with type checking)
func InterfaceToMap[K comparable, V any](m map[string]interface{}, keyConverter func(string) (K, error)) (map[K]V, error) {
	result := make(map[K]V, len(m))
	var zeroV V
	targetType := reflect.TypeOf(zeroV)

	for strKey, value := range m {
		// Convert key
		key, err := keyConverter(strKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert key %q: %w", strKey, err)
		}

		// Convert value
		if value == nil {
			result[key] = zeroV
			continue
		}

		valueType := reflect.TypeOf(value)
		if !valueType.AssignableTo(targetType) {
			return nil, fmt.Errorf("cannot convert value for key %v: type %v is not assignable to %v", key, valueType, targetType)
		}

		result[key] = value.(V)
	}

	return result, nil
}

// Intersection returns a map containing keys common to both maps
func MapIntersection[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m1 {
		if _, exists := m2[k]; exists {
			result[k] = v
		}
	}
	return result
}

// Difference returns a map containing keys in m1 that are not in m2
func MapDifference[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m1 {
		if _, exists := m2[k]; !exists {
			result[k] = v
		}
	}
	return result
}
