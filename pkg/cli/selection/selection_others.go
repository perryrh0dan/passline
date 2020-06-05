// +build !windows

package selection

import (
	"fmt"
	"os"

	"passline/pkg/cli/list"
	"passline/pkg/cli/screenbuf"
	"passline/pkg/cli/terminal"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

func Default(message string, items []string) (int, error) {
	// Print Message
	fmt.Println(message)

	sb := screenbuf.New(os.Stdout)

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

	printList(list, sb)

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
			printList(list, sb)
		}
	}

	clearScreen(sb)
	return selected, nil
}

func printList(list *list.List, sb *screenbuf.ScreenBuf) {
	values, selected := list.Items()

	items := make([]string, len(values))
	for i, v := range values {
		items[i] = fmt.Sprint(v)
	}

	for index, item := range items {
		if index != selected {
			text := fmt.Sprintf("[ ] %s", item)
			sb.WriteString(text)
		} else {
			d := color.New(color.FgGreen)
			text := d.Sprintf("[x] %s", item)
			sb.WriteString(text)
		}
	}

	sb.Flush()
}

func clearScreen(sb *screenbuf.ScreenBuf) {
	sb.Reset()
	sb.Clear()
	sb.Flush()
}