package utils

func AppendFront[T any](el T, arr []T) []T {
	result := make([]T, len(arr)+1)
	result[0] = el
	for i, prev := range arr {
		result[i+1] = prev
	}
	return result
}
