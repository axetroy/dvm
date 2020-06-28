package deno

import (
	"runtime"

	"github.com/Masterminds/semver/v3"
)

// get Deno os for current platform
func GetDenoOS(version string) (*string, error) {
	var denoOS string

	v, err := semver.NewVersion(version)

	if err != nil {
		return nil, err
	}

	v1, err := semver.NewVersion("0.39.0")

	if err != nil {
		return nil, err
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
