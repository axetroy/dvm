package command

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
)

func Uninstall(version string) error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Name() == version {
			if err := os.RemoveAll(path.Join(core.ReleaseDir, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}
