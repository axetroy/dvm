package command

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
	"github.com/pkg/errors"
)

func Uninstall(version string) error {
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
