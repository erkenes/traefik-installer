package main

import (
	"fmt"
)

/*
Ask for confirmation
*/
func askForConfirmation(withPrefix bool) bool {
	var response string
	if withPrefix == true {
		fmt.Printf("> " + colorGreen)
	}
	fmt.Scanln(&response)
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println(colorRed + "Please type yes or no and then press enter:" + colorReset)
		return askForConfirmation(withPrefix)
	}
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}
