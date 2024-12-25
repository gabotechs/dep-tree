package utils

// RemoveOverlap removes elements from a that are also on b.
func RemoveOverlap[T comparable](a, b []T) []T {
	res := make([]T, 0, len(a))
	bSet := SetFromSlice(b)
	for _, el := range a {
		if _, ok := bSet[el]; !ok {
			res = append(res, el)
		}
	}

	return res
}
