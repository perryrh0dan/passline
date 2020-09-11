package out

import "github.com/fatih/color"

func UpdateFound() {
	d := color.New(color.FgGreen)
	d.Printf("A new release is available\n")
}
