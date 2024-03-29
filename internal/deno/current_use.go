package deno

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func getInstalledDenoFilepath() string {
	if denoPath, err := exec.LookPath("deno"); err != nil {
		return ""
	} else {
		return strings.TrimSpace(denoPath)
	}
}

// get current using version
func GetCurrentUsingVersion() (*string, error) {
	denoFilepath := getInstalledDenoFilepath()

	if denoFilepath == "" {
		return nil, nil
	}

	args := []string{"--version"}
	cmd := exec.Command(denoFilepath, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, errors.Wrapf(err, "`deno --version` failed\n%s", string(output))
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return nil, errors.New(string(output))
	}

	arr := strings.Split(strings.Split(string(output), "\n")[0], " ")

	version := strings.TrimSpace("v" + strings.TrimSpace(arr[1]))

	return &version, nil
}
