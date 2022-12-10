package utils

func Clamp(min int, n, max int) int {
	switch {
	case n < min:
		return min
	case n > max:
		return max
	default:
		return n
	}
}
