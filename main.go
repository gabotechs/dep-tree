package main

import (
	"dep-tree/cmd"
	"os"
)

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
