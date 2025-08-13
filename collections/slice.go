package collections

import (
	"fmt"
	"reflect"
)

// Contains checks if a slice contains a specific element
func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// ContainsAny checks if a slice contains any of the given elements
func ContainsAny[T comparable](slice []T, elements ...T) bool {
	for _, element := range elements {
		if Contains(slice, element) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of element in slice
func IndexOf[T comparable](slice []T, element T) int {
	for i, item := range slice {
		if item == element {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of element in slice
func LastIndexOf[T comparable](slice []T, element T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

// Remove removes the first occurrence of element from slice
func Remove[T comparable](slice []T, element T) []T {
	index := IndexOf(slice, element)
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// RemoveAll removes all occurrences of element from slice
func RemoveAll[T comparable](slice []T, element T) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if item != element {
			result = append(result, item)
		}
	}
	return result
}

// RemoveAt removes element at the specified index
func RemoveAt[T any](slice []T, index int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return slice, fmt.Errorf("index %d out of bounds for slice of length %d", index, len(slice))
	}
	return append(slice[:index], slice[index+1:]...), nil
}

// Insert inserts element at the specified index
func Insert[T any](slice []T, index int, element T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return slice, fmt.Errorf("index %d out of bounds for slice of length %d", index, len(slice))
	}

	// Extend slice by one element
	slice = append(slice, *new(T))

	// Shift elements to the right
	copy(slice[index+1:], slice[index:])

	// Insert new element
	slice[index] = element

	return slice, nil
}

// Unique removes duplicate elements from slice
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Reverse reverses the slice in place
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Reversed returns a new reversed slice
func Reversed[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[len(slice)-1-i] = item
	}
	return result
}

// Filter returns a new slice containing only elements that match the predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms each element of the slice using the provided function
func Map[T, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = transform(item)
	}
	return result
}

// Reduce reduces the slice to a single value using the provided function
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Find returns the first element that matches the predicate
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element that matches the predicate
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// All checks if all elements match the predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Any checks if any element matches the predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Count counts elements that match the predicate
func Count[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, item := range slice {
		if predicate(item) {
			count++
		}
	}
	return count
}

// Chunk splits slice into chunks of specified size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	chunks := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Flatten flattens a slice of slices into a single slice
func Flatten[T any](slices [][]T) []T {
	totalLen := 0
	for _, slice := range slices {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)
	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}

// Intersection returns elements common to both slices
func Intersection[T comparable](slice1, slice2 []T) []T {
	set1 := make(map[T]bool)
	for _, item := range slice1 {
		set1[item] = true
	}

	result := make([]T, 0)
	seen := make(map[T]bool)

	for _, item := range slice2 {
		if set1[item] && !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

// Union returns all unique elements from both slices
func Union[T comparable](slice1, slice2 []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice1)+len(slice2))

	for _, item := range slice1 {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	for _, item := range slice2 {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

// Difference returns elements in slice1 that are not in slice2
func Difference[T comparable](slice1, slice2 []T) []T {
	set2 := make(map[T]bool)
	for _, item := range slice2 {
		set2[item] = true
	}

	result := make([]T, 0)
	for _, item := range slice1 {
		if !set2[item] {
			result = append(result, item)
		}
	}

	return result
}

// SlicesEqual checks if two slices are equal
func SlicesEqual[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, item := range slice1 {
		if item != slice2[i] {
			return false
		}
	}

	return true
}

// IsEmpty checks if slice is empty
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// IsNotEmpty checks if slice is not empty
func IsNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}

// Clone creates a shallow copy of the slice
func Clone[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

// SliceToInterface converts a typed slice to []interface{}
func SliceToInterface[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, item := range slice {
		result[i] = item
	}
	return result
}

// InterfaceToSlice converts []interface{} to a typed slice (with type checking)
func InterfaceToSlice[T any](interfaces []interface{}) ([]T, error) {
	result := make([]T, len(interfaces))
	var zero T
	targetType := reflect.TypeOf(zero)

	for i, item := range interfaces {
		if item == nil {
			result[i] = zero
			continue
		}

		itemType := reflect.TypeOf(item)
		if !itemType.AssignableTo(targetType) {
			return nil, fmt.Errorf("cannot convert item at index %d: type %v is not assignable to %v", i, itemType, targetType)
		}

		result[i] = item.(T)
	}

	return result, nil
}
