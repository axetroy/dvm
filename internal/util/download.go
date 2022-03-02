package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/axetroy/dvm/internal/fs"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

// Download file from URL to the filepath
func DownloadFile(filepathStr string, url string) error {
	tmpl := fmt.Sprintf(`{{string . "prefix"}}{{ "%s" }} {{counters . }} {{ bar . "[" "=" ">" "-" "]"}} {{percent . }} {{speed . }}{{string . "suffix"}}`, filepathStr)

	// Get the data
	response, err := http.Get(url)

	if err != nil {
		return errors.Wrapf(err, "Download `%s` fail", url)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	// 404
	if response.StatusCode == http.StatusNotFound {
		return errors.New(http.StatusText(response.StatusCode))
	}

	if response.StatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("download file with status code %d", response.StatusCode))
	}

	if err := fs.EnsureDir(filepath.Dir(filepathStr)); err != nil {
		return errors.Wrapf(err, "ensure `%s` fail", filepath.Dir(filepathStr))
	}

	// Create the file
	writer, err := os.Create(filepathStr)

	if err != nil {
		return errors.Wrapf(err, "Create `%s` fail", filepathStr)
	}

	defer func() {
		err = writer.Close()

		if err != nil {
			_ = os.Remove(filepathStr)
		}
	}()

	bar := pb.ProgressBarTemplate(tmpl).Start64(response.ContentLength)

	bar.SetWriter(os.Stdout)

	barReader := bar.NewProxyReader(response.Body)

	_, err = io.Copy(writer, barReader)

	bar.Finish()

	if err != nil {
		err = errors.Wrap(err, "copy fail")
	}

	return err
}
