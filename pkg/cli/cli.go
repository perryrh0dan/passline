package cli

func ArgOrInput(args []string, index int, message string, values []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		input = Input(message, values)
	}

	return input, nil
}

func ArgOrSelect(args []string, index int, message string, items []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		selection := Select(message, items)
		input = items[selection]
	}

	return input, nil
}
