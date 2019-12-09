package command

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/axetroy/dvm/internal/core"
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

	if err := os.RemoveAll(core.RootDir); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Dir(core.DenoDownloadDir)); err != nil {
		return err
	}

	fmt.Println("Uninstall successfully! see u.")

	return nil
}
