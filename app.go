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
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add an existing password for a website",
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
					Usage: "Password",
				},
			},
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
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
			ArgsUsage: "<name>",
			Action: func(c *cli.Context) error {
				_ = core.DeleteItem(c)
				return nil
			},
		},
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
			},
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				_ = core.GenerateForSite(c)
				return nil
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"l"},
			Usage:     "List all items",
			ArgsUsage: " ",
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
