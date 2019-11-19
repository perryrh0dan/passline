package ui

import (
	"bufio"
	"fmt"
	"github.com/perryrh0dan/passline/pkg/cli"
	"os"
	"strings"
)

func GetInput(message string, values []string) string {
	// Print message
	if len(values) == 0 {
		fmt.Print(message)
	} else {
		fmt.Printf(message, values[0])
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	// TODO Test if working for Linux
	text = strings.TrimSuffix(text, "\r")
	return text
}

func GetSelect(message string, items []string) string {
	selection := cli.Select(message, items)
	return items[selection]
}

func GetArgOrInput(args []string, index int, message string, values []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		input = GetInput(message, values)
	}

	return input, nil
}

func GetArgOrSelect(args []string, index int, message string, items []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		input = GetSelect(message, items)
	}

	return input, nil

}
