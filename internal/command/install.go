package command

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/axetroy/dvm/internal/fs"
	"github.com/axetroy/dvm/internal/util"
	"github.com/pkg/errors"
)

// install Deno
func Install(version string) error {
	v, filename, downloadURL, err := deno.GetRemoteDownloadURL(version)

	if err != nil {
		return errors.WithStack(err)
	}

	version = v

	cacheFilepath := filepath.Join(core.CacheDir, filename)

	quitAndCleanCache := make(chan os.Signal)
	signal.Notify(quitAndCleanCache, util.GetAbortSignals()...)

	go func() {
		<-quitAndCleanCache
		// make sure dir been remove if it exit
		_ = os.RemoveAll(core.CacheDir)
		os.Exit(255)
	}()

	defer signal.Stop(quitAndCleanCache)

	fmt.Printf("Download %s\n", downloadURL)

	if err := util.DownloadFile(cacheFilepath, downloadURL); err != nil {
		if err.Error() == http.StatusText(http.StatusNotFound) {
			return errors.Wrapf(err, "Deno %s is not yet released or available", version)
		}
		return errors.Wrapf(err, "download remote file `%s` fail", downloadURL)
	}

	denoFilepath, err := deno.Unzip(cacheFilepath, filepath.Dir(cacheFilepath))

	if err != nil {
		return errors.Wrap(err, "unzip tar file fail")
	}

	currentVersionWorkspaceDir := filepath.Join(core.ReleaseDir, version)

	if err := fs.EnsureDir(currentVersionWorkspaceDir); err != nil {
		return errors.Wrap(err, "ensure workspace fail")
	}

	quitAndCleanWorkspace := make(chan os.Signal)
	signal.Notify(quitAndCleanWorkspace, util.GetAbortSignals()...)

	go func() {
		<-quitAndCleanWorkspace
		// make sure dir been remove if it exit
		if err != nil {
			_ = os.RemoveAll(currentVersionWorkspaceDir)
		}
		os.Exit(255)
	}()

	defer signal.Stop(quitAndCleanWorkspace)

	defer func() {
		if err != nil {
			_ = os.RemoveAll(currentVersionWorkspaceDir)
		}
	}()

	newDenoFilepath := filepath.Join(currentVersionWorkspaceDir, filepath.Base(*denoFilepath))

	if err := os.Rename(*denoFilepath, newDenoFilepath); err != nil {
		return errors.Wrapf(err, "rename Deno executable file fail. try move `%s` to `%s` by manual", *denoFilepath, newDenoFilepath)
	}

	fmt.Printf("Install successfully! try `dvm use %s` before you use it.\n", version)

	if currentUseVersion, err := deno.GetCurrentUsingVersion(); err != nil {
		// ignore error
		return nil
	} else if currentUseVersion == nil {
		// if not use Deno yet. then use it.
		return Use(version)
	} else {
		// if Deno have been use. then install and do nothing.
		return nil
	}
}
