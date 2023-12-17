package utils

func Max[T any](arr []T, f func(T) int) int {
	result := 0
	for _, el := range arr {
		v := f(el)
		if v > result {
			result = v
		}
	}
	return result
}
