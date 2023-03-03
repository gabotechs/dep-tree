package language

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/require"
)

type ExportsResultBuilder map[string]*ExportsResult

func (e *ExportsResultBuilder) Build() map[string]*ExportsResult {
	return *e
}

func (e *ExportsResultBuilder) Entry(inId string, toId string, names ...string) *ExportsResultBuilder {
	var result *ExportsResult
	var ok bool
	if result, ok = (*e)[inId]; !ok {
		result = &ExportsResult{}
	}

	if len(names) == 1 && names[0] == "*" {
		result.Exports = append(result.Exports, ExportEntry{
			All: true,
			Id:  toId,
		})
	} else {
		var n []ExportName
		for _, name := range names {
			if strings.HasPrefix(name, "as ") {
				n[len(n)-1].Alias = strings.TrimLeft(name, "as ")
			} else {
				n = append(n, ExportName{Original: name})
			}
		}
		result.Exports = append(result.Exports, ExportEntry{
			Names: n,
			Id:    toId,
		})
	}
	(*e)[inId] = result
	return e
}

func b() *ExportsResultBuilder {
	return &ExportsResultBuilder{}
}

func TestParser_parseExports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	lang := TestLanguage{
		exports: b().
			Entry("1", "1", "A").
			Build(),
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
			Exports: b().
				Entry("1", "1", "A").
				Build(),
			Expected: makeStringOm(
				"A", "1",
			),
		},
		{
			Name: "one proxy",
			Id:   "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "2", "A").
				Build(),
			Expected: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double proxy",
			Id:   "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "A").
				Entry("3", "3", "A").
				Build(),
			Expected: makeStringOm(
				"A", "3",
			),
		},
		{
			Name: "double proxy with alias",
			Id:   "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "B", "as A").
				Entry("3", "3", "C", "as B").
				Build(),
			Expected: makeStringOm(
				"A", "3",
			),
		},
		{
			Name: "double proxy with all export",
			Id:   "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "*").
				Entry("3", "3", "C", "as A").
				Build(),
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
