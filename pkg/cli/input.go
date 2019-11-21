package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(message string, value string) (string, error) {
	// Print message

	if value != "" {
		fmt.Printf(message, value)
	} else {
		fmt.Printf(message)
	}

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	// TODO Test if working for Linux
	text = strings.TrimSuffix(text, "\r")
	return text, nil
}
