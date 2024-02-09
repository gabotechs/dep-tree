package js

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLanguage_Display(t *testing.T) {
	tests := []struct {
		Name            string
		Path            string
		ExpectedRelPath string
		ExpectedPackage string
	}{
		{
			Name:            "with a parent package.json",
			Path:            filepath.Join(resolverTestFolder, "src", "utils", "sum.ts"),
			ExpectedPackage: "test-project",
			ExpectedRelPath: "src/utils/sum.ts",
		},
		{
			Name:            "with a parent package.json (same as above for checking cache)",
			Path:            filepath.Join(resolverTestFolder, "src", "utils", "sum.ts"),
			ExpectedPackage: "test-project",
			ExpectedRelPath: "src/utils/sum.ts",
		},
		{
			Name:            "with two parent package.json, one without name",
			Path:            filepath.Join(resolverTestFolder, "src", "module", "main.ts"),
			ExpectedPackage: "test-project",
			ExpectedRelPath: "src/module/main.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeJsLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)
			absPath, _ := filepath.Abs(tt.Path)
			file, err := lang.ParseFile(absPath)
			a.NoError(err)
			a.Equal(tt.ExpectedPackage, file.Package)
			a.Equal(tt.ExpectedRelPath, file.RelPath)
		})
	}
}
