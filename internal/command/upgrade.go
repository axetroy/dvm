package command

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/util"
	"github.com/fatih/color"
)

func Upgrade(version string, force bool) error {
	if version == "" {
		res, err := http.Get("https://api.github.com/repos/axetroy/dvm/releases/latest")

		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return err
		}

		type Asset struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}

		type Response struct {
			TagName string  `json:"tag_name"`
			Assets  []Asset `json:"assets"`
		}

		response := Response{}

		if err = json.Unmarshal(body, &response); err != nil {
			return err
		}

		version = response.TagName
	}

	// get current dvm version
	dvmExecutablePath, err := os.Executable()

	if err != nil {
		return err
	}

	cmd := exec.Command(dvmExecutablePath, "version")

	cmdOutput, err := cmd.CombinedOutput()

	if !force && version == "v"+string(cmdOutput) {
		fmt.Printf("You are using the latest version `%s`\n", color.GreenString(version))
		return nil
	}

	tarFilename := fmt.Sprintf("dvm_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)

	tarFilepath := path.Join(core.CacheDir, tarFilename)

	downloadURL := fmt.Sprintf("https://github.com/axetroy/dvm/releases/download/%s/%s", version, tarFilename)

	defer os.RemoveAll(core.CacheDir)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		_ = os.RemoveAll(core.CacheDir)
		os.Exit(1)
	}()

	if err := util.DownloadFile(tarFilepath, downloadURL); err != nil {
		return nil
	}

	// decompress the tag
	if err := decompress(tarFilepath, core.CacheDir); err != nil {
		return err
	}

	newDvmFilepath := path.Join(core.CacheDir, "dvm")

	if runtime.GOOS == "windows" {
		newDvmFilepath += ".exe"
	}

	// cover the binary file
	if err := os.Rename(newDvmFilepath, dvmExecutablePath); err != nil {
		return err
	}

	ps := exec.Command(dvmExecutablePath, "--help")

	ps.Stderr = os.Stderr
	ps.Stdout = os.Stdout

	if err := ps.Run(); err != nil {
		return err
	}

	return nil
}

// decompress tar.gz
func decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		filename := dest + hdr.Name

		file, err := os.Create(filename)

		if err != nil {
			return err
		}

		if err := file.Chmod(os.FileMode(hdr.Mode)); err != nil {
			return err
		}

		if _, err := io.Copy(file, tr); err != nil {
			return err
		}
	}
	return nil
}
