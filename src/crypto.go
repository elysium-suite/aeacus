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
	"errors"
	"io"
)

// These hashes are used for XORing the plaintext. Again-- not
// cryptographically genius.
const (
	randomHashOne = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	randomHashTwo = "NowThatsWhatICallARandomString"
)

// encryptConfig takes the plainText config and returns an encrypted string
// that should be written to the encrypted scoring data file.
func encryptConfig(plainText string) (string, error) {
	if verboseEnabled {
		infoPrint("Encrypting configuration...")
	}

	// Generate key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Compress the file with zlib.
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)

	// Write zlib compressed data into encryptedFile
	_, err := writer.Write([]byte(plainText))
	if err != nil {
		if debugEnabled {
			failPrint("Unable to zlib compress scoring data: " + err.Error())
		}
		return "", err
	}
	writer.Close()

	// XOR the encrypted file with our key.
	return xor(key, encryptedFile.String()), err
}

// decryptConfig is used to decrypt the scoring data file.
func decryptConfig(cipherText string) (string, error) {
	// Create our key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Apply the XOR key to decrypt the zlib-compressed data.
	cipherText = xor(key, cipherText)

	// Create the zlib reader.
	reader, err := zlib.NewReader(bytes.NewReader([]byte(cipherText)))
	if err != nil {
		if debugEnabled {
			failPrint("Error creating archive reader for scoring data.")
		}
		return "", errors.New("Error creating zLib reader")
	}
	defer reader.Close()

	// Read into our created buffer.
	dataBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(dataBuffer, reader)
	if err != nil {
		if debugEnabled {
			failPrint("Error decompressing scoring data.")
		}
		return "", errors.New("Error decompressing zlib data.")
	}

	// Check that decryptedConfig is not empty.
	decryptedConfig := string(dataBuffer.Bytes())
	if decryptedConfig == "" {
		if debugEnabled {
			failPrint("Scoring data is empty!")
		}
		return "", errors.New("Decrypted config is empty!")
	}

	return decryptedConfig, err
}
