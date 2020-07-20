package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// For compatibility with Windows ANSI/UNICODE/etcetc
// and if Linux ever decides to use weird encoding
// Usage: newString, err := decodeString("bruhhh")
//        if err != nil { [handler] }
//        [do something with newString]
func decodeString(fileContent string) (string, error) {
	return fileContent, nil
}

// Read a file into a string
// Usage: contents, err := readFile("/etc/test")
//        if err != nil { [handler] }
//        [do something with contents]
func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	return string(fileContent), err
}

var aeacusVersion = "1.1.1"

func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func timeCheck(mc *metaConfig) {
	if mc.Config.EndDate != "" {
		endDate, err := time.Parse("2006/01/02 15:04:05 MST", mc.Config.EndDate)
		if err != nil {
			failPrint("Your EndDate value in the configuration is invalid.")
		} else {
			if time.Now().After(endDate) {
				destroyImage(mc)
			}
		}
	}
}

func runningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}
