package command

import (
	"fmt"
	"github.com/axetroy/dvm/internal/deno"
)

func Current() error {
	version, err := deno.GetCurrentUseVersion()

	if err != nil {
		return err
	}

	if version == nil {
		fmt.Println("You haven't used Deno yet.")
		return nil
	}

	fmt.Println(*version)

	return nil
}
