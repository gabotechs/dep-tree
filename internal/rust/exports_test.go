package rust

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

func TestLanguage_ParseExports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(path.Join(testFolder))

	tests := []struct {
		Name     string
		Expected []language.ExportEntry
		Errors   []error
	}{
		{
			Name: "lib.rs",
			Expected: []language.ExportEntry{
				{
					Names: []language.ExportName{{Original: "div"}},
					Id:    path.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					Names: []language.ExportName{{Original: "abs"}},
					Id:    path.Join(absTestFolder, "src", "abs", "abs.rs"),
				},
				{
					Names: []language.ExportName{{Original: "div"}},
					Id:    path.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					Names: []language.ExportName{{Original: "avg"}},
					Id:    path.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Names: []language.ExportName{{Original: "sum"}},
					Id:    path.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					All: true,
					Id:  path.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					Names: []language.ExportName{{Original: "run"}},
					Id:    path.Join(absTestFolder, "src", "lib.rs"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeRustLanguage(path.Join(testFolder, "src", "lib.rs"))
			a.NoError(err)

			lang := _lang.(*Language)

			file, err := lang.ParseFile(path.Join(absTestFolder, "src", tt.Name))
			a.NoError(err)

			_, exports, err := lang.ParseExports(context.Background(), file)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Exports)
			a.Equal(tt.Errors, exports.Errors)
		})
	}
}
