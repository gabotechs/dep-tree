package utils

import (
	"crypto/md5" //nolint:gosec
	"fmt"
)

func hash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text))) //nolint:gosec
}

var Hash = Cached(hash)
