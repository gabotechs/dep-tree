package language

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParser_parseFile_IsCached(t *testing.T) {
	a := require.New(t)
	lang := TestLanguage{}

	parser := lang.testParser("1")

	start := time.Now()
	_, err := parser.parseFile("1")
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, err = parser.parseFile("1")
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}
