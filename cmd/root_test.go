package cmd

import (
	"bytes"
	"path"
	"strings"
	"testing"

	"dep-tree/internal/utils"
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
			Name: "render random.rs",
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
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			root := NewRoot()
			b := bytes.NewBufferString("")
			root.SetOut(b)
			root.SetArgs(strings.Split(tt.Name, " "))
			err := root.Execute()
			if err != nil {
				utils.GoldenTest(t, path.Join(testFolder, tt.Name+".txt"), err.Error())
			} else {
				utils.GoldenTest(t, path.Join(testFolder, tt.Name+".txt"), b.String())
			}
		})
	}
}
