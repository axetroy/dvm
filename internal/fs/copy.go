package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// copy a file or dir
func Copy(dest, src string) error {
	var (
		srcFile    *os.File
		targetFile *os.File
		fileInfo   os.FileInfo
		files      []os.FileInfo
		err        error
	)

	if fileInfo, err = os.Stat(src); err != nil {
		return errors.WithStack(err)
	}

	if fileInfo.IsDir() {
		// read dir and copy one by one
		if files, err = ioutil.ReadDir(src); err != nil {
			return errors.WithStack(err)
		}

		if err = EnsureDir(dest); err != nil {
			return errors.WithStack(err)
		}

		for _, file := range files {
			filename := file.Name()
			src = filepath.Join(src, filename)
			dest = filepath.Join(dest, filename)
			if err = Copy(dest, src); err != nil {
				return errors.WithStack(err)
			}
		}

	} else {
		if srcFile, err = os.Open(src); err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			_ = srcFile.Close()
		}()

		if targetFile, err = os.Create(dest); err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			_ = targetFile.Close()
		}()

		if _, err = io.Copy(targetFile, srcFile); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
