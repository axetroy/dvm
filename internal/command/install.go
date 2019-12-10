package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/deno"
	"github.com/axetroy/dvm/internal/fs"
	"github.com/axetroy/dvm/internal/util"
	"github.com/pkg/errors"
)

func isHaveInstall(version string) (bool, error) {
	files, err := ioutil.ReadDir(core.ReleaseDir)

	if err != nil {
		return false, errors.Wrapf(err, "read dir `%s` fail", core.ReleaseDir)
	}

	for _, f := range files {
		if version == f.Name() && f.IsDir() {
			denoFilepath := path.Join(core.ReleaseDir, f.Name(), core.ExecutableFilename)
			if isExist, err := fs.PathExists(denoFilepath); err != nil {
				return false, errors.Wrapf(err, "detect path `%s` exist fail", denoFilepath)
			} else if isExist {
				return true, nil
			}
		}
	}

	return false, nil
}

func Install(version string) error {
	filename, err := deno.GetRemoteTarFilename()

	if err != nil {
		return errors.Wrap(err, "get remote tar filename fail")
	}

	downloadURL, err := deno.GetRemoteDownloadURL(version)

	if err != nil {
		return errors.Wrap(err, "get remote download url fail")
	}

	cacheFilepath := path.Join(core.CacheDir, *filename)

	quitAndCleanCache := make(chan os.Signal)
	signal.Notify(quitAndCleanCache, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quitAndCleanCache
		// make sure dir been remove if it exit
		if err != nil {
			_ = os.RemoveAll(core.CacheDir)
		}
		os.Exit(255)
	}()

	defer signal.Stop(quitAndCleanCache)

	if err := util.DownloadFile(cacheFilepath, *downloadURL); err != nil {
		return errors.Wrap(err, "download remote file fail")
	}

	denoFilepath, err := util.Unzip(cacheFilepath, path.Dir(cacheFilepath))

	if err != nil {
		return errors.Wrap(err, "unzip tar file fail")
	}

	currentVersionWorkspaceDir := path.Join(core.ReleaseDir, version)

	if err := fs.EnsureDir(currentVersionWorkspaceDir); err != nil {
		return errors.Wrap(err, "ensure workspace fail")
	}

	quitAndCleanWorkspace := make(chan os.Signal)
	signal.Notify(quitAndCleanWorkspace, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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

	newDenoFilepath := path.Join(currentVersionWorkspaceDir, path.Base(*denoFilepath))

	if err := os.Rename(*denoFilepath, newDenoFilepath); err != nil {
		return errors.Wrapf(err, "rename Deno executable file fail. try move `%s` to `%s` by manual", *denoFilepath, newDenoFilepath)
	}

	fmt.Printf("Install successfully! try `dvm use %s` before you use it.\n", version)

	if currentUseVersion, err := deno.GetCurrentUseVersion(); err != nil {
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
