package utils

func Clamp(min int, n, max int) int {
	if n < min {
		return min
	} else if n > max {
		return max
	} else {
		return n
	}
}
