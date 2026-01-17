package slices

import "slices"

// Contains returns true if the slice contains the given element.
// Uses comparable constraint for type-safe equality comparison.
func Contains[T comparable](slice []T, element T) bool {
	return slices.Contains(slice, element)
}

// ContainsAny returns true if the slice contains any of the given elements.
func ContainsAny[T comparable](slice []T, elements []T) bool {
	for _, element := range elements {
		if Contains(slice, element) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if the slice contains all of the given elements.
func ContainsAll[T comparable](slice []T, elements []T) bool {
	for _, element := range elements {
		if !Contains(slice, element) {
			return false
		}
	}
	return true
}

// IndexOf returns the index of the first occurrence of element in the slice,
// or -1 if the element is not found.
func IndexOf[T comparable](slice []T, element T) int {
	for i, item := range slice {
		if item == element {
			return i
		}
	}
	return -1
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms each element in the slice using the given function.
func Map[T any, R any](slice []T, transform func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = transform(item)
	}
	return result
}

// Unique returns a new slice with duplicate elements removed.
// Preserves the order of first occurrence.
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// First returns the first element of the slice and true, or zero value and false if empty.
func First[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Last returns the last element of the slice and true, or zero value and false if empty.
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Copy returns a shallow copy of the slice.
func Copy[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}
