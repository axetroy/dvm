package command

import (
	"os"
	"path/filepath"

	"github.com/axetroy/dvm/internal/dvm"
	"github.com/pkg/errors"
)

// unused Deno, but will not remove it.
func Unused() error {
	denoFilepath := filepath.Join(dvm.DenoBinDir, dvm.ExecutableFilename)
	if err := os.RemoveAll(denoFilepath); err != nil {
		return errors.Wrapf(err, "unused Deno fail. try remove `%s` by manual", denoFilepath)
	}

	return nil
}
