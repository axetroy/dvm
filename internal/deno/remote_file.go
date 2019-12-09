package deno

import (
	"fmt"
	"runtime"
)

func GetRemoteTarFilename() (*string, error) {
	os, err := GetDenoOS()

	if err != nil {
		return nil, err
	}

	arch, err := GetDenoArch()

	if err != nil {
		return nil, err
	}

	extensionName := "gz"

	if runtime.GOOS == "windows" {
		extensionName = "zip"
	}

	filename := fmt.Sprintf("deno_%s_%s.%s", *os, *arch, extensionName)

	return &filename, nil
}

func GetRemoteDownloadURL(version string) (*string, error) {
	filename, err := GetRemoteTarFilename()

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/%s", version, *filename)

	return &url, nil
}
