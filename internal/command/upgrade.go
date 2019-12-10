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

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func getLatestVersion() (*string, error) {
	res, err := http.Get("https://api.github.com/repos/axetroy/dvm/releases/latest")

	if err != nil {
		return nil, errors.Wrap(err, "fetch remote version information fail")
	}

	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("fetch remote version information and get status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.Wrap(err, "read from response body fail")
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
		return nil, errors.Wrap(err, "unmarshal response body fail")
	}

	version := response.TagName

	return &version, nil
}

func Upgrade(version string, force bool) error {
	var (
		err         error
		tarFilename = fmt.Sprintf("dvm_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
		tarFilepath = path.Join(core.CacheDir, tarFilename)
	)

	if version == "" {
		if v, err := getLatestVersion(); err != nil {
			return errors.Wrap(err, "get latest version fail")
		} else {
			version = *v
		}
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
		return err
	}

	cmd := exec.Command(dvmExecutablePath, "version")

	cmdOutput, err := cmd.CombinedOutput()

	currentDvmVersion := "v" + string(cmdOutput)

	if !force && version == currentDvmVersion {
		fmt.Printf("You are using the latest version `%s`\n", color.GreenString(version))
		return nil
	}

	fmt.Printf("Upgrade dvm `%s` to `%s`\n", currentDvmVersion, version)

	defer os.RemoveAll(core.CacheDir)

	quit := make(chan os.Signal)
	signal.Notify(quit, util.GetAbortSignals()...)

	go func() {
		<-quit
		fmt.Printf("What made you cancel the download? you can download the file via `%s` and update manually.\n", downloadURL)
		fmt.Println("Good Luck :)")
		_ = os.RemoveAll(core.CacheDir)
		os.Exit(255)
	}()

	if err = util.DownloadFile(tarFilepath, downloadURL); err != nil {
		return errors.Wrap(err, "download fail")
	}

	// decompress the tag
	if err := decompress(tarFilepath, core.CacheDir); err != nil {
		return errors.Wrap(err, "unzip fail")
	}

	downloadedDvmFilepath := path.Join(core.CacheDir, "dvm")

	if runtime.GOOS == "windows" {
		downloadedDvmFilepath += ".exe"
	}

	// cover the binary file
	if err := os.Rename(downloadedDvmFilepath, dvmExecutablePath); err != nil {
		return errors.Wrapf(err, "rename `%s` to `%s` fail", downloadedDvmFilepath, dvmExecutablePath)
	}

	ps := exec.Command(dvmExecutablePath, "--help")

	ps.Stderr = os.Stderr
	ps.Stdout = os.Stdout

	if err := ps.Run(); err != nil {
		return errors.Wrap(err, "upgrade fail")
	}

	fmt.Printf("dvm upgrade success at `%s`\n", dvmExecutablePath)

	return nil
}

// decompress tar.gz
func decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)

	if err != nil {
		return errors.Wrapf(err, "open file `%s` fail", tarFile)
	}

	defer srcFile.Close()

	gr, err := gzip.NewReader(srcFile)

	if err != nil {
		return errors.Wrapf(err, "read zip file fail")
	}

	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "read from zip file fail")
		}

		filename := path.Join(dest, hdr.Name)

		file, err := os.Create(filename)

		if err != nil {
			return errors.Wrapf(err, "unzip and create file `%s` fail\n", filename)
		}

		if runtime.GOOS != "windows" {
			if err := file.Chmod(os.FileMode(hdr.Mode)); err != nil {
				_ = file.Close()
				return errors.Wrap(err, "change file mode fail")
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
