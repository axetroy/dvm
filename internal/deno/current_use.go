package deno

import (
	"bytes"
	"os/exec"
	"path"
	"strings"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/fs"
)

func GetCurrentUseVersion() (*string, error) {
	denoFilepath := path.Join(core.DenoDownloadDir, core.ExecutableFilename)

	if exist, err := fs.PathExists(denoFilepath); err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}

	var stdout bytes.Buffer
	cmd := exec.Command(denoFilepath, []string{"--version"}...)

	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := stdout.String()

	arr := strings.Split(strings.Split(output, "\n")[0], " ")

	version := strings.TrimSpace("v" + strings.TrimSpace(arr[1]))

	return &version, nil
}
