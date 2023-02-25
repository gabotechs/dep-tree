package language

import (
	"context"
	"errors"
	"path"
	"testing"
	"time"

	om "github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/graph"
)

const testFolder = ".parser_test"

type TestLanguageData struct{}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

func (t *TestLanguage) MakeNode(id string) (*graph.Node[TestLanguageData], error) {
	return &graph.Node[TestLanguageData]{
		Id:     id,
		Errors: make([]error, 0),
		Data:   TestLanguageData{},
	}, nil
}

func (t *TestLanguage) ParseImports(ctx context.Context, id string) (context.Context, *ImportsResult, error) {
	time.Sleep(time.Millisecond)
	if imports, ok := t.imports[id]; ok {
		return ctx, imports, nil
	} else {
		return ctx, imports, errors.New(id + " not found")
	}
}

func (t *TestLanguage) ParseExports(ctx context.Context, id string) (context.Context, *ExportsResult, error) {
	time.Sleep(time.Millisecond)
	if exports, ok := t.exports[id]; ok {
		return ctx, exports, nil
	} else {
		return ctx, exports, errors.New(id + " not found")
	}
}

func TestParser_Entrypoint(t *testing.T) {
	a := require.New(t)
	id := path.Join(testFolder, t.Name()+".js")

	parser, err := ParserBuilder(func(s string) (Language[TestLanguageData], error) {
		return &TestLanguage{}, nil
	})(id)
	a.NoError(err)
	entrypoint, err := parser.Entrypoint()
	a.NoError(err)

	a.NoError(err)
	a.Equal(id, entrypoint.Id)
}

func newOm(entries map[string][]string) *om.OrderedMap[string, []string] {
	m := om.NewOrderedMap[string, []string]()
	for k, v := range entries {
		m.Set(k, v)
	}
	return m
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name     string
		Id       string
		Imports  map[string]*ImportsResult
		Exports  map[string]*ExportsResult
		Expected []string
	}{
		{
			Name: "Simple",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: newOm(map[string][]string{
						"2": {"*"},
					}),
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: map[string]string{"Exported": "2"},
				},
			},
			Expected: []string{
				"2",
			},
		},
		{
			Name: "Proxy export",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: newOm(map[string][]string{
						"2": {"*"},
					}),
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: map[string]string{"Exported": "3"},
				},
				"3": {
					Exports: map[string]string{"Another-one": "3"},
				},
			},
			Expected: []string{
				"3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			parser, err := ParserBuilder(func(entrypoint string) (Language[TestLanguageData], error) {
				return &TestLanguage{
					imports: tt.Imports,
					exports: tt.Exports,
				}, nil
			})(tt.Id)

			a.NoError(err)
			node, err := parser.Entrypoint()
			a.NoError(err)
			_, deps, err := parser.Deps(context.Background(), node)
			a.NoError(err)
			result := make([]string, len(deps))
			for i, dep := range deps {
				result[i] = dep.Id
			}

			a.Equal(tt.Expected, result)
		})
	}
}
