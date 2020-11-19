package dvm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

var (
	version string
)

func init() {
	version = os.Getenv("DVM_VERSION")
}

// get current using dvm version with `v` prefix
func GetCurrentUsingVersion() string {
	return "v" + version
}

// get latest version with `v` prefix
func GetLatestRemoteVersion() (string, error) {
	res, err := http.Get("https://api.github.com/repos/axetroy/dvm/releases/latest")

	if err != nil {
		return "", errors.Wrap(err, "fetch remote version information fail")
	}

	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return "", errors.New(fmt.Sprintf("fetch remote version information and get status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", errors.Wrap(err, "read from response body fail")
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
		return "", errors.Wrap(err, "unmarshal response body fail")
	}

	version := response.TagName

	return version, nil
}
