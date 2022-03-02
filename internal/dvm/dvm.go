package dvm

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/axetroy/dvm/internal/fs"
	"github.com/fatih/color"
	"github.com/pkg/errors"
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

	DenoBinDir = filepath.Join(HomeDir, ".deno", "bin")
	RootDir = filepath.Join(HomeDir, ".dvm")
	ReleaseDir = filepath.Join(RootDir, "releases")

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
		CacheDir = filepath.Join(c, "dvm")

		if e := fs.EnsureDir(CacheDir); e != nil {
			err = e
			return
		}
	}
}

func CheckEnv() error {
	paths := strings.Split(os.Getenv("PATH"), ":")

	for _, p := range paths {
		if strings.ToLower(p) == strings.ToLower(DenoBinDir) {
			return nil
		}
	}

	envStr := strings.ReplaceAll(DenoBinDir, HomeDir, "$HOME")

	msg := fmt.Sprintf(`For using dvm, you need to put '%s' into $PATH environment.`, envStr)

	fmt.Fprintf(os.Stderr, "can not found '%s' in '$PATH'\n", envStr)

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		fishPath, _ := exec.LookPath("fish")
		bashPath, _ := exec.LookPath("bash")
		zshPath, _ := exec.LookPath("zsh")

		if fishPath != "" {
			cmd := color.GreenString(fmt.Sprintf(`echo 'set PATH $PATH "%s"' >> %s`, envStr, "$HOME/.config/fish/config.fish"))
			_, _ = fmt.Fprintf(os.Stderr, "If you are using fish, try to run the following command %s\n", cmd)
		}

		if zshPath != "" {
			cmd := color.GreenString(fmt.Sprintf(`echo 'export PATH=$PATH:%s' >> %s`, envStr, "$HOME/.bash_profile"))
			_, _ = fmt.Fprintf(os.Stderr, "If you are using zsh, try to run the following command %s\n", cmd)
		}

		if bashPath != "" {
			cmd := color.GreenString(fmt.Sprintf(`echo 'export PATH=$PATH:%s' >> %s`, envStr, "$HOME/.zshrc"))
			_, _ = fmt.Fprintf(os.Stderr, "If you are using zsh, try to run the following command %s\n", cmd)
		}
	}

	return errors.New(msg)
}
