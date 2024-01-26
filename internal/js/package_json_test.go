package js

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const packageJsonTestDir = ".package_json_test"

func TestPackageJson(t *testing.T) {
	absPath, _ := filepath.Abs(packageJsonTestDir)

	tests := []struct {
		Name     string
		Expected packageJson
		Error    string
	}{
		{
			Name: "package.json",
			Expected: packageJson{
				absPath: absPath,
				Name:    "test",
				Main:    "foo.js",
				Workspaces: []interface{}{
					"foo",
				},
			},
		},
		{
			Name: "folder",
			Expected: packageJson{
				absPath: filepath.Join(absPath, "folder"),
				Name:    "",
				Main:    "bar.js",
			},
		},
		{
			Name:  "non existent",
			Error: "no such file or directory",
		},
		{
			Name:  "badly_formed",
			Error: "cannot unmarshal number",
		},
		{
			Name: "empty",
			Expected: packageJson{
				absPath: filepath.Join(absPath, "empty"),
				Name:    "",
				Main:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			pckJson, err := readPackageJson(filepath.Join(packageJsonTestDir, tt.Name))
			if tt.Error != "" {
				a.ErrorContains(err, tt.Error)
			} else {
				a.NoError(err)
				a.Equal(tt.Expected, *pckJson)
			}
		})
	}
}
