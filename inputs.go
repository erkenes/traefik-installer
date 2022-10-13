package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/mail"
	"strconv"
	"syscall"
)

/*
Get text input
*/
func getTextInput(question string, required bool, value string) string {
	var output string

	fmt.Println(question)

	changeValue := true
	if value != "" {
		changeValue = getConfirmInput("The current value is "+colorCyan+value+colorReset+" do you want to change it?", false)
		output = value

		if changeValue {
			fmt.Println(question)
		}
	}

	if changeValue {
		fmt.Printf("> " + colorGreen)
		fmt.Scanln(&output)
		fmt.Println(colorReset)
	}

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getTextInput(question, required, "")
	}

	return output
}

/*
Get number input
*/
func getNumberInput(question string, required bool, value int) int {
	var output int
	var err error
	var valueForInput string

	if value == 0 {
		valueForInput = ""
	} else {
		valueForInput = strconv.Itoa(value)
	}
	text := getTextInput(question, required, valueForInput)
	output, err = strconv.Atoi(text)
	if err != nil {
		fmt.Println(colorRed + "Please enter a valid port" + colorReset)
		output = getNumberInput(question, required, 0)
	}

	return output
}

/*
Ask for confirmation
*/
func getConfirmInput(question string, value bool) bool {
	var output bool

	fmt.Println(question)

	changeValue := true
	if value {
		changeValue = getConfirmInput("Currently the option is active. Do you want to change it?", false)
		output = value

		if changeValue {
			fmt.Println(question)
		}
	}

	if changeValue {
		output = askForConfirmation(true)
		fmt.Println(colorReset)
	}

	return output
}

/*
Get email address and validate it
*/
func getEmailAddressInput(question string, required bool, value string) string {
	var output string

	fmt.Println(question)

	changeValue := true
	if value != "" {
		changeValue = getConfirmInput("The current value is "+colorCyan+value+colorReset+" do you want to change it?", false)
		output = value

		if changeValue {
			fmt.Println(question)
		}
	}

	if changeValue {
		fmt.Printf("> " + colorGreen)
		fmt.Scanln(&output)
		fmt.Println(colorReset)
	}

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getEmailAddressInput(question, required, "")
	}

	_, err := mail.ParseAddress(output)

	if err != nil {
		fmt.Println(colorRed + "The email address is not valid. Please try again." + colorReset)
		output = getEmailAddressInput(question, required, "")
	}

	return output
}

/*
Get password input (hidden input)
*/
func getPasswordInput(question string, required bool, value string) string {
	var output string
	var passwd []byte
	var err error

	fmt.Println(question)

	changeValue := true
	if value != "" {
		changeValue = getConfirmInput("The current value is "+colorCyan+value+colorReset+" do you want to change it?", false)
		passwd = []byte(value)

		if changeValue {
			fmt.Println(question)
		}
	}

	if changeValue {
		fmt.Printf("> " + colorGreen)
		passwd, err = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println(colorReset)
	}

	if err != nil {
		fmt.Println(colorRed + "An error occurred. Please try again." + colorReset)
		output = getPasswordInput(question, required, "")
	} else {
		output = string(passwd)
	}

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getPasswordInput(question, required, "")
	}

	return output
}

/*
Print a section header
*/
func printSectionHeader(title string) {
	fmt.Println("\n\n\n --- " + colorCyan + title + colorReset + " --- \n")
}
