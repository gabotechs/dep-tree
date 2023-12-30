package rust

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeRustLanguage_Errors(t *testing.T) {
	cwd, _ := os.Getwd()

	tests := []struct {
		Name       string
		Entrypoint string
		Expected   string
	}{
		{
			Name:       "invalid entrypoint",
			Entrypoint: cwd,
			Expected:   "could not find Cargo.toml in any parent directory",
		},
		{
			Name:       "empty project",
			Entrypoint: path.Join(cwd, ".empty_project"),
			Expected:   "could not find any of the possible entrypoint paths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, err := MakeRustLanguage(tt.Entrypoint, nil)
			a.ErrorContains(err, tt.Expected)
		})
	}
}
