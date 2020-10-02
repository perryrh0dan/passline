package out

import (
	"github.com/blang/semver"
	"github.com/fatih/color"
)

func UpdateFound(actual, new semver.Version) {
	d := color.New(color.FgGreen)
	d.Printf("A new release is available: %s -> %s\n", actual, new)
}
