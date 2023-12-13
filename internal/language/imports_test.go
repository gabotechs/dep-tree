package language

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParser_parseImports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	lang := TestLanguage{
		imports: map[string]*ImportsResult{
			"1": {
				Imports: []ImportEntry{
					{All: true, Path: "2"},
				},
			},
		},
	}
	parser := lang.testParser("1")

	start := time.Now()
	ctx, _, err := parser.CachedParseImports(ctx, "1")
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = parser.CachedParseImports(ctx, "1")
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}
