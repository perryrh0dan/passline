package action

import (
	"log"
	"os"
	"passline/pkg/cli/input"
	"passline/pkg/renderer"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	ucli "github.com/urfave/cli/v2"
)

const (
	repo = "perryrh0dan/passline"
)

func (s *Action) Update(c *ucli.Context) error {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		renderer.DetectVersionError(err)
		return err
	}

	if !found || latest.Version.LTE(s.version) {
		renderer.NoUpdatesFound()
		return nil
	}

	message := "Do you want to update to: " + s.version.String() + "? (y/n): "
	confirm, err := input.Confirmation(message)

	if !confirm {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return err
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		renderer.UpdateError(err)
		return err
	}

	renderer.SuccessfulUpdated(latest.Version.String())
	renderer.DisplayReleaseNotes(latest.ReleaseNotes)

	return nil
}
