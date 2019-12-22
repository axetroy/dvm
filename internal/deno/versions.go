package deno

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// get latest remote version
func GetLatestRemoteVersion() (string, error) {
	versions, err := GetRemoteVersions()

	if err != nil {
		return "", errors.Wrap(err, "get remote versions fail")
	}

	if len(versions) == 0 {
		return "", errors.New("no distributions found")
	}

	latest := versions[len(versions)-1]

	return latest, nil
}

// get remote versions
func GetRemoteVersions() ([]string, error) {
	url := "https://api.github.com/repos/denoland/deno/git/refs/tags"

	r, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrapf(err, "request `%s` fail", url)
	}

	defer r.Body.Close()

	if r.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("download file with status code %d", r.StatusCode))
	}

	type node struct {
		Ref    string `json:"ref"`
		NodeID string `json:"node_id"`
		URL    string `json:"url"`
		Object struct {
			SHA  string `json:"sha"`
			Type string `json:"type"`
			URL  string `json:"url"`
		}
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, errors.Wrap(err, "real response body error")
	}

	var res []node
	var versions = make([]string, 0)

	if err := json.Unmarshal(b, &res); err != nil {
		return nil, errors.Wrap(err, "JSON parse fail")
	}

	for _, n := range res {
		versions = append(versions, strings.TrimSpace(strings.TrimLeft(n.Ref, "refs/tags/")))
	}

	return versions, nil
}
