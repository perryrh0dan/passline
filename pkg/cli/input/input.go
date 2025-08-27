package input

import (
	"bufio"
	"fmt"
	"os"
	"passline/pkg/crypt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	ucli "github.com/urfave/cli/v3"
	"golang.org/x/term"
)

func ArgOrInput(args ucli.Args, index int, message, defaultValue, rules string) (string, error) {
	userInput := args.Get(index)

	if userInput == "" {
		message := fmt.Sprintf("Please enter a %s", message)
		if defaultValue != "" {
			message += " [%s]: "
		} else {
			message += ": "
		}

		var err error
		userInput, err = Default(message, defaultValue, rules)
		if err != nil {
			return "", err
		}
	}

	return userInput, nil
}

func Default(message, defaultValue, rules string) (string, error) {
	// input validation
	var input string
	for true {
		// find if %s is in string
		rgx := regexp.MustCompile("%s")
		matches := rgx.FindAllStringIndex(message, -1)

		// print message
		if len(matches) == 0 {
			fmt.Print(message)
		} else {
			fmt.Printf(message, defaultValue)
		}

		reader := bufio.NewReader(os.Stdin)

		var err error
		input, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r") // TODO Test if working for Linux

		// if input empty assign default value also valid if defaultValue is empty
		if input == "" {
			input = defaultValue
		}

		if validate(input, rules) {
			break
		}
	}

	return input, nil
}

func Confirmation(message string) bool {
	result := ""
	for result != "y" && result != "n" {
		var err error
		result, err = Default(message, "", "required")
		if err != nil {
			return false
		}
	}

	return result == "y"
}

func Password(message string) []byte {
	// Now get the password.
	fmt.Print(message)
	p, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		if runtime.GOOS != "windows" {
			panic(err)
		}
	}

	// Return the password as a string.
	return p
}

func MasterPassword(encryptedEncryptionKey string, reason string) ([]byte, error) {
	// If encrypted encryption key exists decrypt it
	envKey := []byte(os.Getenv("PASSLINE_MASTER_KEY"))
	if len(envKey) > 0 {
		encryptionKey, err := crypt.DecryptKey(envKey, encryptedEncryptionKey)
		if err == nil {
			return []byte(encryptionKey), nil
		}
	}

	prompt := "Enter master password: "
	if reason != "" {
		prompt = fmt.Sprintf("Enter master password %s: ", reason)
	}

	counter := 0
	for counter < 3 {
		password := Password(prompt)
		fmt.Println()

		encryptionKey, err := crypt.DecryptKey(password, encryptedEncryptionKey)
		if err == nil {
			return []byte(encryptionKey), nil
		} else if counter != 2 {
			fmt.Println("Wrong password! Please try again")
		}

		counter++
	}

	return []byte{}, fmt.Errorf("Wrong password")

}

func validate(input string, meta string) bool {
	input = strings.TrimSpace(input)
	rules := strings.Split(meta, ",")

	required, _ := regexp.Compile("required")
	length, _ := regexp.Compile("length:\\d+")
	number, _ := regexp.Compile("number:\\d+")

	valid := true

	for _, rule := range rules {
		if required.MatchString(rule) {
			if len(input) < 1 {
				valid = false
			}
		} else if length.MatchString(rule) {
			value, _ := strconv.Atoi(strings.Split(rule, ":")[1])
			if len(input) < value {
				valid = false
			}
		} else if number.MatchString(rule) {
			value, _ := strconv.Atoi(strings.Split(rule, ":")[1])
			number, _ := strconv.Atoi(input)
			if number < value {
				valid = false
			}
		}
	}

	return valid
}
