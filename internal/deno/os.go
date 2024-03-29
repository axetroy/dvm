package deno

import (
	"runtime"

	semver "github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

// get Deno os for current platform
func GetDenoOS(version string) (*string, error) {
	var denoOS string

	v, err := semver.NewVersion(version)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	v1, err := semver.NewVersion("0.39.0")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch runtime.GOOS {
	case "darwin":
		if v.LessThan(v1) {
			denoOS = "osx"
		} else {
			denoOS = "apple-darwin"
		}
	case "windows":
		if v.LessThan(v1) {
			denoOS = "win"
		} else {
			denoOS = "pc-windows-msvc"
		}
	default:
		// default to linux
		if v.LessThan(v1) {
			denoOS = "linux"
		} else {
			denoOS = "unknown-linux-gnu"
		}
	}

	return &denoOS, nil
}
