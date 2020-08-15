package main

import (
	"errors"
	"os"
	"runtime"
	"time"
)

// readScoringData is a convenience function around readData and decodeString,
// which parses the encrypted scoring configuration file.
func readScoringData() error {
	decryptedData, err := readData(scoringData)
	if err != nil {
		failPrint("Error reading in scoring data: " + err.Error())
		return err
	} else if decryptedData == "" {
		failPrint("Scoring data is empty! Is the file corrupted?")
		return errors.New("Scoring data is empty!")
	}
	parseConfig(decryptedData)
	return nil
}

// checkConfig parses and checks the validity of the current
// `scoring.conf` file.
func checkConfig(fileName string) {
	fileContent, err := readFile(mc.DirPath + fileName)
	if err != nil {
		failPrint("Configuration file (" + fileName + ") not found!")
		os.Exit(1)
	}
	parseConfig(fileContent)
	printConfig()
}

// fillConstants determines the correct constants, such as DirPath, for the
// given runtime and environment.
func fillConstants() {
	if runtime.GOOS == "linux" {
		mc.DirPath = linuxDir
	} else if runtime.GOOS == "windows" {
		mc.DirPath = windowsDir
	} else {
		failPrint("This operating system (" + runtime.GOOS + ") is not supported!")
		os.Exit(1)
	}
}

// runningPermsCheck is a convenience function wrapper around
// adminCheck, which prints an error indicating that admin
// permissions are needed.
func runningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}

// timeCheck calls destroyImage if the configured EndDate for the image has
// passed. Its purpose is to dissuade or prevent people using an image after
// the round ends.
func timeCheck() {
	if mc.Config.EndDate != "" {
		endDate, err := time.Parse("2006/01/02 15:04:05 MST", mc.Config.EndDate)
		if err != nil {
			failPrint("Your EndDate value in the configuration is invalid.")
		} else {
			if time.Now().After(endDate) {
				destroyImage()
			}
		}
	}
}
