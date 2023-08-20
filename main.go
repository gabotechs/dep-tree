package main

import (
	"os"

	"dep-tree/cmd"
)

func main() {
	err := cmd.NewRoot().Execute()
	if err != nil {
		os.Exit(1)
	}
}
