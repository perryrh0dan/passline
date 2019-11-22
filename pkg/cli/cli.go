package cli

import (
	"fmt"
)

func ArgOrInput(args []string, index int, message string, value string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		message := fmt.Sprintf("Please enter a %s []: ", message)
		var err error
		input, err = Input(message, value)
		if err != nil {
			return "", err
		}
	}

	return input, nil
}

func ArgOrSelect(args []string, index int, message string, items []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		if len(items) > 1 {
			message := fmt.Sprintf("Please select a %s: ", message)
			selection, err := Select(message, items)
			if err != nil {
				return "", err
			}

			input = items[selection]
			fmt.Printf("%s%s\n", message, input)
		} else if len(items) == 1 {
			fmt.Printf("Selected %s: %s\n", message, items[0])
			return items[0], nil
		}
	}

	return input, nil
}
