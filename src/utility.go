package main

import (
	"io/ioutil"
	"regexp"
)

const (
	aeacusVersion = "1.4.0"
	scoringConf   = "scoring.conf"
	scoringData   = "scoring.dat"
	linuxDir      = "/opt/aeacus/"
	windowsDir    = "C:\\aeacus\\"
)

var (
	verboseEnabled = false
	debugEnabled   = false
	yesEnabled     = false
	mc             = &metaConfig{}
)

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		failPrint("Error writing file: " + err.Error())
	}
}

// grepString acts like grep, taking in a pattern to search for, and the
// fileText to search in. It returns the line which contains the string
// (if any).
func grepString(patternText, fileText string) string {
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + patternText + ".*$")
	return string(re.Find([]byte(fileText)))
}
