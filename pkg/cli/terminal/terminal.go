package terminal

import (
	"fmt"
	"os"
	"runtime"

	"github.com/k0kubun/go-ansi"
)

// Move cursor up relative the current position
func MoveCursorUp(n int) {
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

func ClearLines(lines int) {
	if runtime.GOOS == "windows" {
		for i := 1; i <= lines; i++ {
			MoveCursorUp(1)
			ansi.EraseInLine(3)
		}
	} else {
		MoveCursorUp(lines)
		fmt.Fprintf(os.Stdout, "\u001b[0J")
	}
}

func GetCursor() {
	fmt.Fprintf(os.Stdout, "\033[6n")
}

func HideCursor() {
	ansi.CursorHide()
}

func ShowCursor() {
	ansi.CursorShow()
}
