package main

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"bytes"
    "strings"
	"io/ioutil"

	// crypto magic
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"

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

	configFile, err := os.Open(mc.DirPath + "scoring.conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	info, _ := configFile.Stat()
	var size int64 = info.Size()
	configBuffer := make([]byte, size)
	buffer := bufio.NewReader(configFile)
	_, err = buffer.Read(configBuffer)

	if mc.Cli.Bool("v") {
		infoPrint("Encrypting configuration...")
	}

	info, err = os.Stat(mc.DirPath + "scoring.conf")
	if err != nil {
		failPrint("Crypto magic can not occur! No configuration file found.")
		os.Exit(1)
	}

    // Additionally, XOR it with ModTime
	modifiedTime := info.ModTime().Format("01/02/2006")
	modifiedTimeHash := "3208c653a58297997ae22a3ea21be68fb2f4d06"
	if err == nil {
		hasher := sha1.New()
		hasher.Write([]byte(modifiedTime))
		modifiedTimeHash = hex.EncodeToString(hasher.Sum(nil))
	}
	key := xor(modifiedTimeHash, getXORKey())

	// zlib compress
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write(configBuffer)
	writer.Close()

	// apply xor key
	xordFile := xor(key, encryptedFile.String())

	if mc.Cli.Bool("v") {
		infoPrint("Writing data to " + mc.DirPath + "...")
	}
	writeFile(mc.DirPath+"scoring.dat", xordFile)
}

////////////////////
// DECRYPT CONFIG //
////////////////////

func readData(mc *metaConfig) string {
	if mc.Cli.Bool("v") {
		infoPrint("Decrypting data from " + mc.DirPath + "scoring.dat...")
	}

	dataFile, err := readFile(mc.DirPath + "scoring.dat")
	if err != nil {
		failPrint("Data file not found.")
		os.Exit(1)
	}

    // Apply jank ModTime hash
	info, err := os.Stat(mc.DirPath + "scoring.dat")
	if err != nil {
		failPrint("Oops, you yoinked scoring.dat? Uncool.")
		os.Exit(1)
	}
	modifiedTime := info.ModTime().Format("01/02/2006")
	modifiedTimeHash := "3208c653a58297997ae22a3ea21be68fb2f4d06"
	if err == nil {
		hasher := sha1.New()
		hasher.Write([]byte(modifiedTime))
		modifiedTimeHash = hex.EncodeToString(hasher.Sum(nil))
	}
	key := xor(modifiedTimeHash, getXORKey())

	// decrypt with xor key
	dataFile = xor(key, dataFile)

	// zlib decompress
	reader, err := zlib.NewReader(bytes.NewReader([]byte(dataFile)))
	if err != nil {
		failPrint("Error decrypting scoring.dat. You naughty little competitor. Commencing self destruct...")
        destroyImage()
		os.Exit(1)
	}
	defer reader.Close()

	dataBuffer := bytes.NewBuffer(nil)
	io.Copy(dataBuffer, reader)

	return string(dataBuffer.Bytes())
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

func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	return string(fileContent), err
}

func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println(err)
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
