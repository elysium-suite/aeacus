package cmd

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

// parseConfig takes the config content as a string and attempts to parse it
// into the mc.Config struct based on the TOML spec.
func parseConfig(configContent string) {
	if configContent == "" {
		failPrint("Configuration is empty!")
		os.Exit(1)
	}

	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		failPrint("error decoding TOML: " + err.Error())
		os.Exit(1)
	}

	// If there's no remote, local must be enabled.
	if mc.Config.Remote == "" {
		mc.Config.Local = true
	}

	if mc.Config.Remote != "" {
		if mc.Config.Name == "" {
			failPrint("Need image name in config if remote is enabled.")
			os.Exit(1)
		}
	}
}

// WriteConfig reads the plaintext configuration from sourceFile, and writes
// the encrypted configuration into the destFile name passed.
func WriteConfig(sourceFile, destFile string) {
	infoPrint("Reading configuration from " + mc.DirPath + sourceFile + "...")

	configFile, err := readFile(mc.DirPath + sourceFile)
	if err != nil {
		failPrint("Can't open scoring configuration file (" + sourceFile + "): " + err.Error())
		os.Exit(1)
	}
	parseConfig(configFile)
	configFile = ""

	obfuscateConfig()
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(mc.Config); err != nil {
		failPrint(err.Error())
		os.Exit(1)
		return
	}

	encryptedConfig, err := encryptConfig(buf.String())
	if err != nil {
		failPrint("Encrypting config failed: " + err.Error())
		os.Exit(1)
	} else if verboseEnabled {
		infoPrint("Writing data to " + mc.DirPath + "...")
	}

	writeFile(mc.DirPath+destFile, encryptedConfig)
}

// readData is a wrapper around decryptData, taking the scoring data fileName,
// and reading its content. It returns the decrypt config.
func readData(fileName string) (string, error) {
	// Read in the encrypted configuration filei
	dataFile, err := readFile(mc.DirPath + ScoringData)
	if err != nil {
		return "", err
	} else if dataFile == "" {
		return "", errors.New("Scoring data is empty!")
	}
	decryptedConfig, err := decryptConfig(dataFile)
	if err != nil {
		return "", err
	}
	return decryptedConfig, nil
}

// printConfig offers a printed representation of the config, as parsed
// by readData and parseConfig.
func printConfig() {
	passPrint("Configuration " + mc.DirPath + ScoringConf + " check passed!")
	fmt.Println("Title:", mc.Config.Title)
	fmt.Println("Name:", mc.Config.Name)
	fmt.Println("OS:", mc.Config.OS)
	fmt.Println("User:", mc.Config.User)
	fmt.Println("Remote:", mc.Config.Remote)
	fmt.Println("Local:", mc.Config.Local)
	fmt.Println("EndDate:", mc.Config.EndDate)
	fmt.Println("NoDestroy:", mc.Config.NoDestroy)
	fmt.Println("DisableShell:", mc.Config.DisableShell)
	fmt.Println("Checks:")
	for i, check := range mc.Config.Check {
		fmt.Printf("\tCheck %d (%d points):\n", i+1, check.Points)
		fmt.Println("\t\tMessage:", check.Message)
		if check.Pass != nil {
			fmt.Println("\t\tPassConditions:")
			for _, condition := range check.Pass {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
		if check.PassOverride != nil {
			fmt.Println("\t\tPassOverrideConditions:")
			for _, condition := range check.PassOverride {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
		if check.Fail != nil {
			fmt.Println("\t\tFailConditions:")
			for _, condition := range check.Fail {
				fmt.Printf("\t\t\t%s: %s %s %s %s\n", condition.Type, condition.Arg1, condition.Arg2, condition.Arg3, condition.Arg4)
			}
		}
	}
}

func obfuscateConfig() {
	infoPrint("Obfuscating configuration...")
	obfuscateData(&mc.Config.Password)
	for i, check := range mc.Config.Check {
		obfuscateData(&mc.Config.Check[i].Message)
		if check.Pass != nil {
			for x := range check.Pass {
				obfuscateData(&mc.Config.Check[i].Pass[x].Type)
				obfuscateData(&mc.Config.Check[i].Pass[x].Arg1)
				obfuscateData(&mc.Config.Check[i].Pass[x].Arg2)
				obfuscateData(&mc.Config.Check[i].Pass[x].Arg3)
				obfuscateData(&mc.Config.Check[i].Pass[x].Arg4)
			}
		}
		if check.PassOverride != nil {
			for x := range check.Pass {
				obfuscateData(&mc.Config.Check[i].PassOverride[x].Type)
				obfuscateData(&mc.Config.Check[i].PassOverride[x].Arg1)
				obfuscateData(&mc.Config.Check[i].PassOverride[x].Arg2)
				obfuscateData(&mc.Config.Check[i].PassOverride[x].Arg3)
				obfuscateData(&mc.Config.Check[i].PassOverride[x].Arg4)
			}
		}
		if check.Fail != nil {
			for x := range check.Fail {
				obfuscateData(&mc.Config.Check[i].Fail[x].Type)
				obfuscateData(&mc.Config.Check[i].Fail[x].Arg1)
				obfuscateData(&mc.Config.Check[i].Fail[x].Arg2)
				obfuscateData(&mc.Config.Check[i].Fail[x].Arg3)
				obfuscateData(&mc.Config.Check[i].Fail[x].Arg4)
			}
		}
	}
}

// confirmPrint will prompt the user with the given toPrint string, and
// exit the program if N or n is input.
func ConfirmPrint(toPrint string) {
	printer(color.FgYellow, "CONF", "")
	fmt.Print(toPrint + " [Y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(strings.TrimSpace(resp)) == "n" {
		os.Exit(1)
	}
}

func passPrint(toPrint string) {
	printStr := printer(color.FgGreen, "PASS", toPrint)
	if verboseEnabled {
		fmt.Printf(printStr)
	}
}

func failPrint(toPrint string) {
	fmt.Printf(printer(color.FgRed, "FAIL", toPrint))
}

func warnPrint(toPrint string) {
	fmt.Printf(printer(color.FgYellow, "WARN", toPrint))
}

func debugPrint(toPrint string) {
	printStr := printer(color.FgMagenta, "DEBUG", toPrint)
	if debugEnabled {
		fmt.Printf(printStr)
	}
}

func infoPrint(toPrint string) {
	printStr := printer(color.FgCyan, "INFO", toPrint)
	if verboseEnabled {
		fmt.Printf(printStr)
	}
}

func printer(colorChosen color.Attribute, messageType, toPrint string) string {
	printer := color.New(colorChosen, color.Bold)
	printStr := fmt.Sprintf("[")
	printStr += printer.Sprintf(messageType)
	printStr += fmt.Sprintf("] %s", toPrint)
	if toPrint != "" {
		printStr += fmt.Sprintf("\n")
	}
	return printStr
}

func xor(key, plaintext string) string {
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = key[i%len(key)] ^ plaintext[i]
	}
	return string(ciphertext)
}

func hexEncode(inputString string) string {
	return hex.EncodeToString([]byte(inputString))
}

func hexDecode(inputString string) (string, error) {
	result, err := hex.DecodeString(inputString)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
