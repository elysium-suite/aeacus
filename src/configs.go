package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

func parseConfig(configContent string) {
	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeConfig(sourceFile, destFile string) {
	if verboseEnabled {
		infoPrint("Reading configuration from " + mc.DirPath + sourceFile + "...")
	}

	configFile, err := readFile(mc.DirPath + sourceFile)
	if err != nil {
		failPrint("Can't open scoring configuration file (" + sourceFile + "): " + err.Error())
		os.Exit(1)
	}
	encryptedConfig := encryptConfig(configFile)
	if verboseEnabled {
		infoPrint("Writing data to " + mc.DirPath + "...")
	}
	writeFile(mc.DirPath+destFile, encryptedConfig)
}

func readData(fileName string) string {
	if verboseEnabled {
		infoPrint("Decrypting data from " + mc.DirPath + fileName + "...")
	}
	// Read in the encrypted configuration file.
	dataFile, err := readFile(mc.DirPath + "scoring.dat")
	if err != nil {
		failPrint("Data file (" + fileName + ") not found.")
		os.Exit(1)
	}
	return decryptConfig(dataFile)
}

func printConfig() {
	passPrint("Configuration " + mc.DirPath + "scoring.conf" + " check passed!")
	fmt.Printf("Title: %s (%s)\n", mc.Config.Title, mc.Config.Name)
	fmt.Printf("User: %s\n", mc.Config.User)
	if mc.Config.Remote == "" {
		fmt.Printf("Remote: None (local scoring only)\n")
	} else {
		fmt.Printf("Remote: %s\n", mc.Config.Remote)
	}
	if mc.Config.EndDate == "" {
		fmt.Printf("Valid Until: None (image lasts forever)\n")
	} else {
		fmt.Printf("Valid Until: %s\n", mc.Config.EndDate)
	}
	fmt.Println("Checks:")
	for i, check := range mc.Config.Check {
		fmt.Printf("\tCheck %d (%d points):\n", i+1, check.Points)
		fmt.Printf("\t\tMessage: %s\n", check.Message)
		if check.Pass != nil {
			fmt.Printf("\t\tPassConditions:\n")
			for _, condition := range check.Pass {
				fmt.Printf("\t\t\t%s: %s", condition.Type, condition.Arg1)
				if condition.Arg2 != "" {
					fmt.Printf(", %s\n", condition.Arg2)
				} else {
					fmt.Printf("\n")
				}
			}
		}
		if check.Fail != nil {
			fmt.Printf("\t\tFailConditions:\n")
			for _, condition := range check.Fail {
				fmt.Printf("\t\t\t%s: %s, %s\n", condition.Type, condition.Arg1, condition.Arg2)
			}
		}
	}
}

func passPrint(toPrint string) {
	printer(color.FgGreen, "PASS", toPrint)
}

func failPrint(toPrint string) {
	printer(color.FgRed, "FAIL", toPrint)
}

func warnPrint(toPrint string) {
	printer(color.FgYellow, "WARN", toPrint)
}

func infoPrint(toPrint string) {
	printer(color.FgBlue, "INFO", toPrint)
}

func printer(colorChosen color.Attribute, messageType string, toPrint string) {
	printer := color.New(colorChosen, color.Bold)
	fmt.Printf("[")
	printer.Printf(messageType)
	fmt.Printf("] %s", toPrint)
	if toPrint != "" {
		fmt.Printf("\n")
	}
}

func printerNoNewLine(colorChosen color.Attribute, messageType string, toPrint string) {
	printer := color.New(colorChosen, color.Bold)
	fmt.Printf("[")
	printer.Printf(messageType)
	fmt.Printf("] %s", toPrint)
}

func printerPrompt(toPrint string) {
	printer(color.FgBlue, "?", toPrint)
}

func xor(key string, plaintext string) string {
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = key[i%len(key)] ^ plaintext[i]
	}
	return string(ciphertext)
}

func hexEncode(inputString string) string {
	return hex.EncodeToString([]byte(inputString))
}
