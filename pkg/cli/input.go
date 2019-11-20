package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(message string, values []string) string {
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
