package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/axetroy/dvm/internal/fs"
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

// Download file from URL to the filepath
func DownloadFile(filepath string, url string) error {
	tmpl := fmt.Sprintf(`{{string . "prefix"}}{{ "%s" }} {{counters . }} {{ bar . "[" "=" ">" "-" "]"}} {{percent . }} {{speed . }}{{string . "suffix"}}`, filepath)

	// Get the data
	response, err := http.Get(url)

	if err != nil {
		return errors.Wrapf(err, "Download `%s` fail", url)
	}

	defer response.Body.Close()

	// 404
	if response.StatusCode == http.StatusNotFound {
		return errors.New(http.StatusText(response.StatusCode))
	}

	if response.StatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("download file with status code %d", response.StatusCode))
	}

	if err := fs.EnsureDir(path.Dir(filepath)); err != nil {
		return errors.Wrapf(err, "ensure `%s` fail", path.Dir(filepath))
	}

	// Create the file
	writer, err := os.Create(filepath)

	if err != nil {
		return errors.Wrapf(err, "Create `%s` fail", filepath)
	}

	defer func() {
		err = writer.Close()

		if err != nil {
			_ = os.Remove(filepath)
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
