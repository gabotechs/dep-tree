package js

import (
	"path/filepath"
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"
)

func TestLanguage_Display(t *testing.T) {
	tests := []struct {
		Name     string
		Path     string
		Expected graph.DisplayResult
	}{
		{
			Name: "with a parent package.json",
			Path: filepath.Join(resolverTestFolder, "src", "utils", "sum.ts"),
			Expected: graph.DisplayResult{
				Name:  "src/utils/sum.ts",
				Group: "test-project",
			},
		},
		{
			Name: "with a parent package.json (same as above for checking cache)",
			Path: filepath.Join(resolverTestFolder, "src", "utils", "sum.ts"),
			Expected: graph.DisplayResult{
				Name:  "src/utils/sum.ts",
				Group: "test-project",
			},
		},
		{
			Name: "with two parent package.json, one without name",
			Path: filepath.Join(resolverTestFolder, "src", "module", "main.ts"),
			Expected: graph.DisplayResult{
				Name:  "src/module/main.ts",
				Group: "test-project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeJsLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)
			abs, _ := filepath.Abs(tt.Path)
			a.Equal(tt.Expected, lang.Display(abs))
		})
	}
}
