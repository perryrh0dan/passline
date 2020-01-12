package cli

import (
	"errors"
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

func Select(message string, items []string) (int, error) {
	selected := 0

	// Print Message
	fmt.Println(message)

	// Print Initial Selection
	printSelect(items, selected)

	// Open keyboard
	err := keyboard.Open()
	if err != nil {
		return -1, err
	}
	defer keyboard.Close()

	// Hide Cursor
	hideCursor()
	defer showCursor()

	for {
		_, key, _ := keyboard.GetKey()
		update := false
		switch key {
		case keyboard.KeyEsc:
			return -1, errors.New("Canceled")
		case keyboard.KeyCtrlC:
			return -1, errors.New("Canceled")
		case keyboard.KeyEnter:
			clearLines(len(items) + 1)
			return selected, nil
		case keyboard.KeyArrowUp:
			if selected > 0 {
				selected--
				update = true
			}
		case keyboard.KeyArrowDown:
			if selected < len(items)-1 {
				selected++
				update = true
			}
		}

		if update {
			moveCursorUp(len(items))
			// printSelect(items, selected)
		}
	}
}

func printSelect(items []string, selected int) {
	for index, item := range items {
		if index != selected {
			fmt.Printf("[ ] %s\n", item)
		} else {
			d := color.New(color.FgGreen)
			d.Printf("[x] %s\n", item)
		}
	}
}
