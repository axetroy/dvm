package dvm

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	semver "github.com/Masterminds/semver/v3"
	"github.com/axetroy/dvm/internal/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// upgrade dvm
func Upgrade(version string, force bool) error {
	var (
		err         error
		tarFilename = fmt.Sprintf("dvm_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
		tarFilepath = filepath.Join(CacheDir, tarFilename)
	)

	if version == "" {
		if v, err := GetLatestRemoteVersion(); err != nil {
			return errors.WithStack(err)
		} else {
			version = v
		}
	}

	_, err = semver.NewVersion(version)

	if err != nil {
		return errors.WithStack(err)
	}

	downloadURL := fmt.Sprintf("https://github.com/axetroy/dvm/releases/download/%s/%s", version, tarFilename)

	defer func() {
		if err != nil {
			fmt.Printf("If the upgrade fails, download from the `%s` and upgrade manually.\n", downloadURL)
		}
	}()

	// get current dvm version
	dvmExecutablePath, err := os.Executable()

	if err != nil {
		return errors.WithStack(err)
	}

	currentDvmVersion := GetCurrentUsingVersion()

	if !force && version == GetCurrentUsingVersion() {
		fmt.Printf("You are using the latest version `%s`\n", color.GreenString(version))
		return nil
	}

	fmt.Printf("Upgrade dvm `%s` to `%s`\n", currentDvmVersion, version)

	defer func() {
		_ = os.RemoveAll(CacheDir)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, util.GetAbortSignals()...)

	go func() {
		<-quit
		fmt.Printf("What made you cancel the download? you can download the file via `%s` and update manually.\n", downloadURL)
		fmt.Println("Good Luck :)")
		_ = os.RemoveAll(CacheDir)
		os.Exit(255)
	}()

	fmt.Printf("Download %s\n", downloadURL)

	if err = util.DownloadFile(tarFilepath, downloadURL); err != nil {
		return errors.WithStack(err)
	}

	defer signal.Stop(quit)

	// decompress the tag
	if err := decompress(tarFilepath, CacheDir); err != nil {
		return errors.WithStack(err)
	}

	downloadedDvmFilepath := filepath.Join(CacheDir, "dvm")

	if runtime.GOOS == "windows" && !strings.HasSuffix(downloadedDvmFilepath, ".exe") {
		// Ensure to add '.exe' to given path on Windows
		downloadedDvmFilepath += ".exe"
	}

	if err := replaceExecutableFile(downloadedDvmFilepath, dvmExecutablePath); err != nil {
		return errors.WithStack(err)
	}

	ps := exec.Command(dvmExecutablePath, "--help")

	ps.Stderr = os.Stderr
	ps.Stdout = os.Stdout

	if err := ps.Run(); err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("dvm upgrade success at `%s`\n", dvmExecutablePath)

	return nil
}

// decompress tar.gz
func decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)

	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		_ = srcFile.Close()
	}()

	gr, err := gzip.NewReader(srcFile)

	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		_ = gr.Close()
	}()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.WithStack(err)
		}

		filename := filepath.Join(dest, hdr.Name)

		file, err := os.Create(filename)

		if err != nil {
			return errors.WithStack(err)
		}

		if runtime.GOOS != "windows" {
			if err := file.Chmod(os.FileMode(hdr.Mode)); err != nil {
				_ = file.Close()
				return errors.WithStack(err)
			}
		}

		if _, err := io.Copy(file, tr); err != nil {
			_ = file.Close()
			return errors.Wrap(err, "copy file from zip fail")
		}

		_ = file.Close()
	}

	return nil
}

// replace executable file
// In Windows, if this executable file is running, it cannot be replaced
func replaceExecutableFile(newFilepath, oldFilepath string) (err error) {
	var (
		old os.FileInfo
	)

	old, err = os.Stat(oldFilepath)

	if err != nil {
		return err
	}

	// cover the binary file
	if runtime.GOOS == "windows" {
		oldFilepathBackup := filepath.Join(CacheDir, old.Name()) + fmt.Sprintf(".%d.old", time.Now().UnixNano())
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
