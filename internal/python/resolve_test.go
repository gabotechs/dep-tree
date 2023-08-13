package python

import (
	"github.com/stretchr/testify/require"

	"path"
	"path/filepath"
	"testing"
)

const resolverTestFolder = ".resolve_test"

// TODO: what happens if you import a folder with no __init__.py
//  ANSWER: that is valid, each file gets imported as a module.

func TestResolveAbsolute(t *testing.T) {
	absPath, _ := filepath.Abs(resolverTestFolder)

	tests := []struct {
		Name       string
		Entrypoint string
		Slices     []string
		PythonPath []string
		Expected   *ResolveResult
	}{
		{
			Name:       "File in folder at top level",
			Entrypoint: "main.py",
			Slices:     []string{"foo", "foo"},
			Expected:   &ResolveResult{File: "foo/foo.py"},
		},
		{
			Name:       "File in folder in folder at top level",
			Entrypoint: "main.py",
			Slices:     []string{"foo", "bar", "bar"},
			Expected:   &ResolveResult{File: "foo/bar/bar.py"},
		},
		{
			Name:       "__init__.py in folder",
			Entrypoint: "main.py",
			Slices:     []string{"foo"},
			Expected:   &ResolveResult{InitModule: "foo/__init__.py"},
		},
		{
			Name:       "Import in a non-top level folder",
			Entrypoint: "foo/foo.py",
			Slices:     []string{"baz", "baz"},
			Expected:   nil,
		},
		{
			Name:       "Import in a non-top level folder with PYTHONPATH set",
			Entrypoint: "foo/foo.py",
			Slices:     []string{"baz", "baz"},
			PythonPath: []string{absPath},
			Expected:   &ResolveResult{File: "baz/baz.py"},
		},
		{
			Name:       "Import whole dir in a non-top level folder with PYTHONPATH set",
			Entrypoint: "foo/foo.py",
			Slices:     []string{"baz"},
			PythonPath: []string{absPath},
			Expected:   &ResolveResult{Directory: "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			_lang, err := MakePythonLanguage(path.Join(absPath, tt.Entrypoint))
			a.NoError(err)
			lang := _lang.(*Language)
			lang.PythonPath = append(lang.PythonPath, tt.PythonPath...)

			resolved, err := lang.ResolveAbsolute(tt.Slices)
			a.NoError(err)
			switch {
			case tt.Expected == nil:
				// nothing.
			case tt.Expected.File != "":
				tt.Expected.File = path.Join(absPath, tt.Expected.File)
			case tt.Expected.Directory != "":
				tt.Expected.Directory = path.Join(absPath, tt.Expected.Directory)
			case tt.Expected.InitModule != "":
				tt.Expected.InitModule = path.Join(absPath, tt.Expected.InitModule)
			}

			a.Equal(tt.Expected, resolved)
		})
	}
}

func TestResolveRelative(t *testing.T) {
	absPath, _ := filepath.Abs(resolverTestFolder)

	tests := []struct {
		Name      string
		Slices    []string
		Dir       string
		StepsBack int
		Expected  *ResolveResult
	}{
		{
			Name:      "Import from nested dir",
			Slices:    []string{"foo", "foo"},
			Dir:       absPath,
			StepsBack: 0,
			Expected:  &ResolveResult{File: "foo/foo.py"},
		},
		{
			Name:      "Import from same dir",
			Slices:    []string{"foo"},
			Dir:       path.Join(absPath, "foo"),
			StepsBack: 0,
			Expected:  &ResolveResult{File: "foo/foo.py"},
		},
		{
			Name:      "Import from nested dir with init",
			Slices:    []string{"foo"},
			Dir:       absPath,
			StepsBack: 0,
			Expected:  &ResolveResult{InitModule: "foo/__init__.py"},
		},
		{
			Name:      "Import from parent dir",
			Slices:    []string{"foo"},
			Dir:       path.Join(absPath, "foo", "bar"),
			StepsBack: 1,
			Expected:  &ResolveResult{File: "foo/foo.py"},
		},
		{
			Name:      "Import from double parent dir",
			Slices:    []string{"baz", "baz"},
			Dir:       path.Join(absPath, "foo", "bar"),
			StepsBack: 2,
			Expected:  &ResolveResult{File: "baz/baz.py"},
		},
		{
			Name:      "Import from parent dir to __init__.py",
			Slices:    []string{},
			Dir:       path.Join(absPath, "foo", "bar"),
			StepsBack: 1,
			Expected:  &ResolveResult{InitModule: "foo/__init__.py"},
		},
		{
			Name:     "Import baz directory from root",
			Slices:   []string{"baz"},
			Dir:      absPath,
			Expected: &ResolveResult{Directory: "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			resolved, err := ResolveRelative(tt.Slices, tt.Dir, tt.StepsBack)
			a.NoError(err)
			switch {
			case tt.Expected == nil:
				// nothing.
			case tt.Expected.File != "":
				tt.Expected.File = path.Join(absPath, tt.Expected.File)
			case tt.Expected.Directory != "":
				tt.Expected.Directory = path.Join(absPath, tt.Expected.Directory)
			case tt.Expected.InitModule != "":
				tt.Expected.InitModule = path.Join(absPath, tt.Expected.InitModule)
			}

			a.Equal(tt.Expected, resolved)
		})
	}
}
