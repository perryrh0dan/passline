package action

import (
	ucli "github.com/urfave/cli/v2"
)

// GetCommands returns the ucli commands exported by this module
func (s *Action) GetCommands() []*ucli.Command {
	return []*ucli.Command{
		{
			Name:      "backup",
			Aliases:   []string{"b"},
			Usage:     "Create a backup",
			ArgsUsage: "<path>",
			Action:    s.Backup,
		},
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Add an existing password for a website",
			ArgsUsage: "<name> <username> <password>",
			Action:    s.Add,
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete an item",
			ArgsUsage: "<name> <username>",
			Action:    s.Delete,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "Edit an item",
			ArgsUsage: "<name> <username>",
			Action:    s.Edit,
		},
		{
			Name:      "generate",
			Aliases:   []string{"g"},
			Usage:     "Generate a password for an item",
			ArgsUsage: "<name> <username>",
			Flags: []ucli.Flag{
				&ucli.BoolFlag{
					Name:    "advanced",
					Aliases: []string{"a"},
					Usage:   "Enable advanced mode",
				},
				&ucli.BoolFlag{
					Name:    "print",
					Aliases: []string{"p"},
					Usage:   "Print the generated password to the terminal",
				},
			},
			Action: s.Generate,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    s.List,
		},
		{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Change master password",
			Action:  s.Password,
		},
		{
			Name:      "restore",
			Aliases:   []string{"r"},
			Usage:     "Restore backup",
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
			Usage:   "Update to the latest release",
			Action:  s.Update,
		},
	}
}
