package cmd

import (
	"bytes"
	"path"
	"strings"
	"testing"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const testFolder = ".root_test"

func TestRoot(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			Name: "help",
		},
		{
			Name: "render",
		},
		{
			Name: "render random.js",
		},
		{
			Name: "render random.pdf",
		},
		{
			Name: "render random.py",
		},
		{
			Name: "random.py",
		},
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
			Name: ".root_test/main.py --json",
		},
		{
			Name: "render .root_test/main.py --json --config .root_test/.dep-tree.yml",
		},
		{
			Name: ".root_test/main.py --json --config .root_test/.dep-tree.yml",
		},
		{
			Name: "render .root_test/main.py --json --config .root_test/.dep-tree.yml-bad-path",
		},
		{
			Name: ".root_test/main.py --json --config .root_test/.dep-tree.yml-bad-path",
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
			name := strings.ReplaceAll(tt.Name+".txt", "/", "|")
			if err != nil {
				utils.GoldenTest(t, path.Join(testFolder, name), err.Error())
			} else {
				utils.GoldenTest(t, path.Join(testFolder, name), b.String())
			}
		})
	}
}
