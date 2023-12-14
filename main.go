package main

import (
	"os"

	"github.com/gabotechs/dep-tree/cmd"
)

func main() {
	err := cmd.NewRoot(nil).Execute()
	if err != nil {
		os.Exit(1)
	}
}
