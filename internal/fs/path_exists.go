package fs

import (
	"github.com/pkg/errors"
	"os"
)

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "stat file `%s` fail", path)
	}

	return true, nil
}
