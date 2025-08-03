package util

import (
	"bufio"
	"fmt"
	"os"
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
