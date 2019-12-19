package command

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/pkg/errors"
)

// uninstall Deno
func Uninstall(version string) error {
	currentUseVersion, err := deno.GetCurrentUsingVersion()

	if err != nil {
		return err
	}

	if currentUseVersion != nil && *currentUseVersion == version {
		if err := os.Remove(path.Join(core.DenoBinDir, core.ExecutableFilename)); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	for _, f := range files {
		if f.Name() == version {
			if err := os.RemoveAll(path.Join(core.ReleaseDir, f.Name())); err != nil {
				return errors.Wrapf(err, "uninstall deno@`%s` fail", version)
			}
		}
	}

	return nil
}
