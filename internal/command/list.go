package command

import (
	"fmt"
	"io/ioutil"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/fatih/color"
)

func List() error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return err
	}

	currentDenoVersion, err := deno.GetCurrentUseVersion()

	if err != nil {
		return err
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
