package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const testFolder = ".root_test"

func TestRoot(t *testing.T) {
	tests := []struct {
		Name              string
		JustExpectAtLeast int
	}{
		{
			Name:              "",
			JustExpectAtLeast: 100,
		},
		{
			Name:              "help",
			JustExpectAtLeast: 100,
		},
		{
			Name: "entropy .root_test/main.py --no-browser-open",
		},
		{
			Name: ".root_test/main.py --no-browser-open",
		},
		{
			Name: "tree",
		},
		// TODO: this test now refers to an absolute path, so it's not golden testable
		// {
		//   Name: "tree random.pdf",
		// },
		// TODO: these will change once globstar entrypoints are allowed.
		// {
		//	 Name: "tree random.js",
		// },
		// {
		//	 Name: "tree random.py",
		// },
		// {
		//	 Name: "random.py",
		// },
		{
			Name: "check",
		},
		{
			Name: "check --config .root_test/.dep-tree.yml",
		},
		{
			Name: "--config .root_test/.dep-tree.yml",
		},
		{
			Name: "check --config .root_test/.dep-tree.yml-bad-path",
		},
		{
			Name: "tree .root_test/main.py --json",
		},
		{
			Name: "tree .root_test/main.py --json --config .root_test/.dep-tree.yml",
		},
		{
			Name: "tree .root_test/main.py --json --config .root_test/.dep-tree.yml",
		},
		{
			Name: "tree .root_test/main.py --json --config .root_test/.dep-tree.yml-bad-path",
		},
		{
			Name: "tree .root_test/main.py --json --config .root_test/.dep-tree.yml-bad-path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			args := strings.Split(tt.Name, " ")
			if tt.Name == "" {
				args = []string{}
			}
			root := NewRoot(args)
			b := bytes.NewBufferString("")
			root.SetOut(b)
			err := root.Execute()
			name := tt.Name + ".txt"
			name = strings.ReplaceAll(name, "/", "_")
			name = strings.ReplaceAll(name, "-", "_")
			if tt.JustExpectAtLeast > 0 {
				a := require.New(t)
				a.Greater(len(b.String()), tt.JustExpectAtLeast)
			} else {
				if err != nil {
					utils.GoldenTest(t, filepath.Join(testFolder, name), err.Error())
				} else {
					utils.GoldenTest(t, filepath.Join(testFolder, name), b.String())
				}
			}
		})
	}
}
