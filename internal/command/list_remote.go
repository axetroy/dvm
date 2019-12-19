package command

import (
	"fmt"

	"github.com/axetroy/dvm/internal/deno"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func ListRemote() error {
	versions, err := deno.GetRemoteVersions()

	if err != nil {
		return errors.Wrap(err, "get remote version fail")
	}

	currentDenoVersion, err := deno.GetCurrentUseVersion()

	if err != nil {
		// ignore error
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
			fmt.Printf(" (%s)", color.GreenString("Currently using"))
		} else if isLatest {
			fmt.Printf(" (%s)", color.GreenString("Latest"))
		}

		fmt.Println()
	}

	return nil
}
