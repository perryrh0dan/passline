// +build windows

package selection

import (
	"fmt"

	"passline/pkg/cli/list"
	"passline/pkg/cli/terminal"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func Default(message string, items []string) (int, error) {
	// Print Message
	fmt.Println(message)

	list, err := list.New(items, 10)
	if err != nil {
		return 0, err
	}

	selected := -1

	// Open keyboard
	err = keyboard.Open()
	if err != nil {
		return selected, err
	}
	defer keyboard.Close()

	printList(list)

	// Hide Cursor
	terminal.HideCursor()
	defer terminal.ShowCursor()
	var open = true

	for open {
		_, key, _ := keyboard.GetKey()
		update := false
		switch key {
		case keyboard.KeyEsc:
			open = false
		case keyboard.KeyCtrlC:
			open = false
		case keyboard.KeyEnter:
			selected = list.Index()
			open = false
		case keyboard.KeyArrowUp:
			list.Prev()
			update = true
		case keyboard.KeyArrowDown:
			list.Next()
			update = true
		}

		if update {
			clearScreen(list)
			printList(list)
		}
	}

	clearScreen(list)
	return selected, nil
}

func printList(list *list.List) {
	values, selected := list.Items()

	items := make([]string, len(values))
	for i, v := range values {
		items[i] = fmt.Sprint(v)
	}

	for index, item := range items {
		if index != selected {
			text := fmt.Sprintf("[ ] %s", item)
			fmt.Println(text)
		} else {
			d := color.New(color.FgGreen)
			text := d.Sprintf("[x] %s", item)
			fmt.Println(text)
		}
	}
}

func clearScreen(list *list.List) {
	values, _ := list.Items()

	for i := 1; i <= len(values); i++ {
		ansi.CursorUp(1)
		ansi.EraseInLine(3)
	}
}
