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
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

// These hashes are used for XORing the plaintext. Again-- not
// cryptographically genius.
const (
	randomHashOne = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	randomHashTwo = "NowThatsWhatICallARandomString"
)

// encryptConfig takes the plainText config and returns an encrypted string
// that should be written to the encrypted scoring data file.
func encryptConfig(plainText string) string {
	if verboseEnabled {
		infoPrint("Encrypting configuration...")
	}

	// Generate key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Compress the file with zlib.
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write([]byte(plainText))
	writer.Close()

	// XOR the encrypted file with our key.
	return xor(key, encryptedFile.String())
}

// decryptConfig is used to decrypt the scoring data file.
func decryptConfig(cipherText string) string {
	// Create our key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Apply the XOR key to decrypt the zlib-compressed data.
	//
	// XOR is special in that when you apply it twice, you get the original data
	// as long as the key was the same.
	cipherText = xor(key, cipherText)

	// Decompress zlib data.
	reader, err := zlib.NewReader(bytes.NewReader([]byte(cipherText)))
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
