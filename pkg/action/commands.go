package action

import (
	"github.com/urfave/cli/v2"
)

// GetCommands returns the cli commands exported by this module
func (s *Action) GetCommands() []*cli.Command {
	return []*cli.Command{
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
			Flags: []cli.Flag{
				&cli.BoolFlag{
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
			Usage:     "List all items",
			ArgsUsage: "<name>",
			Action:    s.List,
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
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "timeout",
					Usage: "Time to wait",
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
