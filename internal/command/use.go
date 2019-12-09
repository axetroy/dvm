package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Use(version string) error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return err
	}

	var match bool

	for _, f := range files {
		v := f.Name()

		if v == version {
			match = true
			oldDenoFilepath := path.Join(core.HomeDir, ".deno/bin/deno")

			// remove it before anyway
			_ = os.Remove(oldDenoFilepath)

			if err := os.Symlink(path.Join(core.ReleaseDir, v, core.ExecutableFilename), oldDenoFilepath); err != nil {
				return err
			}
		}
	}

	if match == false {
		return errors.New(fmt.Sprintf("N/A: version `%s` is not yet installed. try install with the following command `%s` before use it", version, color.GreenString("dvm install "+version)))
	}

	return nil
}
