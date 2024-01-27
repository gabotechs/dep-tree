package rust

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const cargoTomlTestFolder = ".cargo_toml_test"

func TestCargoToml(t *testing.T) {
	absPath, _ := filepath.Abs(cargoTomlTestFolder)

	tests := []struct {
		Name     string
		Expected CargoToml
		Error    string
	}{
		{
			Name: "normal",
			Expected: CargoToml{
				path: filepath.Join(absPath, "normal"),
				Dependencies: map[string]localDependency{
					"bar": {"../bar"},
					"baz": {"../baz"},
					"foo": {"../foo"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			result, err := readCargoToml(filepath.Join(cargoTomlTestFolder, tt.Name))
			if tt.Error != "" {
				a.ErrorContains(err, tt.Error)
			} else {
				a.NoError(err)
				a.Equal(tt.Expected, *result)
			}
		})
	}
}
