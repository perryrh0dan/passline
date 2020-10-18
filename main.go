package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"passline/pkg/out"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	name = "passline"
	repo = "perryrh0dan/passline"
)

var (
	// Version is the released version of passline
	version string = "1.4.0"
	// BuildTime is the time the binary was built
	date string
)

func main() {
	ctx := context.Background()

	// Get the initial state of the terminal.
	initialTermState, _ := terminal.GetState(int(syscall.Stdin))

	//trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
			exit(initialTermState)
		case <-ctx.Done():
			cancel()
			exit(initialTermState)
		}
	}()

	sv := getVersion()

	// check for updates
	latest, found, _ := selfupdate.DetectLatest(repo)
	if found && latest.Version.GT(sv) {
		out.UpdateFound(sv, latest.Version)
	}

	ctx, app := setupApp(ctx, sv)
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func exit(initialTermState *terminal.State) {
	_ = terminal.Restore(int(syscall.Stdin), initialTermState)
	fmt.Println()
	os.Exit(1)
}

func getVersion() semver.Version {
	sv, err := semver.Parse(strings.TrimPrefix(version, "v"))
	if err == nil {
		return sv
	}

	return semver.Version{
		Major: 1,
		Minor: 9,
		Patch: 2,
		Pre: []semver.PRVersion{
			{VersionStr: "git"},
		},
		Build: []string{"HEAD"},
	}
}
