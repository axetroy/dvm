package command

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/axetroy/dvm/internal/deno"
	"github.com/axetroy/dvm/internal/dvm"
	"github.com/pkg/errors"
)

// uninstall Deno
func Uninstall(versions []string) error {
	currentUseVersion, err := deno.GetCurrentUsingVersion()

	if err != nil {
		return err
	}

	for _, version := range versions {
		if currentUseVersion != nil && *currentUseVersion == version {
			if err := os.Remove(filepath.Join(dvm.DenoBinDir, dvm.ExecutableFilename)); err != nil {
				return err
			}
		}

		files, err := ioutil.ReadDir(dvm.ReleaseDir)

		if err != nil {
			return errors.Wrapf(err, "read dir `%s` fail", dvm.ReleaseDir)
		}

		for _, f := range files {
			if f.Name() == version {
				if err := os.RemoveAll(filepath.Join(dvm.ReleaseDir, f.Name())); err != nil {
					return errors.Wrapf(err, "uninstall deno@`%s` fail", version)
				}
			}
		}
	}

	return nil
}
