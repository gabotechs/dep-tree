package utils

func InLimits(i int, arr []any) bool {
	return i >= 0 && i < len(arr)
}
