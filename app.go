package main

import (
	"sort"

	"github.com/perryrh0dan/passline/pkg/core"
	"github.com/urfave/cli"
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Passline"
	app.Usage = "Password manager"
	app.HelpName = "passline"
	app.Version = "0.0.1"
	app.Description = "Password manager for the command line"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password, p",
			Usage: "Password",
		},
	}

	// default command to get password
	app.Action = func(c *cli.Context) error {
		_ = core.DisplayBySite(c)
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Add an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name>",
			Action: func(c *cli.Context) error {
				_ = core.DeleteItem(c)
				return nil
			},
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<website> <username>",
			Action: func(c *cli.Context) error {
				_ = core.GenerateForSite(c)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all items",
			Action: func(c *cli.Context) error {
				_ = core.ListSites()
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
