package cli

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// Move cursor up relative the current position
func moveCursorUp(bias int) {
	out := windows.Handle(os.Stdout.Fd())
	var inMode uint32
	if err := windows.GetConsoleMode(out, &inMode); err == nil {
		// Validate that windows.ENABLE_VIRTUAL_TERMINAL_INPUT is supported, but do not set it.
		windows.SetConsoleMode(out, windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
	} else {
		fmt.Printf("failed to get console mode for stdin: %v\n", err)
	}

	fmt.Fprintf(os.Stdout, "\033[%dA", bias)
}

func clearToEnd() {
	fmt.Fprintf(os.Stdout, "\u001b[0J")
}
