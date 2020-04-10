package deno

import (
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
)

// get remote Deno tar filename
func GetRemoteTarFilename(version string) (*string, error) {
	os, err := GetDenoOS(version)

	if err != nil {
		return nil, err
	}

	arch, err := GetDenoArch(version)

	if err != nil {
		return nil, err
	}

	extensionName := "gz"

	if runtime.GOOS == "windows" {
		extensionName = "zip"
	}

	v, err := semver.NewVersion(version)

	if err != nil {
		return nil, err
	}

	v1, err := semver.NewVersion("0.39.0")

	if err != nil {
		return nil, err
	}

	var filename string

	if v.LessThan(v1) {
		filename = fmt.Sprintf("deno_%s_%s.%s", *os, *arch, extensionName)
	} else {
		// use the new release file
		filename = fmt.Sprintf("deno-%s-%s.zip", *arch, *os)
	}

	return &filename, nil
}

// get remote Deno tar download URL for specified version
func GetRemoteDownloadURL(version string) (v string, url string, err error) {

	if version == "latest" {
		if latest, err := GetLatestRemoteVersion(); err != nil {
			return "", "", err
		} else {
			version = latest
		}
	}

	filename, err := GetRemoteTarFilename(version)

	if err != nil {
		return "", "", err
	}

	downloadURL := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/%s", version, *filename)

	return version, downloadURL, nil
}
