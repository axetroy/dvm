package core

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/axetroy/dvm/internal/fs"
)

var (
	HomeDir            string // $HOME
	RootDir            string // $HOME/.dvm
	CacheDir           string // cache dir
	ReleaseDir         string // $HOME/.dvm/releases
	ExecutableFilename string // in Unix: deno. in Windows: deno.exe
	DenoBinDir         string // deno bin dir. defaults to $HOME/.deno/bin
)

func init() {
	var err error

	defer func() {
		if err != nil {
			log.Panicln(err)
		}
	}()

	ExecutableFilename = "deno"

	if runtime.GOOS == "windows" {
		ExecutableFilename += ".exe"
	}

	if h, e := os.UserHomeDir(); e != nil {
		err = e
		return
	} else {
		HomeDir = h
	}

	DenoBinDir = path.Join(HomeDir, ".deno", "bin")
	RootDir = path.Join(HomeDir, ".dvm")
	ReleaseDir = path.Join(RootDir, "releases")

	if e := fs.EnsureDir(DenoBinDir); e != nil {
		err = e
		return
	}

	if e := fs.EnsureDir(ReleaseDir); e != nil {
		err = e
		return
	}

	if c, e := os.UserCacheDir(); e != nil {
		err = e
		return
	} else {
		CacheDir = path.Join(c, "dvm")

		if e := fs.EnsureDir(CacheDir); e != nil {
			err = e
			return
		}
	}
}
