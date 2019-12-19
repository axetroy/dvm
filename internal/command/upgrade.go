package command

import (
	"github.com/axetroy/dvm/internal/dvm"
)

// upgrade dvm
func Upgrade(version string, force bool) error {
	return dvm.Upgrade(version, force)
}
