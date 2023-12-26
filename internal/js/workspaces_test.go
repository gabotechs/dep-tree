package js

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const workspacesTestDir = ".workspaces_test"

func TestNewWorkspaces(t *testing.T) {
	a := require.New(t)

	for _, entry := range []string{
		path.Join(workspacesTestDir, "nested", "c"),
		path.Join(workspacesTestDir, "a"),
		workspacesTestDir,
	} {
		result, err := NewWorkspaces(entry)
		a.NoError(err)
		a.NotNil(result)
		abs, _ := filepath.Abs(workspacesTestDir)

		a.Equal(map[string]WorkspaceEntry{
			"a":      {absPath: path.Join(abs, "a")},
			"c":      {absPath: path.Join(abs, "nested", "c")},
			"f":      {absPath: path.Join(abs, "r-nested", "1", "f")},
			"g":      {absPath: path.Join(abs, "r-nested", "2", "3", "g")},
			"h":      {absPath: path.Join(abs, "r-nested", "h")},
			"@foo/k": {absPath: path.Join(abs, "foo", "k")},
		}, result.ws)
	}
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
			Resolved:   path.Join(abs, "a", "src", "a.js"),
		},
		{
			Name:       `From npm org`,
			Unresolved: "@foo/k/src/k",
			Resolved:   path.Join(abs, "foo", "k", "src", "k.ts"),
		},
		{
			Name:       "Nested",
			Unresolved: "g/g",
			Resolved:   path.Join(abs, "r-nested", "2", "3", "g", "g.tsx"),
		},
		{
			Name:       "Index in source",
			Unresolved: "c",
			Resolved:   path.Join(abs, "nested", "c", "src", "index.ts"),
		},
		{
			Name:       "Index in root",
			Unresolved: "f",
			Resolved:   path.Join(abs, "r-nested", "1", "f", "index.jsx"),
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
