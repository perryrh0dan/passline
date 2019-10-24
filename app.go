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
	app.EnableBashCompletion = true

	// default command to get password
	app.Action = core.DisplayByName

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
			Name:      "create",
			Aliases:   []string{"c"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<name> <username>",
			Action:    core.GenerateForSite,
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name> <username>",
			Action:    core.DeleteItem,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edit an item",
			ArgsUsage: "<name> <username>",
			Action:    core.EditItem,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    core.ListSites,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
