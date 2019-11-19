package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/pkg/term"
)

func Select(message string, items []string) int {
	selected := 0
	fmt.Println(message)
	printSelect(items, selected)
	for {
		ascii, keyCode, _ := getChar()
		if ascii == 13 {
			moveCursorUp(len(items) + 1)
			clearToEnd()
			return selected
		}
		switch keyCode {
		case 38:
			if selected > 0 {
				selected -= 1
			}
		case 40:
			if selected < len(items)-1 {
				selected = selected + 1
			}
		}
		moveCursorUp(len(items))
		printSelect(items, selected)
	}
}

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}

func printSelect(items []string, selected int) {
	for index, item := range items {
		if index != selected {
			fmt.Println(item)
		} else {
			d := color.New(color.FgGreen)
			d.Printf(item + "\n")
		}
	}
}
