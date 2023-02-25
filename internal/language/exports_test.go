package language

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParser_parseExports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	lang := TestLanguage{
		exports: map[string]*ExportsResult{
			"1": {
				Exports: map[string]string{"a": "1"},
			},
		},
	}

	parser, err := makeParser("1", func(_ string) (Language[TestLanguageData, TestFile], error) {
		return &lang, nil
	})
	a.NoError(err)

	start := time.Now()
	ctx, _, err = parser.CachedParseExports(ctx, "1")
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = parser.CachedParseExports(ctx, "1")
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}
