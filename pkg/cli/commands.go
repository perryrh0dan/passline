package cli

import (
	"fmt"
	"os"
)

// Move cursor up relative the current position
func moveCursorUp(bias int) {
	fmt.Fprintf(os.Stdout, "\033[%dA", bias)
}

func clearToEnd() {
	fmt.Fprintf(os.Stdout, "\u001b[0J")
}
