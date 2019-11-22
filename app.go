package main

import (
	"fmt"
	"sort"

	"github.com/perryrh0dan/passline/pkg/core"
	"github.com/urfave/cli"
)

func setupApp() *cli.App {
	app := cli.NewApp()
	pl := core.NewPassline()
	app.Name = "Passline"
	app.Usage = "Password manager"
	app.HelpName = "passline"
	app.Version = "0.0.1"
	app.Description = "Password manager for the command line"
	app.EnableBashCompletion = true

	// Append website information to default helper print
	app.CustomAppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: 
   https://github.com/perryrh0dan/passline

	`, cli.AppHelpTemplate)

	// default command to get password
	app.Action = pl.DisplayItem

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Add an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action:    pl.AddItem,
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name> <username>",
			Action:    pl.DeleteItem,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edit an item",
			ArgsUsage: "<name> <username>",
			Action:    pl.EditItem,
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<name> <username>",
			Action:    pl.GenerateItem,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    pl.ListSites,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
