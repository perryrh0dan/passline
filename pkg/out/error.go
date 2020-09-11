package out

import "github.com/fatih/color"

func PasswordTooShort() {
	d := color.New(color.FgRed)
	d.Printf("Password must be at least 6 characters long\n")
}

func InvalidGeneratorOptions() {
	d := color.New(color.FgRed)
	d.Printf("Invalid options\n")
}
