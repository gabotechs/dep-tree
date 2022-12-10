package main

import (
	"os"

	"dep-tree/cmd"
)

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
