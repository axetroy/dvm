package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/axetroy/dvm/internal/dvm"
	"github.com/pkg/errors"
)

// uninstall dvm
func Destroy() error {
	confirm := false

	prompt := &survey.Confirm{
		Message: "Do you want to uninstall dvm?",
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return errors.Wrap(err, "prompt fail")
	}

	if !confirm {
		return nil
	}

	//remove $HOME/.dvm
	if err := os.RemoveAll(dvm.RootDir); err != nil {
		return errors.Wrap(err, "remove `$HOME/.dvm` fail")
	}

	// remove $HOME/.deno/bin/deno
	currentUseDenoFilepath := filepath.Join(dvm.DenoBinDir, dvm.ExecutableFilename)
	if err := os.RemoveAll(currentUseDenoFilepath); err != nil {
		return errors.Wrapf(err, "remove `$HOME/.deno/bin/%s` fail", dvm.ExecutableFilename)
	}

	dvmFilepath, err := os.Executable()

	if err != nil {
		return errors.Wrap(err, "get dvm executable filepath fail")
	}

	// remove dvm executable file
	if err := os.Remove(dvmFilepath); err != nil {
		return errors.Wrap(err, "remove dvm executable filepath fail")
	}

	fmt.Println("Uninstall successfully! see u.")

	return nil
}
