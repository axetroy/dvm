package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/fs"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// use Deno
func Use(version string) error {
	if version == "" {
		data, _ := ioutil.ReadFile(".dvmrc")

		if data == nil {
			return errors.New(fmt.Sprintf("No .dvmrc file found \n\nrequire argument <%s>", "version"))
		} else {
			version = setCharV(strings.TrimSpace(string(data)))
		}
	}

	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	var match bool

	for _, f := range files {
		v := f.Name()

		if v == version {
			match = true
			oldDenoFilepath := filepath.Join(core.DenoBinDir, core.ExecutableFilename)

			// remove it before anyway
			if err := os.Remove(oldDenoFilepath); err != nil {
				if !os.IsNotExist(err) {
					return errors.Wrapf(err, "remove `%s` fail", oldDenoFilepath)
				}
			}

			p := filepath.Join(core.ReleaseDir, v, core.ExecutableFilename)

			if err := os.Symlink(p, oldDenoFilepath); err != nil {
				// Windows requires permission for soft link
				// Use copy as fallback
				if runtime.GOOS == "windows" {
					if err = fs.Copy(oldDenoFilepath, p); err != nil {
						return errors.Wrapf(err, "copy `%s` to `%s` fail", p, oldDenoFilepath)
					}
				} else {
					return errors.Wrapf(err, "use `%s` fail", version)
				}
			}
		}
	}

	if !match {
		return errors.New(fmt.Sprintf("N/A: version `%s` is not yet installed. try install with the following command `%s` before use it", version, color.GreenString("dvm install "+version)))
	}

	fmt.Printf("Currently using Deno %s\n", version)

	return nil
}

func setCharV(version string) string {
	if string(version[0]) != "v" {
		return "v" + string(version)
	}

	return string(version)
}
