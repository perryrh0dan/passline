package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/perryrh0dan/passline/core"
	"github.com/urfave/cli"
)

func main() {
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
			Usage:     "Add a existing password for a website",
			ArgsUsage: "<website> <username> <password>",
			Action: func(c *cli.Context) error {
				fmt.Println("Add")
				return nil
			},
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for a website",
			ArgsUsage: "<website> <username>",
			Action: func(c *cli.Context) error {
				_ = core.GenerateForSite(c)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all websites",
			Action: func(c *cli.Context) error {
				_ = core.ListSites()
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
