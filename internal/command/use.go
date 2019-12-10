package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/fs"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Use(version string) error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	var match bool

	for _, f := range files {
		v := f.Name()

		if v == version {
			match = true
			oldDenoFilepath := path.Join(core.DenoBinDir, core.ExecutableFilename)

			// remove it before anyway
			if err := os.Remove(oldDenoFilepath); !os.IsNotExist(err) {
				return errors.Wrapf(err, "remove `%s` fail", oldDenoFilepath)
			}

			p := path.Join(core.ReleaseDir, v, core.ExecutableFilename)

			if err := os.Symlink(p, oldDenoFilepath); err != nil {
				// Windows requires permission for soft link
				// Use copy as fallback
				if runtime.GOOS == "windows" {
					err = nil
					if err := fs.Copy(oldDenoFilepath, p); err != nil {
						return errors.Wrapf(err, "copy `%s` to `%s` fail", p, oldDenoFilepath)
					}
				} else {
					return errors.Wrapf(err, "use `%s` fail", version)
				}
			}
		}
	}

	if match == false {
		return errors.New(fmt.Sprintf("N/A: version `%s` is not yet installed. try install with the following command `%s` before use it", version, color.GreenString("dvm install "+version)))
	} else {
		fmt.Printf("Currently using Deno %s\n", version)
	}

	return nil
}
