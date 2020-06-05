package selection

import (
	"context"
	"fmt"
	"os"

	"passline/pkg/cli/list"
	"passline/pkg/cli/screenbuf"
	"passline/pkg/cli/terminal"
	"passline/pkg/util"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/gopasspw/gopass/pkg/action"
	ucli "github.com/urfave/cli/v2"
)

func ArgOrSelect(ctx context.Context, args ucli.Args, index int, message string, items []string) (string, error) {
	userInput := ""
	if args.Len()-1 >= index {
		userInput = args.Get(index)

		// if input is no item name use as filter
		if !util.ArrayContains(items, userInput) {
			items = util.FilterArray(items, userInput)
			if len(items) == 0 {
				return "", action.ExitError(ctx, action.ExitNotFound, nil, "No items with filter: %s found", userInput)
			}
			userInput = ""
		}
	}
	if userInput == "" {
		if len(items) > 1 {
			message := fmt.Sprintf("Please select a %s: ", message)
			selection, err := Default(message, items)
			if err != nil {
				return "", err
			}
			if selection == -1 {
				os.Exit(1)
			}

			userInput = items[selection]
			terminal.MoveCursorUp(1)
			terminal.ClearLines(1)
			fmt.Printf("%s%s\n", message, userInput)
		} else if len(items) == 1 {
			fmt.Printf("Selected %s: %s\n", message, items[0])
			return items[0], nil
		}
	}

	return userInput, nil
}

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

	sb.Reset()
	sb.Clear()
	sb.Flush()
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
