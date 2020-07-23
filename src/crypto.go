// crypto.go is an example to provides basic cryptographical functions
// for aeacus.
//
// This file is not a shining example for cryptographically secure operations.
//
// Practically, it is more important that your implemented solution is
// different than the example, to make reverse engineering much more difficult.
//
// You could radically change the crypto.go file each time you release an
// image, which would make things very difficult for a would-be hacker.
//
// At the very least, edit some strings. Add some ciphers and operations if
// you're feeling spicy.

package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// This key is used as a backup for remote encryption when there is no password
// specified in the configuration.
//
// This must be the same value specified in Minos (or any other reporting
// endpoint) as a backup password.
var remoteBackupKey = "ThisIsAReallyCoolAndSecureKeyLol"

// These hashes are used for XORing the plaintext. Again-- not
// cryptographically genius.
var randomHashOne = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
var randomHashTwo = "NowThatsWhatICallARandomString"

// writeCryptoConfig takes the metaConfig (context) and writes to the hardcoded
// file `scoring.dat`, in the DirPath specified in the metaConfig.
// writeCryptoConfig is used to create the encrypted `scoring.dat` from the
// plaintext configuration `scoring.conf`.
func writeCryptoConfig(mc *metaConfig) string {

	// Open the hardcoded file path to the plaintext configuration.
	configFile, err := os.Open(mc.DirPath + "scoring.conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Read the file into a buffer.
	info, _ := configFile.Stat()
	var size int64 = info.Size()
	configBuffer := make([]byte, size)
	buffer := bufio.NewReader(configFile)
	_, err = buffer.Read(configBuffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verboseEnabled {
		infoPrint("Encrypting configuration...")
	}

	// Generate key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Compress the file with zlib.
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write(configBuffer)
	writer.Close()

	// XOR the encrypted file with our key.
	return xor(key, encryptedFile.String())

}

// readCryptoConfig is used to decrypt the `scoring.dat` file, which contains
// the configuration for aeacus.
func readCryptoConfig(mc *metaConfig) string {

	// Read in the encrypted configuration file.
	dataFile, err := readFile(mc.DirPath + "scoring.dat")
	if err != nil {
		failPrint("Data file not found.")
		os.Exit(1)
	}

	// Create our key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Apply the XOR key to decrypt the zlib-compressed data.
	//
	// XOR is special in that when you apply it twice, you get the original data
	// as long as the key was the same.
	dataFile = xor(key, dataFile)

	// Decompress zlib data.
	reader, err := zlib.NewReader(bytes.NewReader([]byte(dataFile)))
	if err != nil {
		failPrint("Error decrypting scoring.dat. You naughty little competitor. Commencing self destruct...")
		destroyImage(mc)
		os.Exit(1)
	}
	defer reader.Close()
	dataBuffer := bytes.NewBuffer(nil)
	io.Copy(dataBuffer, reader)

	return string(dataBuffer.Bytes())
}
