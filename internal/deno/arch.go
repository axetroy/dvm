package deno

import (
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

// get Deno arch for current platform
func GetDenoArch(version string) (*string, error) {
	var denoArch string

	v, err := semver.NewVersion(version)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	v1, err := semver.NewVersion("0.39.0")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch runtime.GOARCH {
	case "amd64":
		fallthrough
	case "arm64":
		if v.LessThan(v1) {
			denoArch = "x64"
		} else {
			denoArch = "x86_64"
		}
	default:
		return nil, errors.New("not support your platform")
	}

	return &denoArch, nil
}
