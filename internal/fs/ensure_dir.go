package fs

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

func EnsureDir(dir string) error {
	parent := path.Dir(dir)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if err := EnsureDir(parent); err != nil {
			return errors.Wrapf(err, "ensure dir `%s` fail", dir)
		}
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return errors.Wrapf(err, "ensure dir `%s` fail", dir)
		}
	}
	return nil
}
