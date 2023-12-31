package js

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const workspacesTestDir = ".workspaces_test"

func TestNewWorkspaces(t *testing.T) {
	a := require.New(t)

	for _, entry := range []string{
		filepath.Join(workspacesTestDir, "nested", "c"),
		filepath.Join(workspacesTestDir, "a"),
		workspacesTestDir,
	} {
		result, err := NewWorkspaces(entry)
		a.NoError(err)
		a.NotNil(result)
		abs, _ := filepath.Abs(workspacesTestDir)

		a.Equal(map[string]WorkspaceEntry{
			"a":      {absPath: filepath.Join(abs, "a")},
			"c":      {absPath: filepath.Join(abs, "nested", "c")},
			"f":      {absPath: filepath.Join(abs, "r-nested", "1", "f")},
			"g":      {absPath: filepath.Join(abs, "r-nested", "2", "3", "g")},
			"h":      {absPath: filepath.Join(abs, "r-nested", "h")},
			"@foo/k": {absPath: filepath.Join(abs, "foo", "k")},
		}, result.ws)
	}
}

func TestNewWorkspaces_parses_packages(t *testing.T) {
	a := require.New(t)
	result, err := NewWorkspaces(filepath.Join(workspacesTestDir, "other"))
	a.NoError(err)
	abs, _ := filepath.Abs(workspacesTestDir)
	a.Equal(map[string]WorkspaceEntry{
		"foo": {absPath: filepath.Join(abs, "other", "foo")},
	}, result.ws)
}

func TestWorkspaces_ResolveFromWorkspaces(t *testing.T) {
	abs, _ := filepath.Abs(workspacesTestDir)

	tests := []struct {
		Name       string
		Unresolved string
		Resolved   string
		Error      string
	}{
		{
			Name:       "Basic",
			Unresolved: "a/src/a",
			Resolved:   filepath.Join(abs, "a", "src", "a.js"),
		},
		{
			Name:       `From npm org`,
			Unresolved: "@foo/k/src/k",
			Resolved:   filepath.Join(abs, "foo", "k", "src", "k.ts"),
		},
		{
			Name:       "Nested",
			Unresolved: "g/g",
			Resolved:   filepath.Join(abs, "r-nested", "2", "3", "g", "g.tsx"),
		},
		{
			Name:       "Index in source",
			Unresolved: "c",
			Resolved:   filepath.Join(abs, "nested", "c", "src", "index.ts"),
		},
		{
			Name:       "Index in root",
			Unresolved: "f",
			Resolved:   filepath.Join(abs, "r-nested", "1", "f", "index.jsx"),
		},
		{
			Name:       "No index",
			Unresolved: "h",
			Error:      "has no index file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			ws, err := NewWorkspaces(workspacesTestDir)
			a.NoError(err)
			a.NotNil(ws)
			result, err := ws.ResolveFromWorkspaces(tt.Unresolved)
			if tt.Error != "" {
				a.ErrorContains(err, tt.Error)
			} else {
				a.NoError(err)
				a.Equal(tt.Resolved, result)
			}
		})
	}
}
