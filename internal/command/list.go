package command

import (
	"fmt"
	"io/ioutil"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func List() error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	if len(files) == 0 {
		fmt.Printf("You have not installed Deno yet. try install with the following command `%s` before use it\n", color.GreenString("dvm install v0.25.0"))
		return nil
	}

	currentDenoVersion, err := deno.GetCurrentUseVersion()

	if err != nil {
		// ignore error
	}

	for _, f := range files {
		fmt.Print(f.Name())

		if currentDenoVersion != nil && *currentDenoVersion == f.Name() {
			fmt.Printf(" (%s)", color.GreenString("currently using"))
		}

		fmt.Println()
	}

	return nil
}
