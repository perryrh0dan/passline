package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	// default command to get password
	app.Action = func(c *cli.Context) error {
		displayBySite(c.Args())
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a password for a website",
			Action: func(c *cli.Context) error {
				password := generate()
				fmt.Println(password)
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a existing password for a website",
			Action: func(c *cli.Context) error {
				fmt.Println("Add")
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
