package main

import (
	"fmt"
	"os"

	"github.com/axetroy/dvm/internal/command"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "dvm"
	app.Usage = "version management for Deno"
	app.Version = "0.1.2"
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
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
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
			Usage:     "Download and install a <version> from source.",
			ArgsUsage: "<version>",
			Action: func(c *cli.Context) error {
				return command.Install(c.Args().First())
			},
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall specified Deno version",
			ArgsUsage: "<version>",
			Action: func(c *cli.Context) error {
				return command.Uninstall(c.Args().First())
			},
		},
		{
			Name:      "use",
			Usage:     "Use specified Deno version",
			ArgsUsage: "<version>",
			Action: func(c *cli.Context) error {
				return command.Use(c.Args().First())
			},
		},
		{
			Name:  "unuse",
			Usage: "Unuse specified Deno version",
			Action: func(c *cli.Context) error {
				return command.Unuse()
			},
		},
		// here is the commands for dvm self
		{
			Name:  "version",
			Usage: "Print dvm version info to stdout",
			Action: func(context *cli.Context) error {
				_, err := os.Stdout.Write([]byte(app.Version))

				if err != nil {
					return errors.Wrap(err, "write to stdout fail")
				}

				return nil
			},
		},
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
			Usage: "Uninstall dvm and remove all the thing about Deno",
			Action: func(c *cli.Context) error {
				return command.Destroy()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
