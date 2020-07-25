package main

import (
	"io/ioutil"
	"os"
	"time"
)

var verboseEnabled = false
var aeacusVersion = "1.2.0"

// writeFile wraps ioutil's WriteFule function, and prints
// the error the screen if one occurs.
func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		failPrint("Error writing file: " + err.Error())
	}
}

// newImageData returns an empty/default imageData struct.
func newImageData() imageData {
	return imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0, []string{"green", "OK", "green", "OK", "green", "OK"}, false}
}

// clearImageData resets the imageData values pertaining to scoring.
func clearImageData(id *imageData) {
	id.Score = 0
	id.ScoredVulns = 0
	id.TotalPoints = 0
	id.Contribs = 0
	id.Detracts = 0
	id.Points = []scoreItem{}
	id.Penalties = []scoreItem{}
}

// timeCheck calls destroyImage if the configured EndDate for the image has
// passed. Its purpose is to dissuade or prevent people using an image after
// the round ends.
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

// runningPermsCheck is a convenience function wrapper around
// adminCheck, which prints an error indicating that admin
// permissions are needed.
func runningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}
