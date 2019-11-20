package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/k0kubun/go-ansi"
)

// Move cursor up relative the current position
func moveCursorUp(n int) {
	if runtime.GOOS == "windows" {
		ansi.CursorUp(n)
	} else {
		fmt.Fprintf(os.Stdout, "\033[%dA", n)
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
