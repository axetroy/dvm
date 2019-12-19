package deno

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

func decompressZip(tarFile, dest string) (*string, error) {
	r, err := zip.OpenReader(tarFile)

	if err != nil {
		return nil, errors.Wrapf(err, "read zip file `%s` fail", tarFile)
	}

	defer r.Close()

	if len(r.File) > 1 {
		return nil, errors.New("window .zip file should only contain single file")
	}

	f := r.File[0]

	newFilepath := path.Join(dest, f.Name)

	src, err := f.Open()

	if err != nil {
		return nil, err
	}

	defer src.Close()

	dst, err := os.Create(newFilepath)

	if err != nil {
		return nil, err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	return &newFilepath, nil
}

func decompressGz(tarFile, dest string) (*string, error) {
	fileReader, err := os.Open(tarFile)

	if err != nil {
		return nil, errors.Wrapf(err, "open file `%s` fail", tarFile)
	}

	defer fileReader.Close()

	gzipReader, err := gzip.NewReader(fileReader)

	if err != nil {
		return nil, errors.Wrapf(err, "gzip decode fail")
	}

	defer gzipReader.Close()

	newFilepath := path.Join(dest, "deno")

	fileWriter, err := os.Create(newFilepath)

	if err != nil {
		return nil, errors.Wrapf(err, "create file `%s` fail", newFilepath)
	}

	defer func() {
		err = fileWriter.Close()

		if err != nil {
			err = os.Remove(newFilepath)
		}
	}()

	if _, err = io.Copy(fileWriter, gzipReader); err != nil {
		return nil, err
	}

	if err := fileWriter.Chmod(os.FileMode(0755)); err != nil {
		return nil, errors.Wrap(err, "change file mod fail")
	}

	return &newFilepath, nil
}

func Unzip(tarFilepath, destDir string) (*string, error) {
	switch path.Ext(tarFilepath) {
	case ".zip":
		return decompressZip(tarFilepath, destDir)
	case ".gz":
		return decompressGz(tarFilepath, destDir)
	default:
		return nil, errors.New(fmt.Sprintf("not support unzip the file `%s`", tarFilepath))
	}
}
