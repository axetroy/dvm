package deno

import (
	"runtime"

	"github.com/pkg/errors"
)

func GetDenoOS() (*string, error) {
	var denoOS string

	switch runtime.GOOS {
	case "darwin":
		denoOS = "osx"
		break
	case "openbsd":
		fallthrough
	case "freebsd":
		fallthrough
	case "linux":
		denoOS = "linux"
		break
	case "windows":
		denoOS = "win"
		break
	default:
		return nil, errors.New("not support your platform")
	}

	return &denoOS, nil
}
