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

package cmd

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

// These hashes are used for XORing the plaintext. Again-- not
// cryptographically genius.
const (
	randomHashOne = "bc1f55699832f80b63da0de463c0f4b030ee3e219371b29cab292b2e439194bc8e517051f0e3dc44a858134e1757fda0445d3a2203ac1e383d20f491105ed5d279b418b85fdc6c2f3003791568af9d85ce7a5b0bf0be90f259e52f089a3ee5682eac0c2b53af5be18fb85c9de8980c2fb32a14fb7fa971881463655fa3dd817d"
	randomHashTwo = "1b29cab292b2e439194bc8e517051f0e3dc44a858134e1757fda0445d3a2203ac1e383d20f491105ed51b29cab292b2e439194bc8e517051f0e3dc44a858134e1757fda0445d3a2203ac1e383d20f491105ed5d279b418b85fdc6c2f3003791568af9d85ce7a5b0bf0bd279b418b85fdc6c2f3003791568af9d85ce7a5b0bf0b"
)

var byteKey = []byte{0x23, 0xf3, 0x24, 0x32, 0x54, 0x76, 0x37, 0x37, 0x86, 0x12, 0x26, 0x07, 0x43, 0x12, 0x26, 0x07, 0x43}

// encryptConfig takes the plainText config and returns an encrypted string
// that should be written to the encrypted scoring data file.
func encryptConfig(plainText string) (string, error) {
	// Generate key by XORing two strings.
	key := xor(randomHashOne, randomHashTwo)

	// Compress the file with zlib.
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)

	// Write zlib compressed data into encryptedFile
	_, err := writer.Write([]byte(plainText))
	if err != nil {
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
		return "", errors.New("error creating zlib reader")
	}
	defer reader.Close()

	// Read into our created buffer.
	dataBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(dataBuffer, reader)
	if err != nil {
		failPrint("error decompressing scoring data")
		return "", errors.New("error decompressing zlib data")
	}

	// Check that decryptedConfig is not empty.
	decryptedConfig := dataBuffer.String()
	if decryptedConfig == "" {
		return "", errors.New("decrypted config is empty")
	}

	return decryptedConfig, err
}

// tossKey is responsible for changing up the byteKey.
func tossKey() []byte {
	// Add your cool byte array manipulations here!
	return byteKey
}

// obfuscateData encodes the configuration when writing to ScoringData.
// This also makes manipulation of data in use harder, since there is
// a very small opportunity for catching plaintext data, and very tough
// to decode the decrypted ScoringData without source code.
func obfuscateData(datum *string) error {
	var err error
	if *datum == "" {
		return errors.New("empty datum given to obfuscateData")
	}
	if *datum, err = encryptConfig(*datum); err == nil {
		*datum = hexEncode(xor(string(tossKey()), *datum))
	} else {
		failPrint("crypto: failed to obufscate datum: " + err.Error())
		return err
	}
	return nil
}

// deobfuscateData decodes configuration data.
func deobfuscateData(datum *string) error {
	var err error
	if *datum == "" {
		return errors.New("empty datum given to deobfuscateData")
	}
	*datum, err = hexDecode(*datum)
	if err != nil {
		println(*datum)
		failPrint("crypto: failed to deobfuscate datum hex: " + err.Error())
		return err
	}
	*datum = xor(string(tossKey()), *datum)
	if *datum, err = decryptConfig(*datum); err != nil {
		failPrint("crypto: failed to deobufscate datum: " + err.Error())
		return err
	}
	return nil
}
