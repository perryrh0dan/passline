package cli

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/k0kubun/go-ansi"
	"golang.org/x/crypto/ssh/terminal"
)

func getPassword(prompt string) []byte {
	// Get the initial state of the terminal.
	initialTermState, e1 := terminal.GetState(int(syscall.Stdin))
	if e1 != nil {
		panic(e1)
	}

	// Restore it in the event of an interrupt.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = terminal.Restore(int(syscall.Stdin), initialTermState)
		fmt.Println()
		os.Exit(1)
	}()

	// Now get the password.
	fmt.Print(prompt)
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		if runtime.GOOS != "windows" {
			panic(err)
		}
	}

	// Stop looking for ^C on the channel.
	signal.Stop(c)

	// Return the password as a string.
	return p
}

// Move cursor up relative the current position
func moveCursorUp(n int) {
	if runtime.GOOS == "windows" {
		ansi.CursorUp(n)
	} else {
		// Move cursor
		fmt.Fprintf(os.Stdout, "\033[%dA", n)
		// Scroll up
		// fmt.Fprintf(os.Stdout, "\033[M")
		// Clear down
		fmt.Fprintf(os.Stdout, "\033[J")
	}
}

func clearLines(lines int) {
	if runtime.GOOS == "windows" {
		for i := 1; i <= lines; i++ {
			moveCursorUp(1)
			ansi.EraseInLine(3)
		}
	} else {
		moveCursorUp(lines)
		fmt.Fprintf(os.Stdout, "\u001b[0J")
	}
}

func getCursor() {
	fmt.Fprintf(os.Stdout, "\033[6n")
}

func hideCursor() {
	ansi.CursorHide()
}

func showCursor() {
	ansi.CursorShow()
}
