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
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	// default command to get password
	app.Action = func(c *cli.Context) error {
		core.DisplayBySite(c.Args())
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a existing password for a website",
			Action: func(c *cli.Context) error {
				fmt.Println("Add")
				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a password for a website",
			Action: func(c *cli.Context) error {
				core.GenerateForSite(c.Args())
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all websites",
			Action: func(c *cli.Context) error {
				core.ListSites()
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
