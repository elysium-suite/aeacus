package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

func parseConfig(mc *metaConfig, configContent string) {
	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mc.Config.Local = strings.ToLower(mc.Config.Local)
}

////////////////////
// ENCRYPT CONFIG //
////////////////////

func writeConfig(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Reading configuration from " + mc.DirPath + "scoring.conf" + "...")
	}
	encryptedConfig := writeCryptoConfig(mc)
	if mc.Cli.Bool("v") {
		infoPrint("Writing data to " + mc.DirPath + "...")
	}
	writeFile(mc.DirPath+"scoring.dat", encryptedConfig)
}

////////////////////
// DECRYPT CONFIG //
////////////////////

func readData(mc *metaConfig) string {
	if mc.Cli.Bool("v") {
		infoPrint("Decrypting data from " + mc.DirPath + "scoring.dat...")
	}
	return readCryptoConfig(mc)
}

//////////////////////
// HELPER FUNCTIONS //
//////////////////////

func printConfig(mc *metaConfig) {
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
