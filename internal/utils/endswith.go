package utils

import "strings"

func EndsWith(string string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.HasSuffix(string, substring) {
			return true
		}
	}
	return false
}
