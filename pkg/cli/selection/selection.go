package selection

import (
	"context"
	"errors"
	"fmt"

	"passline/pkg/cli/terminal"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func ArgOrSelect(ctx context.Context, args ucli.Args, index int, message string, items []string) (string, error) {
	userInput := args.Get(index)

	// if input is no item name use as filter
	if util.ArrayContains(items, userInput) {
		return userInput, nil
	}
	items = util.FilterArray(items, userInput)
	if len(items) == 0 {
		return "", errors.New("Not found")
	}

	if len(items) == 1 {
		fmt.Printf("Selected %s: %s\n", message, items[0])
		return items[0], nil
	}
	message = fmt.Sprintf("Please select a %s: ", message)
	selection, err := Default(message, items)

	if err != nil {
		return "", err
	}

	if selection == -1 {
    print("test")
		return "", errors.New("Canceled selection")
	}

	terminal.ClearLines(1)
	fmt.Printf("%s%s\n", message, items[selection])
	return items[selection], nil
}
