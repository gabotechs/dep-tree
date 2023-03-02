package language

import (
	"context"
	"testing"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/require"
)

func TestParser_parseExports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	lang := TestLanguage{
		exports: map[string]*ExportsResult{
			"1": {
				Exports: []ExportEntry{{
					Names: []ExportName{{Original: "A"}},
					Id:    "1",
				}},
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

func makeStringOm(args ...string) *orderedmap.OrderedMap[string, string] {
	om := orderedmap.NewOrderedMap[string, string]()
	for i := 0; i < len(args); i += 2 {
		om.Set(args[i], args[i+1])
	}
	return om
}

func TestParser_CachedUnwrappedParseExports(t *testing.T) {
	tests := []struct {
		Name     string
		Id       string
		Exports  map[string]*ExportsResult
		Expected *orderedmap.OrderedMap[string, string]
	}{
		{
			Name: "direct export",
			Id:   "1",
			Exports: map[string]*ExportsResult{
				"1": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "1",
					},
				}},
			},
			Expected: makeStringOm(
				"A", "1",
			),
		},
		{
			Name: "one proxy",
			Id:   "1",
			Exports: map[string]*ExportsResult{
				"1": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "2",
					},
				}},
				"2": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "2",
					},
				}},
			},
			Expected: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double proxy",
			Id:   "1",
			Exports: map[string]*ExportsResult{
				"1": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "2",
					},
				}},
				"2": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "3",
					},
				}},
				"3": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "3",
					},
				}},
			},
			Expected: makeStringOm(
				"A", "3",
			),
		},
		{
			Name: "double proxy with alias",
			Id:   "1",
			Exports: map[string]*ExportsResult{
				"1": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "2",
					},
				}},
				"2": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "B", Alias: "A"}},
						Id:    "3",
					},
				}},
				"3": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "C", Alias: "B"}},
						Id:    "3",
					},
				}},
			},
			Expected: makeStringOm(
				"A", "3",
			),
		},
		{
			Name: "double proxy with all export",
			Id:   "1",
			Exports: map[string]*ExportsResult{
				"1": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "A"}},
						Id:    "2",
					},
				}},
				"2": {Exports: []ExportEntry{
					{
						All: true,
						Id:  "3",
					},
				}},
				"3": {Exports: []ExportEntry{
					{
						Names: []ExportName{{Original: "C", Alias: "A"}},
						Id:    "3",
					},
				}},
			},
			Expected: makeStringOm(
				"A", "3",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			parser, err := makeParser(tt.Id, func(_ string) (Language[TestLanguageData, TestFile], error) {
				return &TestLanguage{
					exports: tt.Exports,
				}, nil
			})
			a.NoError(err)

			_, exports, err := parser.CachedUnwrappedParseExports(context.Background(), "1")
			a.NoError(err)

			a.Equal(tt.Expected, exports.Exports)
			a.Nil(exports.Errors)
		})
	}
}
