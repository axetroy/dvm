package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axetroy/dvm/internal/deno"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func ListRemote() error {
	url := "https://api.github.com/repos/denoland/deno/git/refs/tags"

	r, err := http.Get(url)

	if err != nil {
		return errors.Wrapf(err, "request `%s` fail", url)
	}

	defer r.Body.Close()

	if r.StatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("download file with status code %d", r.StatusCode))
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
		return errors.Wrap(err, "real response body error")
	}

	var res []node
	var versions = make([]string, 0)

	if err := json.Unmarshal(b, &res); err != nil {
		return errors.Wrap(err, "JSON parse fail")
	}

	for _, n := range res {
		versions = append(versions, strings.TrimLeft(n.Ref, "refs/tags/"))
	}

	currentDenoVersion, err := deno.GetCurrentUseVersion()

	if err != nil {
		// TODO: ignore error, remote this line in the future
		// return err
	}

	for i, v := range versions {
		isLatest := i == len(versions)-1
		isCurrentUse := currentDenoVersion != nil && *currentDenoVersion == v

		if isCurrentUse {
			fmt.Print(color.GreenString(v))
		} else {
			fmt.Print(v)
		}

		if isCurrentUse && isLatest {
			fmt.Printf(" (%s)", color.GreenString("Latest and currently using"))
		} else if isCurrentUse {
			fmt.Printf(" (%s)", color.GreenString("currently using"))
		} else if isLatest {
			fmt.Printf(" (%s)", color.GreenString("Latest"))
		}

		fmt.Println()
	}

	return nil
}
