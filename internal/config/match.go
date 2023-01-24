package config

import (
	"github.com/bmatcuk/doublestar/v4"
)

func match(pattern string, check string) (bool, error) {
	return doublestar.Match(pattern, check)
}
