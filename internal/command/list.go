package command

import (
	"fmt"
	"io/ioutil"

	"github.com/axetroy/dvm/internal/deno"
	"github.com/axetroy/dvm/internal/dvm"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// list of installed Deno versions
func List() error {
	files, err := ioutil.ReadDir(dvm.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", dvm.ReleaseDir)
	}

	if len(files) == 0 {
		fmt.Printf("You have not installed Deno yet. try install with the following command `%s` before use it\n", color.GreenString("dvm install v0.25.0"))
		return nil
	}

	currentDenoVersion, err := deno.GetCurrentUsingVersion()

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
