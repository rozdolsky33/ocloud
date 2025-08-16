package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptYesNo prompts the user with a yes or no question and returns true for 'yes' and false for 'no'.
func PromptYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", question)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			fmt.Println("Please enter y or n.")
		}
	}
}

// PromptPort prompts the user to enter a TCP port. If the user enters empty input, the defaultPort is returned.
// It validates the port is in range [1, 65535].
func PromptPort(question string, defaultPort int) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		if defaultPort > 0 {
			fmt.Printf("%s [%d]: ", question, defaultPort)
		} else {
			fmt.Printf("%s: ", question)
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		input = strings.TrimSpace(input)
		if input == "" && defaultPort > 0 {
			return defaultPort, nil
		}
		p, err := strconv.Atoi(input)
		if err != nil || p < 1 || p > 65535 {
			fmt.Println("Please enter a valid port number between 1 and 65535.")
			continue
		}
		return p, nil
	}
}

// PromptString prompts the user to enter a string. If the user enters empty input and defaultVal is provided, defaultVal is returned.
func PromptString(question string, defaultVal string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if strings.TrimSpace(defaultVal) != "" {
		fmt.Printf("%s [%s]: ", question, defaultVal)
	} else {
		fmt.Printf("%s: ", question)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}
