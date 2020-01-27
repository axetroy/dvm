package command

import (
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
	"github.com/pkg/errors"
)

// unused Deno, but will not remove it.
func Unused() error {
	denoFilepath := path.Join(core.DenoBinDir, core.ExecutableFilename)
	if err := os.RemoveAll(denoFilepath); err != nil {
		return errors.Wrapf(err, "unused Deno fail. try remove `%s` by manual", denoFilepath)
	}

	return nil
}
