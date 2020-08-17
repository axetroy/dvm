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
		return nil, errors.WithStack(err)
	}

	defer r.Close()

	if len(r.File) > 1 {
		return nil, errors.WithStack(err)
	}

	f := r.File[0]

	newFilepath := path.Join(dest, f.Name)

	src, err := f.Open()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = src.Close()
	}()

	dst, err := os.OpenFile(newFilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = dst.Close()
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, errors.WithStack(err)
	}

	return &newFilepath, nil
}

func decompressGz(tarFile, dest string) (*string, error) {
	fileReader, err := os.Open(tarFile)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = fileReader.Close()
	}()

	gzipReader, err := gzip.NewReader(fileReader)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = gzipReader.Close()
	}()

	newFilepath := path.Join(dest, "deno")

	fileWriter, err := os.Create(newFilepath)

	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
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
