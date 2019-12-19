package util

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/axetroy/dvm/internal/core"
	"github.com/pkg/errors"
)

// upgrade command with new filepath
func UpgradeCommand(newFilepath, oldFilepath string) (err error) {
	var (
		old os.FileInfo
	)

	old, err = os.Stat(oldFilepath)

	if err != nil {
		return err
	}

	// cover the binary file
	if runtime.GOOS == "windows" {
		oldFilepathBackup := path.Join(core.CacheDir, old.Name()) + fmt.Sprintf(".%d.old", time.Now().UnixNano())
		if err = os.Rename(oldFilepath, oldFilepathBackup); err != nil {
			return errors.Wrapf(err, "backup old version fail")
		}

		defer func() {
			// if upgrade fail
			if err != nil {
				// try rollback
				if err = os.Rename(oldFilepathBackup, oldFilepath); err != nil {
					err = errors.Wrap(err, "rollback fail")
					return
				}
			} else {
				_ = os.Remove(oldFilepathBackup)
				return
			}
		}()

		// rename downloaded dvm to exist dvm filepath
		if err = os.Rename(newFilepath, oldFilepath); err != nil {
			return errors.Wrapf(err, "rename downloaded file to dvm filepath fail")
		}
	} else {
		if err = os.Rename(newFilepath, oldFilepath); err != nil {
			err = errors.Wrapf(err, "rename `%s` to `%s` fail", newFilepath, oldFilepath)
		}
	}

	return
}
