package utils

import (
	"github.com/bmatcuk/doublestar/v4"
)

func GlobstarMatch(pattern string, check string) (bool, error) {
	return doublestar.PathMatch(pattern, check)
}
