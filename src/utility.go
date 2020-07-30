package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"time"
)

var aeacusVersion = "1.3.0"

var verboseEnabled = false
var reverseEnabled = false
var scoringConf = "scoring.conf"
var scoringData = "scoring.dat"
var mc = &metaConfig{}

// writeFile wraps ioutil's WriteFule function, and prints
// the error the screen if one occurs.
func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		failPrint("Error writing file: " + err.Error())
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

func grepString(patternText, fileText string) string {
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + patternText + ".*$")
	return string(re.Find([]byte(fileText)))
}

func fillConstants() {
	if runtime.GOOS == "linux" {
		mc.DirPath = "/opt/aeacus/"
	} else if runtime.GOOS == "windows" {
		mc.DirPath = "C:\\aeacus\\"
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
