package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/axetroy/dvm/internal/core"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// execute with specified Deno version
func Exec(version string, args []string) error {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	var denoPath = ""

	for _, f := range files {
		if f.Name() == version {
			denoPath = path.Join(core.ReleaseDir, version, core.ExecutableFilename)
		}
	}

	if denoPath == "" {
		fmt.Printf("You have not installed Deno %s yet. try install with the following command `%s` before use it\n", version, color.GreenString(fmt.Sprintf("dvm install %s", version)))
		return nil
	}

	ps := exec.Command(denoPath, args...)

	ps.Stdin = os.Stdin
	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr

	if err := ps.Run(); err != nil {
		return err
	}

	return nil
}
