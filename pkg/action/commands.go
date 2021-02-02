package action

import (
	ucli "github.com/urfave/cli/v2"
)

// GetCommands returns the ucli commands exported by this module
func (s *Action) GetCommands() []*ucli.Command {
	return []*ucli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Adds an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Flags: []ucli.Flag{
				&ucli.BoolFlag{
					Name:    "advanced",
					Aliases: []string{"a"},
					Usage:   "Enable advanced mode",
				},
			},
			Action: s.Add,
		},
		{
			Name:      "backup",
			Aliases:   []string{"b"},
			Usage:     "Creates a backup",
			ArgsUsage: "<path>",
			Action:    s.Backup,
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
