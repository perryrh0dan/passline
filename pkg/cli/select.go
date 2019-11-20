package cli

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

func Select(message string, items []string) int {
	selected := 0
	fmt.Println(message)
	printSelect(items, selected)
	keyboard.Open()
	defer keyboard.Close()
	for {
		_, key, _ := keyboard.GetKey()
		update := false
		switch key {
		case 13:
			clearLines(len(items) + 1)
			return selected
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
			fmt.Println(item)
		} else {
			d := color.New(color.FgGreen)
			d.Printf(item + "\n")
		}
	}
}
