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

	"github.com/Masterminds/semver/v3"
	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// upgrade dvm
func Upgrade(version string, force bool) error {
	var (
		err         error
		tarFilename = fmt.Sprintf("dvm_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
		tarFilepath = filepath.Join(core.CacheDir, tarFilename)
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
		_ = os.RemoveAll(core.CacheDir)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, util.GetAbortSignals()...)

	go func() {
		<-quit
		fmt.Printf("What made you cancel the download? you can download the file via `%s` and update manually.\n", downloadURL)
		fmt.Println("Good Luck :)")
		_ = os.RemoveAll(core.CacheDir)
		os.Exit(255)
	}()

	fmt.Printf("Download %s\n", downloadURL)

	if err = util.DownloadFile(tarFilepath, downloadURL); err != nil {
		return errors.WithStack(err)
	}

	defer signal.Stop(quit)

	// decompress the tag
	if err := decompress(tarFilepath, core.CacheDir); err != nil {
		return errors.WithStack(err)
	}

	downloadedDvmFilepath := filepath.Join(core.CacheDir, "dvm")

	if runtime.GOOS == "windows" && !strings.HasSuffix(downloadedDvmFilepath, ".exe") {
		// Ensure to add '.exe' to given path on Windows
		downloadedDvmFilepath += ".exe"
	}

	if err := util.ReplaceExecutableFile(downloadedDvmFilepath, dvmExecutablePath); err != nil {
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
