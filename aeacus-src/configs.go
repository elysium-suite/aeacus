package main

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
    "io"

    // crypto magic
	"io/ioutil"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"

	"github.com/fatih/color"
	"github.com/BurntSushi/toml"
)

func parseConfig(mc *metaConfig, configContent string) {
	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

////////////////////
// ENCRYPT CONFIG //
////////////////////

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

    // "static" hash #1
    hashOne := "7z7551253a53s0f974e3d03d0cf839e7ccfc!879"
	hashOneContent, err := readFile("/bin/bash")
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(hashOneContent))
    	hashOne = hex.EncodeToString(hasher.Sum(nil))
    }

    // "static" hash #2
    hashTwo := "3384b1be7ac2a~9ahc8b4488d4cc2edb5ag497fz"
	hashTwoContent, err := readFile("/usr/lib/apt/apt-helper")
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(hashTwoContent))
    	hashTwo = hex.EncodeToString(hasher.Sum(nil))
    }

    // formulate key with hashes + modified day of config
    key := xor(hashOne, hashTwo)
    info, err = os.Stat(mc.ConfigName)
    if err != nil {
        failPrint("Crypto magic can not ensue! No configuration file found.")
        os.Exit(1)
    }
    modifiedTime := info.ModTime().Format("01/02/2006")
    modifiedTimeHash := "1230-8123nasklnaegklnjwh0-91uiowasfml;3tr23"
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(modifiedTime))
    	modifiedTimeHash = hex.EncodeToString(hasher.Sum(nil))
    }
    key = xor(modifiedTimeHash, key)

    // swap some bytes just 4 fun
    // TODO
    //key = append(key, key[7])
    //key = append(key, key[10])

    // zlib compress
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write(configBuffer)
	writer.Close()

    // apply xor key
    xordFile := xor(key, encryptedFile.String())

	// aes with reversed byte string or something
    // TODO

	if mc.Cli.Bool("v") {
		infoPrint("Writing data to " + mc.DataName + "...")
	}
	err = ioutil.WriteFile(mc.DataName, []byte(xordFile), info.Mode())
}

////////////////////
// DECRYPT CONFIG //
////////////////////

func readData(mc *metaConfig) string {
	if mc.Cli.Bool("v") {
		infoPrint("Decrypting data from " + mc.DataName)
	}

	dataFile, err := readFile(mc.DataName)
	if err != nil {
        failPrint("Data file not found.")
        os.Exit(1)
	}

    // "static" hash #1
    hashOne := "7z7551253a53s0f974e3d03d0cf839e7ccfc!879"
	hashOneContent, err := readFile("/bin/bash")
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(hashOneContent))
    	hashOne = hex.EncodeToString(hasher.Sum(nil))
    }

    // "static" hash #2
    hashTwo := "3384b1be7ac2a~9ahc8b4488d4cc2edb5ag497fz"
	hashTwoContent, err := readFile("/usr/lib/apt/apt-helper")
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(hashTwoContent))
    	hashTwo = hex.EncodeToString(hasher.Sum(nil))
    }

    // formulate key with hashes + modified day of config
    key := xor(hashOne, hashTwo)
    info, err := os.Stat(mc.DataName)
    if err != nil {
        failPrint("Oops, you yoinked scoring.dat? Uncool.")
        os.Exit(1)
    }
    modifiedTime := info.ModTime().Format("01/02/2006")
    modifiedTimeHash := "1230-8123nasklnaegklnjwh0-91uiowasfml;3tr23"
    if err == nil {
    	hasher := sha1.New()
    	hasher.Write([]byte(modifiedTime))
    	modifiedTimeHash = hex.EncodeToString(hasher.Sum(nil))
    }
    key = xor(modifiedTimeHash, key)

    // swap some bytes just 4 fun
    // TODO
    //key = append(key, key[7])
    //key = append(key, key[10])

    // undo aes

    // apply xor key
    dataFile = xor(key, dataFile)

    // zlib decompress
	reader, err := zlib.NewReader(bytes.NewReader([]byte(dataFile)))
    if err != nil {
        failPrint("Error decrypting scoring.dat. You naughty little competitor. Commencing self destruct...")
        // lol jk... for now
        os.Exit(1)
    }
	defer reader.Close()


    dataBuffer := bytes.NewBuffer(nil)
    io.Copy(dataBuffer, reader)

    return string(dataBuffer.Bytes())
}

/////////////////////////////
// CRYPTOGRAPHIC FUNCTIONS //
/////////////////////////////

func xor (key string, plaintext string) string {
    ciphertext := make([]byte, len(plaintext))
    for i := 0; i < len(plaintext); i++ {
        ciphertext[i] = key[i % len(key)] ^ plaintext[i]
    }
    return string(ciphertext)
}

//////////////////////
// HELPER FUNCTIONS //
//////////////////////

func printConfig(mc *metaConfig) {
	passPrint("Configuration " + mc.ConfigName + " check passed!")
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
