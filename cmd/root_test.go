package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
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
		{
			Name: "tree random.pdf",
		},
		{
			Name: "tree random.js",
		},
		{
			Name: "tree random.py",
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
			Name: "tree .root_test/main.py --json",
		},
		{
			Name: "tree .root_test/main.py --json --exclude .root_test/dep.py",
		},
		{
			Name: "tree .root_test/main.py --json --exclude .root_test/*.py",
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
			name = strings.ReplaceAll(name, "*", "_")
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

func TestInferLang(t *testing.T) {
	tests := []struct {
		Name     string
		Files    []string
		Expected language.Language
		Error    string
	}{
		{
			Name:  "zero files",
			Files: []string{},
			Error: "at least 1 file must be provided for infering the language",
		},
		{
			Name:     "only 1 file",
			Files:    []string{"foo.js"},
			Expected: &js.Language{},
		},
		{
			Name:     "majority of files",
			Files:    []string{"foo.js", "bar.rs", "foo.rs", "foo.py"},
			Expected: &rust.Language{},
		},
		{
			Name:     "unrelated files",
			Files:    []string{"foo.py", "foo.pdf"},
			Expected: &python.Language{},
		},
		{
			Name:  "no match",
			Files: []string{"foo.pdf", "bar.docx"},
			Error: "none of the provided files belong to the a supported language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := inferLang(tt.Files, &config.Config{})
			if tt.Error != "" {
				a.ErrorContains(err, tt.Error)
			} else {
				a.NoError(err)
				a.IsType(tt.Expected, lang)
			}
		})
	}
}
