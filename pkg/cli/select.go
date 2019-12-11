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

	// Open keyboard and Hide Cursor
	keyboard.Open()
	defer keyboard.Close()
	hideCursor()
	defer showCursor()

	for {
		_, key, _ := keyboard.GetKey()
		update := false
		switch key {
		case keyboard.KeyEsc:
			return -1, errors.New("Canceled")
		case 13:
			clearLines(len(items) + 1)
			return selected, nil
		case 65517:
			if selected > 0 {
				selected--
				update = true
			}
		case 65516:
			if selected < len(items)-1 {
				selected++
				update = true
			}
		}

		if update {
			moveCursorUp(len(items))
			printSelect(items, selected)
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
