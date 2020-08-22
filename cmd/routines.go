package cmd

import (
	"errors"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

// readScoringData is a convenience function around readData and decodeString,
// which parses the encrypted scoring configuration file.
func ReadScoringData() error {
	infoPrint("Decrypting data from " + mc.DirPath + ScoringData + "...")
	decryptedData, err := readData(ScoringData)
	if err != nil {
		failPrint("error reading in scoring data: " + err.Error())
		return err
	} else if decryptedData == "" {
		failPrint("scoring data is empty! Is the file corrupted?")
		return errors.New("Scoring data is empty!")
	} else {
		infoPrint("Data decrypting successful!")
	}
	parseConfig(decryptedData)
	return nil
}

// checkConfig parses and checks the validity of the current ScoringConf file.
func CheckConfig(fileName string) {
	fileContent, err := readFile(mc.DirPath + fileName)
	if err != nil {
		failPrint("Configuration file (" + fileName + ") not found!")
		os.Exit(1)
	}
	parseConfig(fileContent)
	printConfig()
	obfuscateConfig()
}

// fillConstants determines the correct constants, such as DirPath, for the
// given runtime and environment.
func FillConstants() {
	if runtime.GOOS == "linux" {
		mc.DirPath = LinuxDir
	} else if runtime.GOOS == "windows" {
		mc.DirPath = WindowsDir
	} else {
		failPrint("This operating system (" + runtime.GOOS + ") is not supported!")
		os.Exit(1)
	}
}

func ScoreImage() {
	checkTrace()
	timeCheck()
	infoPrint("Scoring image...")
	scoreImage()
}

// runningPermsCheck is a convenience function wrapper around
// adminCheck, which prints an error indicating that admin
// permissions are needed.
func RunningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}

// ParseFlags sets the global variable values, for example, verboseEnabled.
func ParseFlags(c *cli.Context) {
	if c.Bool("v") {
		verboseEnabled = true
	}
	if c.Bool("d") {
		debugEnabled = true
	}
	if c.Bool("y") {
		YesEnabled = true
	}
}
