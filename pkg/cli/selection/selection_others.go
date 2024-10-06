//go:build !windows
// +build !windows

package selection

import (
	"fmt"
	"os"

	"passline/pkg/cli/list"
	"passline/pkg/cli/screenbuf"
	"passline/pkg/cli/terminal"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
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

	selected := 0

	printList(list, sb)

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
			update = list.Prev()
		case keys.Down:
			update = list.Next()
		case keys.RuneKey:
			if key.String() == "j" {
				update = list.Next()
			} else if key.String() == "k" {
				update = list.Prev()
			}
		}

		if update {
			printList(list, sb)
		}

		return false, nil
	})

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

	printFooter(list, sb)

	sb.Flush()
}

func printFooter(list *list.List, sb *screenbuf.ScreenBuf) {
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
	sb.WriteString("")
	sb.WriteString(text)
}

func clearScreen(sb *screenbuf.ScreenBuf) {
	sb.Reset()
	sb.Clear()
	sb.Flush()
}
