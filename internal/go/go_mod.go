package golang

import (
	"os"

	"github.com/gabotechs/dep-tree/internal/utils"
	"golang.org/x/mod/modfile"
)

type GoMod struct {
	Module string
}

func _ParseGoMod(file string) (*GoMod, error) {
	modBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	goMod, err := modfile.Parse(file, modBytes, nil)
	if err != nil {
		return nil, err
	}
	return &GoMod{
		Module: goMod.Module.Mod.Path,
	}, nil
}

var ParseGoMod = utils.Cached1In1OutErr(_ParseGoMod)
