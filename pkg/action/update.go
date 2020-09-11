package action

import (
	"log"
	"os"
	
	"passline/pkg/cli/input"
	"passline/pkg/out"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	ucli "github.com/urfave/cli/v2"
)

const (
	repo = "perryrh0dan/passline"
)

func (s *Action) Update(c *ucli.Context) error {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		out.DetectVersionError(err)
		return err
	}

	if !found || latest.Version.LTE(s.version) {
		out.NoUpdatesFound()
		return nil
	}

	message := "Do you want to update to: " + latest.Version.String() + "? (y/n): "
	confirm := input.Confirmation(message)

	if !confirm {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return err
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		out.UpdateError(err)
		return err
	}

	out.SuccessfulUpdated(latest.Version.String())
	out.DisplayReleaseNotes(latest.ReleaseNotes)

	return nil
}
