package command

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/axetroy/dvm/internal/core"
	"github.com/pkg/errors"
)

func Destroy() error {
	confirm := false

	prompt := &survey.Confirm{
		Message: "Do you want to uninstall dvm?",
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}

	if confirm == false {
		return nil
	}

	//remove $HOME/.dvm
	if err := os.RemoveAll(core.RootDir); err != nil {
		return errors.Wrap(err, "remove `$HOME/.dvm` fail")
	}

	// remove $HOME/.deno
	if err := os.RemoveAll(path.Dir(core.DenoBinDir)); err != nil {
		return errors.Wrap(err, "remove `$HOME/.deno` fail")
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
