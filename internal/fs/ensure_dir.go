package fs

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

// ensure dir exist
func EnsureDir(dir string) error {
	parent := path.Dir(dir)
	if _, err := os.Stat(parent); err != nil {
		if os.IsNotExist(err) {
			if err := EnsureDir(parent); err != nil {
				return errors.Wrapf(err, "ensure dir `%s` fail", dir)
			}
		} else {
			return errors.WithStack(err)
		}
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				return errors.Wrapf(err, "ensure dir `%s` fail", dir)
			}
		} else {
			return errors.WithStack(err)
		}
	}
	return nil
}
