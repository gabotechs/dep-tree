package rust

import (
	"context"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testFolder = ".sample_project"

func (m *ModTree) Debug(indent int) string {
	msg := strings.Repeat(" ", indent)
	abs, _ := filepath.Abs(testFolder)
	rel, _ := filepath.Rel(abs, m.Path)
	msg += m.Name + " " + rel
	msg += "\n"
	keys := make([]string, len(m.Children))
	i := 0
	for key := range m.Children {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		msg += m.Children[key].Debug(indent + 2)
	}
	return msg
}

func TestMakeModTreeIsCached(t *testing.T) {
	a := require.New(t)
	absPath, err := filepath.Abs(path.Join(testFolder, "src", "lib.rs"))
	a.NoError(err)

	ctx := context.Background()

	start := time.Now().UnixMicro()

	ctx, _, err = MakeModTree(ctx, absPath, "crate", nil)
	a.NoError(err)

	first := time.Now().UnixMicro() - start
	start = time.Now().UnixMicro()

	absPath, err = filepath.Abs(path.Join(testFolder, "src", "sum.rs"))
	a.NoError(err)

	_, _, err = MakeModTree(ctx, absPath, "crate", nil)
	a.NoError(err)

	second := time.Now().UnixMicro() - start

	a.Greater(first, second*10)
}

func TestMakeModTree(t *testing.T) {
	a := require.New(t)
	absPath, err := filepath.Abs(path.Join(testFolder, "src", "lib.rs"))
	a.NoError(err)

	_, modTree, err := MakeModTree(context.Background(), absPath, "crate", nil)
	a.NoError(err)

	result := modTree.Debug(0)
	a.Equal(`crate src/lib.rs
  abs src/abs.rs
    abs src/abs/abs.rs
  avg src/avg.rs
    random src/avg/random.rs
  avg_2 src/avg_2.rs
    avg src/avg_2.rs
  div src/div/mod.rs
    div src/div/div.rs
    div_2 src/div/div_2/mod.rs
      div_2 src/div/div_2/div_2.rs
  sum src/sum.rs
  tests src/lib.rs
`, result)

	base := path.Dir(path.Dir(absPath))

	tests := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "abs",
			Expected: "src/abs.rs",
		},
		{
			Name:     "avg_2 avg",
			Expected: "src/avg_2.rs",
		},
		{
			Name:     "div div_2",
			Expected: "src/div/div_2/mod.rs",
		},
		{
			Name:     "div div_2 div_2",
			Expected: "src/div/div_2/div_2.rs",
		},
		{
			Name:     "div div_2 super div_2",
			Expected: "src/div/div_2/mod.rs",
		},
		{
			Name:     "div div_2 self div_2",
			Expected: "src/div/div_2/div_2.rs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			node := modTree.Search(strings.Split(tt.Name, " "))
			a.NotNil(node)
			a.Equal(path.Join(base, tt.Expected), node.Path)
		})
	}
}

func TestModTree_Errors(t *testing.T) {
	tests := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{
			Name:     "invalid path",
			Path:     path.Join(testFolder, "src", "_bad.rs"),
			Expected: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			absPath, err := filepath.Abs(tt.Path)
			a.NoError(err)

			_, _, err = MakeModTree(context.Background(), absPath, "crate", nil)
			a.ErrorContains(err, tt.Expected)
		})
	}
}
