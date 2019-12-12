package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(message string, defaultValue string) (string, error) {
	// Print message

	if defaultValue != "" {
		fmt.Printf(message, defaultValue)
	} else {
		fmt.Print(message)
	}

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	// TODO Test if working for Linux
	text = strings.TrimSuffix(text, "\r")

	// If input empty assign default value also valid if defaultValue is empty
	if text == "" {
		text = defaultValue
	}
	return text, nil
}
