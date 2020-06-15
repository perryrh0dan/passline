package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"

	ap "passline/pkg/action"
	"passline/pkg/config"
	"passline/pkg/ctxutil"

	"github.com/blang/semver"
	"github.com/fatih/color"
	ucli "github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

func setupApp(ctx context.Context, sv semver.Version) (context.Context, *ucli.App) {
	// try to load config
	cfg, err := config.Get()
	if err != nil {
		os.Exit(ap.ExitConfig)
	}

	// set config values
	ctx = initContext(ctx, cfg)

	action, err := ap.New(cfg, sv)
	if err != nil {
		os.Exit(ap.ExitUnknown)
	}

	app := ucli.NewApp()
	app.Name = "Passline"
	app.Usage = "Password manager"
	app.HelpName = "passline"
	app.Version = sv.String()
	app.Description = "Password manager for the command line"
	app.EnableBashCompletion = true

	// Append website information to default helper print
	app.CustomAppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: 
   https://github.com/perryrh0dan/passline

	`, ucli.AppHelpTemplate)

	// default command to get password
	app.Action = func(c *ucli.Context) error {
		return action.Default(c)
	}

	app.Flags = []ucli.Flag{
		&ucli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force displaying content",
		},
		&ucli.BoolFlag{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Print the generated password to the terminal",
		},
		&ucli.BoolFlag{
			Name:  "yes",
			Usage: "Assume yes on all yes/no questions or use the default on all others",
		},
	}

	app.Commands = action.GetCommands()

	sort.Sort(ucli.FlagsByName(app.Flags))
	sort.Sort(ucli.CommandsByName(app.Commands))

	return ctx, app
}

func initContext(ctx context.Context, cfg *config.Config) context.Context {
	// initialize from config, may be overridden by env vars
	ctx = cfg.WithContext(ctx)

	// support for no-color.org
	if nc := os.Getenv("NO_COLOR"); nc != "" {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
	}

	// only emit color codes when stdout is a terminal
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
		ctx = ctxutil.WithTerminal(ctx, false)
		ctx = ctxutil.WithInteractive(ctx, false)
	}

	// reading from stdin?
	if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeCharDevice == 0 {
		ctx = ctxutil.WithInteractive(ctx, false)
		ctx = ctxutil.WithStdin(ctx, true)
	}

	// disable colored output on windows since cmd.exe doesn't support ANSI color
	// codes. Other terminal may do, but until we can figure that out better
	// disable this for all terms on this platform
	if runtime.GOOS == "windows" {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
	}

	return ctx
}
