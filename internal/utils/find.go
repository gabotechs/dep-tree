package utils

func Find[T any](arr []T, f func(T) bool) *T {
	for _, el := range arr {
		if f(el) {
			return &el
		}
	}
	return nil
}
