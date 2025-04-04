package sliceutils

// RemoveFromSlice removes the first occurrence of an element from the slice.
// Optimization: In-place removal with single allocation.
func RemoveFromSlice[T comparable](slice *[]T, element T) {
	for i, e := range *slice {
		if e == element {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}

// Comparable defines an interface for comparable types.
type Comparable interface {
	Compare(any) bool
}

// RemoveFromSliceComparable removes the first matching element using Compare.
// Optimization: Similar to RemoveFromSlice but with custom comparison.
func RemoveFromSliceComparable[T Comparable](slice *[]T, element T) {
	for i, e := range *slice {
		if element.Compare(e) {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}
