package utils

import (
	"crypto/md5"
	"fmt"
)

func hash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

var Hash = Cached(hash)
