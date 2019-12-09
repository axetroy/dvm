package command

import (
	"os"
	"path"

	"github.com/axetroy/dvm/internal/core"
)

func Unuse() error {
	return os.RemoveAll(path.Join(core.DenoDownloadDir, "deno"))
}
