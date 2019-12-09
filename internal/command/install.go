package command

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/axetroy/dvm/internal/fs"
	"github.com/axetroy/dvm/internal/util"
)

func Install(version string) error {
	filename, err := deno.GetRemoteTarFilename()

	if err != nil {
		return err
	}

	downloadURL, err := deno.GetRemoteDownloadURL(version)

	if err != nil {
		return err
	}

	cacheFilepath := path.Join(core.CacheDir, *filename)

	if err := util.DownloadFile(cacheFilepath, *downloadURL); err != nil {
		return err
	}

	denoFilepath, err := util.Unzip(cacheFilepath, path.Dir(cacheFilepath))

	if err != nil {
		return err
	}

	currentVersionWorkspaceDir := path.Join(core.ReleaseDir, version)

	if err := fs.EnsureDir(currentVersionWorkspaceDir); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = os.RemoveAll(currentVersionWorkspaceDir)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit
		// make sure dir been remove if it exit
		if err != nil {
			_ = os.RemoveAll(currentVersionWorkspaceDir)
		}
		os.Exit(1)
	}()

	newDenoFilepath := path.Join(currentVersionWorkspaceDir, path.Base(*denoFilepath))

	if err := os.Rename(*denoFilepath, newDenoFilepath); err != nil {
		return err
	}

	if currentUseVersion, err := deno.GetCurrentUseVersion(); err != nil {
		return err
	} else if currentUseVersion == nil {
		// if not use Deno yet. then use it.
		return Use(version)
	} else {
		// if Deno have been use. then install and do nothing.
		return nil
	}
}
