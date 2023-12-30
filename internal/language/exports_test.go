package language

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/require"
)

type ExportsResultBuilder map[string]*ExportsEntries

func (e *ExportsResultBuilder) Build() map[string]*ExportsEntries {
	return *e
}

func (e *ExportsResultBuilder) Entry(inId string, toId string, names ...string) *ExportsResultBuilder {
	var result *ExportsEntries
	var ok bool
	if result, ok = (*e)[inId]; !ok {
		result = &ExportsEntries{}
	}

	if len(names) == 1 && names[0] == "*" {
		result.Exports = append(result.Exports, ExportEntry{
			All:  true,
			Path: toId,
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
			Path:  toId,
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
	lang := TestLanguage{
		exports: b().
			Entry("1", "1", "A").
			Build(),
	}

	parser := lang.testParser("1")

	start := time.Now()
	_, err := parser.parseExports("1", false, nil)
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, err = parser.parseExports("1", false, nil)
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
		Name              string
		Path              string
		Exports           map[string]*ExportsEntries
		ExpectedUnwrapped *orderedmap.OrderedMap[string, string]
		ExpectedWrapped   *orderedmap.OrderedMap[string, string]
		ExpectedErrors    []string
	}{
		{
			Name: "direct export",
			Path: "1",
			Exports: b().
				Entry("1", "1", "A").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "1",
			),
			ExpectedWrapped: makeStringOm(
				"A", "1",
			),
		},
		{
			Name: "one proxy",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "2", "A").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "2",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double proxy",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "A").
				Entry("3", "3", "A").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "3",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double proxy with alias",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "B", "as A").
				Entry("3", "3", "C", "as B").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "3",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double proxy with all export",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "*").
				Entry("3", "3", "C", "as A").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "3",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
			),
		},
		{
			Name: "double all export and single named export",
			Path: "1",
			Exports: b().
				Entry("1", "2", "*").
				Entry("2", "2", "D").
				Entry("2", "3", "*").
				Entry("3", "3", "A", "B", "C").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"D", "2",
				"A", "3",
				"B", "3",
				"C", "3",
			),
			ExpectedWrapped: makeStringOm(
				"D", "2",
				"A", "2",
				"B", "2",
				"C", "2",
			),
		},
		{
			Name: "name not found in destination",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("3", "2", "B").
				Build(),
			ExpectedUnwrapped: makeStringOm(),
			ExpectedWrapped:   makeStringOm(),
			ExpectedErrors:    []string{"2 not found"},
		},
		{
			Name: "circular export",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("2", "3", "B", "as A").
				Entry("3", "4", "C", "as B").
				Entry("4", "1", "A", "as C").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				// TODO: I don't know if this is right...
				"A", "4",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
			),
			ExpectedErrors: []string{
				"circular export: cycle detected:\n1\n2\n3\n4\n1",
			},
		},
		{
			Name: "non circular export",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("1", "3", "A").
				Entry("1", "4", "A").
				Entry("2", "2", "A").
				Entry("3", "2", "A").
				Entry("4", "3", "A").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "2",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
				"A", "3",
				"A", "4",
			),
			ExpectedErrors: []string{},
		},
		{
			Name: "non circular export (2)",
			Path: "1",
			Exports: b().
				Entry("1", "2", "A").
				Entry("1", "2", "B").
				Entry("1", "2", "C").
				Entry("2", "2", "A").
				Entry("2", "2", "B").
				Entry("2", "2", "C").
				Build(),
			ExpectedUnwrapped: makeStringOm(
				"A", "2",
				"B", "2",
				"C", "2",
			),
			ExpectedWrapped: makeStringOm(
				"A", "2",
				"B", "2",
				"C", "2",
			),
			ExpectedErrors: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			lang := &TestLanguage{
				exports: tt.Exports,
			}
			parser := lang.testParser(tt.Path)

			for i := 0; i < 2; i++ {
				exports, err := parser.parseExports("1", true, nil)
				a.NoError(err)

				a.Equal(tt.ExpectedUnwrapped, exports.Exports)

				exports, err = parser.parseExports("1", false, nil)
				a.NoError(err)

				a.Equal(tt.ExpectedWrapped, exports.Exports)
				var expectedErrors []error
				for _, expectedError := range tt.ExpectedErrors {
					expectedErrors = append(expectedErrors, errors.New(expectedError))
				}
				a.Equal(expectedErrors, exports.Errors)
			}
		})
	}
}
