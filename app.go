package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	ap "passline/pkg/action"
	"passline/pkg/config"
	"passline/pkg/ctxutil"
	"passline/pkg/util"

	"github.com/blang/semver"
	ucli "github.com/urfave/cli/v3"
	"golang.org/x/term"
)

func setupApp(ctx context.Context, sv semver.Version) (context.Context, *ucli.Command) {
	// try to load config
	cfg, err := config.Get(util.OSFileSystem{})
	if err != nil {
		os.Exit(ap.ExitConfig)
	}

	// set config values
	ctx = initContext(ctx, cfg)

	action, err := ap.New(cfg, sv)
	if err != nil {
		os.Exit(ap.ExitUnknown)
	}

	var app = ucli.Command{
		Name:                  "Passline",
		Usage:                 "Password manager",
		Version:               sv.String(),
		Description:           "Password manager for the command line",
		EnableShellCompletion: true,
	}

	// Append website information to default helper print
	app.CustomRootCommandHelpTemplate = fmt.Sprintf(`%s
WEBSITE: 
   https://github.com/perryrh0dan/passline`)

	app.Flags = []ucli.Flag{
		&ucli.BoolFlag{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Prints the password to the terminal",
		},
		&ucli.BoolFlag{
			Name:  "yes",
			Usage: "Assume yes on all yes/no questions or use the default on all others",
		},
		&ucli.StringFlag{
			Name:    "category",
			Aliases: []string{"c"},
			Usage:   "Select only items with given category",
		},
		&ucli.BoolFlag{
			Name:  "noclip",
			Usage: "Disable copy to clipboard",
		},
	}

	// default command to get password
	app.Action = func(c context.Context, command *ucli.Command) error {
		return action.Default(c, command)
	}

	app.Commands = action.GetCommands()

	sort.Sort(ucli.FlagsByName(app.Flags))

	return ctx, &app
}

func initContext(ctx context.Context, cfg *config.Config) context.Context {
	// initialize from config, may be overridden by env vars
	ctx = cfg.WithContext(ctx)

	// only emit color codes when stdout is a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		ctx = ctxutil.WithColor(ctx, false)
		ctx = ctxutil.WithTerminal(ctx, false)
		ctx = ctxutil.WithInteractive(ctx, false)
	}

	// reading from stdin?
	if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeCharDevice == 0 {
		ctx = ctxutil.WithInteractive(ctx, false)
		ctx = ctxutil.WithStdin(ctx, true)
	}

	return ctx
}
