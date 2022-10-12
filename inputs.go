package main

import (
	"fmt"
	"strconv"
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
	var configSmtpPortError error

	output, configSmtpPortError = strconv.Atoi(getTextInput(question, required))
	if configSmtpPortError != nil {
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

func printSectionHeader(title string) {
	fmt.Println("\n\n\n --- " + colorCyan + title + colorReset + " --- \n")
}
