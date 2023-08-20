//go:build windows
// +build windows

package selection

import (
	"fmt"

	"passline/pkg/cli/list"
	"passline/pkg/cli/terminal"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
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

	selected := 0

	printList(list)

	// Hide Cursor
	terminal.HideCursor()
	defer terminal.ShowCursor()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		update := false

		switch key.Code {
		case keys.Esc:
			selected = -1
			return true, nil
		case keys.CtrlC:
			selected = -1
			return true, nil
		case keys.Enter:
			selected = list.Index()
			return true, nil
		case keys.Up:
			list.Prev()
			update = true
		case keys.Down:
			list.Next()
			update = true
		case keys.RuneKey:
			if key.String() == "j" {
				list.Next()
				update = true
			} else if key.String() == "k" {
				list.Prev()
				update = true
			}
		}

		if update {
			clearScreen(list)
			printList(list)
		}

		return false, nil
	})

	clearScreen(list)
	cursorStartOfLine()

	return selected, nil
}

func printList(list *list.List) {
	values, selected := list.Items()

	items := make([]string, len(values))
	for i, v := range values {
		items[i] = fmt.Sprint(v)
	}

	for index, item := range items {
		cursorStartOfLine()

		if index != selected {
			text := fmt.Sprintf("[ ] %s", item)
			fmt.Println(text)
		} else {
			d := color.New(color.FgGreen)
			d.Printf("[x] %s\n", item)
		}
	}

	printFooter(list)
}

func printFooter(list *list.List) {
	var from = list.Start() + 1
	var size = list.Size()
	var length = list.Length()

	var to int
	if size < length {
		to = list.Start() + size
	} else {
		to = length
	}

	text := fmt.Sprintf("Items %d - %d of %d", from, to, length)
	fmt.Println()
	cursorStartOfLine()
	fmt.Println(text)
}

func clearScreen(list *list.List) {
	values, _ := list.Items()

	// Clear Footer
	ansi.EraseInLine(3)
	ansi.CursorUp(1)
	ansi.EraseInLine(3)

	ansi.CursorUp(1)
	ansi.EraseInLine(3)

	// Clear list items
	for i := 1; i <= len(values); i++ {
		ansi.CursorUp(1)
		ansi.EraseInLine(3)
	}
}

func cursorStartOfLine() {
	ansi.CursorHorizontalAbsolute(0)
}
