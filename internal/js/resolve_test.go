package js

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const resolverTestFolder = ".resolve_test"

func TestParser_ResolvePath(t *testing.T) {
	absPath, _ := filepath.Abs(resolverTestFolder)

	tests := []struct {
		Name          string
		Unresolved    string
		Cwd           string
		Resolved      string
		ExpectedError string
	}{
		{
			Name:       "from relative",
			Cwd:        filepath.Join(resolverTestFolder, "src", "utils"),
			Unresolved: "../foo",
			Resolved:   filepath.Join(absPath, "src", "foo.ts"),
		},
		{
			Name:       "from baseUrl",
			Cwd:        filepath.Join(resolverTestFolder, "src"),
			Unresolved: "foo",
			Resolved:   filepath.Join(absPath, "src", "foo.ts"),
		},
		{
			Name:       "from paths override",
			Cwd:        filepath.Join(resolverTestFolder, "src"),
			Unresolved: "@utils/sum",
			Resolved:   filepath.Join(absPath, "src", "utils", "sum.ts"),
		},
		{
			Name:       "from paths override with glob pattern",
			Cwd:        filepath.Join(resolverTestFolder, "src"),
			Unresolved: "@/helpers/diff",
			Resolved:   filepath.Join(absPath, "src", "helpers", "diff.ts"),
		},
		{
			Name:       "from package.json in folder",
			Cwd:        filepath.Join(resolverTestFolder, "src"),
			Unresolved: "./module",
			Resolved:   filepath.Join(absPath, "src", "module", "main.ts"),
		},
		{
			Name:          "Does not resolve invalid relative import",
			Cwd:           filepath.Join(resolverTestFolder, "src", "utils"),
			Unresolved:    "./foo",
			ExpectedError: "could not perform relative import for './foo' because the file or dir was not found",
		},
		{
			Name:       "Does not resolve invalid import",
			Cwd:        resolverTestFolder,
			Unresolved: "react",
		},
		{
			Name:       "Does not resolve invalid relative import",
			Cwd:        resolverTestFolder,
			Unresolved: filepath.Join("src", "utils", "foo"),
		},
		{
			Name:          "Does not resolve invalid path override import",
			Cwd:           filepath.Join(resolverTestFolder, "src"),
			Unresolved:    "@/helpers/bar",
			ExpectedError: "import '@/helpers/bar' was matched to path '.resolve_test/src/helpers/bar' in tscofing's paths option, but the resolved path did not match an existing file",
		},
		{
			Name:          "Empty name does not panic",
			Unresolved:    "",
			ExpectedError: "import path cannot be empty",
		},
		{
			Name:          "One dot path does not panic",
			Unresolved:    ".",
			ExpectedError: "invalid import path .",
		},
		{
			Name:          "Two dot path works",
			Unresolved:    "..",
			Cwd:           resolverTestFolder,
			ExpectedError: "could not perform relative import for '..'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeJsLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)
			resolved, err := lang.ResolvePath(tt.Unresolved, tt.Cwd)
			if tt.ExpectedError != "" {
				a.ErrorContains(err, tt.ExpectedError)
			} else {
				a.NoError(err)
				a.Equal(tt.Resolved, resolved)
			}
		})
	}
}
