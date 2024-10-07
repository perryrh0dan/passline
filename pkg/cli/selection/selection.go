package selection

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"passline/pkg/cli/terminal"

	ucli "github.com/urfave/cli/v2"
)

type SelectItem struct {
	Value string
	Label string
}

func ArgOrSelect(ctx context.Context, args ucli.Args, index int, message string, items []SelectItem) (string, error) {
	userInput := args.Get(index)

	// if input is no item name use as filter
	if arrayContains(items, userInput) {
		return userInput, nil
	}
	items = filterArray(items, userInput)
	if len(items) == 0 {
		return "", errors.New("Not found")
	}

	if len(items) == 1 {
		fmt.Printf("Selected %s: %s\n", message, items[0].Value)
		return items[0].Value, nil
	}
	message = fmt.Sprintf("Please select a %s: ", message)
	selection, err := Default(message, items)

	if err != nil {
		return "", err
	}

	if selection == -1 {
		return "", errors.New("Canceled selection")
	}

	terminal.ClearLines(1)
	fmt.Printf("%s%s\n", message, items[selection].Value)
	return items[selection].Value, nil
}

func arrayContains(l []SelectItem, i string) bool {
	for _, li := range l {
		if li.Label == i {
			return true
		}
	}

	return false
}

func filterArray(l []SelectItem, filter string) []SelectItem {
	filteredNames := make([]SelectItem, 0)
	for _, i := range l {
		if strings.Contains(i.Label, filter) {
			filteredNames = append(filteredNames, i)
		}
	}
	return filteredNames
}
