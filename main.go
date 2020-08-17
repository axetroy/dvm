package main

import (
	"fmt"
	"os"

	"github.com/axetroy/dvm/internal/command"
	"github.com/axetroy/dvm/internal/core"
	"github.com/axetroy/dvm/internal/dvm"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "dvm"
	app.Usage = "version manager for Deno"
	app.Version = dvm.GetCurrentUsingVersion()
	app.Authors = []*cli.Author{
		{
			Name:  "Axetroy",
			Email: "axetroy.dev@gmail.com",
		},
	}

	cli.AppHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}
USAGE:
    {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
    {{if len .Authors}}
AUTHOR:
    {{range .Authors}}{{ . }}{{end}}
    {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}} {{ .ArgsUsage }}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
    {{.Copyright}}
    {{end}}{{if .Version}}
VERSION:
    {{.Version}}
    {{end}}
EXAMPLES:
    {{.Name}} install v1.0.0
    {{.Name}} use v1.0.0
    {{.Name}} uninstall v1.0.0
    {{.Name}} exec v1.0.0 https://deno.land/std/examples/welcome.ts
    {{.Name}} ls
    {{.Name}} ls-remote

SOURCE CODE:
    https://github.com/axetroy/dvm
`

	app.Commands = []*cli.Command{
		{
			Name:  "current",
			Usage: "Display currently activated version of Deno",
			Action: func(c *cli.Context) error {
				return command.Current()
			},
		},
		{
			Name:  "ls",
			Usage: "List installed versions",
			Action: func(c *cli.Context) error {
				return command.List()
			},
		},
		{
			Name:  "ls-remote",
			Usage: "List remote versions available for install",
			Action: func(c *cli.Context) error {
				return command.ListRemote()
			},
		},
		{
			Name:      "install",
			Usage:     "Download and install specified Deno version",
			ArgsUsage: "<version>",
			Aliases:   []string{"i"},
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return errors.New(fmt.Sprintf("require argument <%s>", "version"))
				}
				return command.Install(c.Args().First())
			},
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall specified Deno version",
			ArgsUsage: "<version1> <version2> ...",
			Aliases:   []string{"rm"},
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return errors.New(fmt.Sprintf("require argument <%s>", "version"))
				}
				return command.Uninstall(c.Args().Slice())
			},
		},
		{
			Name:      "use",
			Usage:     "Use specified Deno version",
			ArgsUsage: "<version>",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return errors.New(fmt.Sprintf("require argument <%s>", "version"))
				}
				return command.Use(c.Args().First())
			},
		},
		{
			Name:  "unused",
			Usage: "Unused Deno",
			Action: func(c *cli.Context) error {
				return command.Unused()
			},
		},
		{
			Name:      "exec",
			Usage:     "Run Deno command on <version>.",
			ArgsUsage: "<version> <args>",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return errors.New(fmt.Sprintf("require argument <%s>", "version"))
				}

				version := c.Args().First()
				args := c.Args().Slice()[1:]

				return command.Exec(version, args)
			},
		},
		// here is the commands for dvm self
		{
			Name:      "upgrade",
			Usage:     "Upgrade dvm",
			ArgsUsage: "[version]",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "Force upgrade", Value: false},
			},
			Action: func(c *cli.Context) error {
				return command.Upgrade(c.Args().First(), c.Bool("force"))
			},
		},
		{
			Name:  "destroy",
			Usage: "Uninstall dvm",
			Action: func(c *cli.Context) error {
				return command.Destroy()
			},
		},
	}

	// regardless of the result, the cache directory should be delete
	if err := app.Run(os.Args); err != nil {
		if os.Getenv("DEBUG") != "" {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err.Error())
			fmt.Printf("run with environment variables %s to print more information\n", color.GreenString("DEBUG=1"))
		}
		_ = os.RemoveAll(core.CacheDir)
		os.Exit(1)
	} else {
		_ = os.RemoveAll(core.CacheDir)
	}
}
