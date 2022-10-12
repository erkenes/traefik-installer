package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

const colorReset = "\033[0m"

const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

func writeFile(filepath string, filename string, content []byte, filePerm os.FileMode) {
	fullFilepath := filename

	if filepath != "." && filepath != "" {
		fullFilepath = filepath + "/" + filename

		if _, err := os.Stat(fullFilepath); os.IsNotExist(err) {
			os.MkdirAll(filepath, 0700)
		}
	}

	err := os.WriteFile(fullFilepath, content, filePerm)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}

func randomString(n int) string {
	var Rando = rand.Reader
	b := make([]byte, n)
	_, _ = Rando.Read(b)
	return hex.EncodeToString(b)
}
