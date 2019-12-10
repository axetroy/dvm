package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// copy a file or dir
func Copy(dest, src string) (err error) {

	var (
		srcFile    *os.File
		targetFile *os.File
		fileInfo   os.FileInfo
		files      []os.FileInfo
	)

	if fileInfo, err = os.Stat(src); err != nil {
		return
	}

	if fileInfo.IsDir() {
		// read dir and copy one by one
		if files, err = ioutil.ReadDir(src); err != nil {
			return
		}

		if err = EnsureDir(dest); err != nil {
			return
		}

		for _, file := range files {
			filename := file.Name()
			src = path.Join(src, filename)
			dest = path.Join(dest, filename)
			if err = Copy(dest, src); err != nil {
				return err
			}
		}

	} else {
		if srcFile, err = os.Open(src); err != nil {
			return
		}

		defer srcFile.Close()

		if targetFile, err = os.Create(dest); err != nil {
			return
		}

		defer targetFile.Close()

		_, err = io.Copy(targetFile, srcFile)
	}
	return
}
