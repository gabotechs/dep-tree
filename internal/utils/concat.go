package utils

func Concat[T any](batches ...interface{}) []T {
	result := make([]T, 0)
	for _, batch := range batches {
		switch b := batch.(type) {
		case []T:
			result = append(result, b...)
		case T:
			result = append(result, b)
		}
	}
	return result
}
