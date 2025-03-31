package slices

func RemoveFromSlice[T comparable](slice *[]T, element T) {
	for i, e := range *slice {
		if e == element {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}

type Comparable interface {
	Compare(any) bool
}

func RemoveFromSliceComparable[T Comparable](slice *[]T, element T) {
	for i, e := range *slice {
		if element.Compare(e) {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}
