// Package ui contains general use case ui printing methods.
// Additionally includes input collectors using bufio.Reader with validation
package ui

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Input struct {
	Reader *bufio.Reader
}

// String gets the raw user input without any validations or parsing.
func (input *Input) String(prompt string) string {
	fmt.Printf("%s", prompt)
	userInput, _ := input.Reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

// Int gets an integer input from the user, validation included.
// The function will keep looping until a valid integer between min and max is inputted.
func (input *Input) Int(prompt string, min, max int) int {
	for {
		fmt.Printf("%s", prompt)
		userInput, _ := input.Reader.ReadString('\n')
		parsed, err := strconv.Atoi(strings.TrimSpace(userInput))
		if err != nil {
			fmt.Printf("INPUT ERROR. Please enter a valid integer.\n")
		} else if parsed < min || parsed > max {
			fmt.Printf("INPUT ERROR. Please enter an integer between %d and %d.\n", min, max)
		} else {
			return parsed
		}
	}
}

// Float gets a float64 input from the user, validation included.
// Same concept as the previous function but with the float64 data type.
func (input *Input) Float(prompt string, min, max float64) float64 {
	for {
		fmt.Printf("%s", prompt)
		userInput, _ := input.Reader.ReadString('\n')
		parsed, err := strconv.ParseFloat(strings.TrimSpace(userInput), 64)
		if err != nil {
			fmt.Printf("INPUT ERROR. Please enter a valid decimal number.\n")
		} else if parsed < min || parsed > max {
			fmt.Printf("INPUT ERROR. Please enter a decimal number between %.2f and %.2f.\n", min, max)
		} else {
			return parsed
		}
	}
}

// Buffer pauses the program temporarily, giving the user time to react.
func (input *Input) Buffer() {
	fmt.Printf("Press enter to proceed. ")
	input.Reader.ReadString('\n')
}
