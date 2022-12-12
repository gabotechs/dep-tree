package utils

func PrefixN(str string, char rune) int {
	for i, c := range str {
		if c != char {
			return i
		}
	}
	return len(str)
}
