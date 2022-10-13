package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/mail"
	"strconv"
	"syscall"
)

func getTextInput(question string, required bool) string {
	var output string

	fmt.Println(question)
	fmt.Printf("> " + colorGreen)
	fmt.Scanln(&output)
	fmt.Println(colorReset)

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getTextInput(question, required)
	}

	return output
}

func getNumberInput(question string, required bool) int {
	var output int
	var err error

	text := getTextInput(question, required)
	output, err = strconv.Atoi(text)
	if err != nil {
		fmt.Println(colorRed + "Please enter a valid port" + colorReset)
		output = getNumberInput(question, required)
	}

	return output
}

func getConfirmInput(question string) bool {
	fmt.Println(question)
	output := askForConfirmation(true)
	fmt.Println(colorReset)

	return output
}

func getEmailAddressInput(question string, required bool) string {
	var output string

	fmt.Println(question)
	fmt.Printf("> " + colorGreen)
	fmt.Scanln(&output)
	fmt.Println(colorReset)

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getEmailAddressInput(question, required)
	}

	_, err := mail.ParseAddress(output)

	if err != nil {
		fmt.Println(colorRed + "The email address is not valid. Please try again." + colorReset)
		output = getEmailAddressInput(question, required)
	}

	return output
}

func getPasswordInput(question string, required bool) string {
	var output string

	fmt.Println(question)
	fmt.Printf("> " + colorGreen)
	passwd, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(colorReset)

	if err != nil {
		fmt.Println(colorRed + "An error occurred. Please try again." + colorReset)
		output = getPasswordInput(question, required)
	} else {
		output = string(passwd)
	}

	if required == true && output == "" {
		fmt.Println(colorRed + "An empty value is not allowed here" + colorReset)
		output = getPasswordInput(question, required)
	}

	return output
}

func printSectionHeader(title string) {
	fmt.Println("\n\n\n --- " + colorCyan + title + colorReset + " --- \n")
}
