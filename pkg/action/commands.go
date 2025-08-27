package action

import (
	ucli "github.com/urfave/cli/v3"
)

// GetCommands returns the ucli commands exported by this module
func (s *Action) GetCommands() []*ucli.Command {
	return []*ucli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Adds an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action:    s.Add,
		},
		{
			Name:      "backup",
			Aliases:   []string{"b"},
			Usage:     "Creates a backup",
			ArgsUsage: "<path>",
			Action:    s.Backup,
		},
		{
			Name:  "category",
			Usage: "Manage categories",
			Commands: []*ucli.Command{
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Usage:   "list all categories",
					Action:  s.CategoryList,
				},
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Deletes an item",
			ArgsUsage: "<name> <username>",
			Action:    s.Delete,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edits an item",
			ArgsUsage: "<name> <username>",
			Action:    s.Edit,
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generates a password for an item",
			ArgsUsage: "<name> <username>",
			Flags: []ucli.Flag{
				&ucli.BoolFlag{
					Name:    "advanced",
					Aliases: []string{"a"},
					Usage:   "Enable advanced mode",
				},
			},
			Action: s.Generate,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "Lists all websites",
			ArgsUsage: "<name>",
			Action:    s.List,
		},
		{
			Name:   "me",
			Usage:  "Displays default username and phone number",
			Action: s.Me,
		},
		{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Changes master password",
			Action:  s.Password,
		},
		{
			Name:      "restore",
			Aliases:   []string{"r"},
			Usage:     "Restores a backup",
			ArgsUsage: "<path>",
			Action:    s.Restore,
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Reapply config, such as encryption mode",
			Action:  s.Sync,
		},
		{
			Name:        "unclip",
			Usage:       "Internal command to clear clipboard",
			Description: "Clear the clipboard if the content matches the checksum.",
			Hidden:      true,
			Flags: []ucli.Flag{
				&ucli.IntFlag{
					Name:  "timeout",
					Usage: "Time to wait",
				}, &ucli.BoolFlag{
					Name:  "force",
					Usage: "Clear clipboard even if checksum mismatches",
				},
			},
			Action: s.Unclip,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Updates to the latest release",
			Action:  s.Update,
		},
	}
}
