package main

import (
	"fmt"
	"sort"

	ucli "github.com/urfave/cli"

	"github.com/perryrh0dan/passline/pkg/cli"
)

func setupApp() *ucli.App {
	app := ucli.NewApp()
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

	`, ucli.AppHelpTemplate)

	// default command to get password
	app.Action = cli.DisplayItem

	app.Commands = []ucli.Command{
		{
			Name:      "create",
			Aliases:   []string{"c"},
			Usage:     "Add an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action:    cli.CreateItem,
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name> <username>",
			Action:    cli.DeleteItem,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edit an item",
			ArgsUsage: "<name> <username>",
			Action:    cli.EditItem,
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<name> <username>",
			Action:    cli.GenerateItem,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    cli.ListItems,
		},
	}

	sort.Sort(ucli.FlagsByName(app.Flags))
	sort.Sort(ucli.CommandsByName(app.Commands))

	return app
}
