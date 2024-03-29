package python

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const resolverTestFolder = ".resolve_test"

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
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "foo", "foo.py")},
			},
		},
		{
			Name:       "File in folder in folder at top level",
			Entrypoint: "main.py",
			Slices:     []string{"foo", "bar", "bar"},
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "foo", "bar", "bar.py")},
			},
		},
		{
			Name:       "__init__.py in folder",
			Entrypoint: "main.py",
			Slices:     []string{"foo"},
			Expected: &ResolveResult{
				InitModule: &InitModuleResult{
					Path: filepath.Join(absPath, "foo", "__init__.py"),
					PythonFiles: []string{
						filepath.Join(absPath, "foo", "__init__.py"),
						filepath.Join(absPath, "foo", "foo.py"),
					},
				},
			},
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
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "baz", "baz.py")},
			},
		},
		{
			Name:       "Import whole dir in a non-top level folder with PYTHONPATH set",
			Entrypoint: "foo/foo.py",
			Slices:     []string{"baz"},
			PythonPath: []string{absPath},
			Expected: &ResolveResult{
				Directory: &DirectoryResult{
					PythonFiles: []string{filepath.Join(absPath, "baz", "baz.py")},
					Path:        filepath.Join(absPath, "baz"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakePythonLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)
			lang.cfg.PythonPath = append(lang.cfg.PythonPath, tt.PythonPath...)
			resolved := lang.ResolveAbsolute(tt.Slices, filepath.Dir(filepath.Join(absPath, tt.Entrypoint)))
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
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "foo", "foo.py")},
			},
		},
		{
			Name:      "Import from same dir",
			Slices:    []string{"foo"},
			Dir:       filepath.Join(absPath, "foo"),
			StepsBack: 0,
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "foo", "foo.py")},
			},
		},
		{
			Name:      "Import from nested dir with init",
			Slices:    []string{"foo"},
			Dir:       absPath,
			StepsBack: 0,
			Expected: &ResolveResult{
				InitModule: &InitModuleResult{
					Path: filepath.Join(absPath, "foo", "__init__.py"),
					PythonFiles: []string{
						filepath.Join(absPath, "foo", "__init__.py"),
						filepath.Join(absPath, "foo", "foo.py"),
					},
				},
			},
		},
		{
			Name:      "Import from parent dir",
			Slices:    []string{"foo"},
			Dir:       filepath.Join(absPath, "foo", "bar"),
			StepsBack: 1,
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "foo", "foo.py")},
			},
		},
		{
			Name:      "Import from double parent dir",
			Slices:    []string{"baz", "baz"},
			Dir:       filepath.Join(absPath, "foo", "bar"),
			StepsBack: 2,
			Expected: &ResolveResult{
				File: &FileResult{Path: filepath.Join(absPath, "baz", "baz.py")},
			},
		},
		{
			Name:      "Import from parent dir to __init__.py",
			Slices:    []string{},
			Dir:       filepath.Join(absPath, "foo", "bar"),
			StepsBack: 1,
			Expected: &ResolveResult{
				InitModule: &InitModuleResult{
					Path: filepath.Join(absPath, "foo", "__init__.py"),
					PythonFiles: []string{
						filepath.Join(absPath, "foo", "__init__.py"),
						filepath.Join(absPath, "foo", "foo.py"),
					},
				},
			},
		},
		{
			Name:   "Import baz directory from root",
			Slices: []string{"baz"},
			Dir:    absPath,
			Expected: &ResolveResult{
				Directory: &DirectoryResult{
					PythonFiles: []string{filepath.Join(absPath, "baz", "baz.py")},
					Path:        filepath.Join(absPath, "baz"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			resolved, err := ResolveRelative(tt.Slices, tt.Dir, tt.StepsBack)
			a.NoError(err)
			a.Equal(tt.Expected, resolved)
		})
	}
}
