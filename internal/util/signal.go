package util

import (
	"os"
	"syscall"
)

func GetAbortSignals() []os.Signal {
	return []os.Signal{
		os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	}
}
