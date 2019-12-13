package main

import (
	"context"
	"fmt"
	"sort"

	ucli "github.com/urfave/cli/v2"

	"github.com/perryrh0dan/passline/pkg/cli"
)

func setupApp(ctx context.Context) *ucli.App {
	cli.Init(ctx)

	app := ucli.NewApp()
	app.Name = "Passline"
	app.Usage = "Password manager"
	app.HelpName = "passline"
	app.Version = "0.2.0"
	app.Description = "Password manager for the command line"
	app.EnableBashCompletion = true

	// Append website information to default helper print
	app.CustomAppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: 
   https://github.com/perryrh0dan/passline

	`, ucli.AppHelpTemplate)

	// default command to get password
	app.Action = func(c *ucli.Context) error { return cli.DisplayItem(ctx, c) }

	app.Commands = []*ucli.Command{
		{
			Name:      "backup",
			Aliases:   []string{"b"},
			Usage:     "Create a backup",
			ArgsUsage: "<path>",
			Action:    func(c *ucli.Context) error { return cli.CreateBackup(ctx, c) },
		},
		{
			Name:      "create",
			Aliases:   []string{"c"},
			Usage:     "Add an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action:    func(c *ucli.Context) error { return cli.CreateItem(ctx, c) },
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name> <username>",
			Action:    func(c *ucli.Context) error { return cli.DeleteItem(ctx, c) },
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edit an item",
			ArgsUsage: "<name> <username>",
			Action:    func(c *ucli.Context) error { return cli.EditItem(ctx, c) },
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<name> <username>",
			Flags: []ucli.Flag{
				&ucli.StringFlag{
					Name:    "mode",
					Aliases: []string{"m"},
					Value:   "default",
					Usage:   "Change between default and advanced mode",
				},
			},
			Action: func(c *ucli.Context) error { return cli.GenerateItem(ctx, c) },
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    func(c *ucli.Context) error { return cli.ListItems(ctx, c) },
		},
		{
			Name:      "restore",
			Aliases:   []string{"r"},
			Usage:     "Restore backup",
			ArgsUsage: "<path>",
			Action:    func(c *ucli.Context) error { return cli.RestoreBackup(ctx, c) },
		},
	}

	sort.Sort(ucli.FlagsByName(app.Flags))
	sort.Sort(ucli.CommandsByName(app.Commands))

	return app
}
