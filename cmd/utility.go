package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

const (
	AeacusVersion = "1.8.3"
	ScoringConf   = "scoring.conf"
	ScoringData   = "scoring.dat"
	LinuxDir      = "/opt/aeacus/"
	WindowsDir    = "C:\\aeacus\\"
)

var (
	YesEnabled            = false
	verboseEnabled        = false
	debugEnabled          = false
	mc                    = &metaConfig{}
	timeStart             = time.Now()
	timeWithoutID, _      = time.ParseDuration("0s")
	withoutIDThreshold, _ = time.ParseDuration("30m")
)

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

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func writeFile(fileName, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0o644)
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

func removeKeys(aeacusPath string) {
	keyFile := path.Join(aeacusPath, ".keys")
	if _, err := os.Stat(keyFile); err == nil {
		if err := os.Remove(keyFile); err != nil {
			failPrint("Failed to remove .keys file")
		}
	} else if os.IsNotExist(err) {
		failPrint("Failed to remove .keys file, does not exist")
	} else {
		failPrint("Failed to stat " + keyFile)
	}
}
