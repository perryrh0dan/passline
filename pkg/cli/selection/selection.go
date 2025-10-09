package selection

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"passline/pkg/cli/terminal"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v3"
)

type SelectItem struct {
	Value string
	Label string
}

type SelectItemWithDistance struct {
	item  SelectItem
	score int
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
	selectItemsWithDistance := make([]SelectItemWithDistance, 0)

	for _, i := range l {
		_, distance := util.LevenshteinDistanceSubstring(i.Label, filter)
		if distance <= max(0, min(len(filter)-2, 2)) {
			selectItemsWithDistance = append(selectItemsWithDistance, SelectItemWithDistance{item: i, score: distance})
		}
	}

	slices.SortFunc(selectItemsWithDistance, func(a, b SelectItemWithDistance) int {
		return a.score - b.score
	})

	filteredItems := make([]SelectItem, 0)
	for _, i := range selectItemsWithDistance {
		filteredItems = append(filteredItems, i.item)
	}

	return filteredItems
}
