package utils

func InArray[T comparable](value T, array []T) bool {
	for _, el := range array {
		if value == el {
			return true
		}
	}
	return false
}
