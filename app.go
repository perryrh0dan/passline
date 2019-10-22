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

	// default command to get password
	app.Action = func(c *cli.Context) error {
		_ = core.DisplayByName(c)
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Generate a password for an item",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of the item",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Global passline password",
				},
			},
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				_ = core.GenerateForSite(c)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete an item",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of the item",
				},
			},
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				_ = core.DeleteByName(c)
				return nil
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				_ = core.ListAllItems()
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
