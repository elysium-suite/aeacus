package main

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
    "io"
	"io/ioutil"
	"compress/zlib"
	"github.com/fatih/color"
	"github.com/BurntSushi/toml"
)

func parseConfig(mc *metaConfig, configContent string) {
	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeConfig(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Reading configuration from " + mc.ConfigName + "...")
	}

	configFile, err := os.Open(mc.ConfigName)
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

	encryptedBuffer := encryptConfig(configBuffer)

	if mc.Cli.Bool("v") {
		infoPrint("Writing data to " + mc.DataName + "...")
	}
	err = ioutil.WriteFile(mc.DataName, encryptedBuffer.Bytes(), info.Mode())
}

func readData(mc *metaConfig) string {
	if mc.Cli.Bool("v") {
		infoPrint("Decrypting data from " + mc.DataName)
	}
	dataFile, err := os.Open(mc.DataName)
	if err != nil {
        failPrint("Data file not found.")
        os.Exit(1)
	}
	defer dataFile.Close()
	return decryptData(dataFile)
}

/////////////////////////////
// CRYPTOGRAPHIC FUNCTIONS //
/////////////////////////////

func encryptConfig(configFile []byte) bytes.Buffer {
	// xor with defined byte string
	// zlib
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write(configFile)
	writer.Close()
	// aes with reversed byte string or something

	return encryptedFile
}

func decryptData(dataFile *os.File) string {
	// aes with reversed byte string

	reader, err := zlib.NewReader(dataFile)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	defer reader.Close()

    dataBuffer := bytes.NewBuffer(nil)
    io.Copy(dataBuffer, reader)

	// xor with defined byte string

	return string(dataBuffer.Bytes())
}

//////////////////////
// HELPER FUNCTIONS //
//////////////////////

func printConfig(mc *metaConfig) {
	passPrint("Configuration " + mc.ConfigName + " check passed!")
	fmt.Printf("Title: %s (%s)\n", mc.Config.Title, mc.Config.Name)
	fmt.Printf("User: %s\n", mc.Config.User)
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
