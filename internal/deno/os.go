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
	case "linux":
		denoOS = "linux"
		break
	case "windows":
		denoOS = "windows"
		break
	default:
		return nil, errors.New("not support your platform")
	}

	return &denoOS, nil
}
