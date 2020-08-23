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
	randomHashOne = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	randomHashTwo = "NowThatsWhatICallARandomString"
)

var byteKey = []byte{0x53, 0xf7, 0xb1, 0xcd, 0x26, 0x7a, 0x6f, 0x9a, 0xa5, 0x61, 0xb0, 0x97, 0x21}

// encryptConfig takes the plainText config and returns an encrypted string
// that should be written to the encrypted scoring data file.
func encryptConfig(plainText string) (string, error) {
	debugPrint("Encrypting data...")
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
	debugPrint("Decrypting data...")
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
	decryptedConfig := string(dataBuffer.Bytes())
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
func obfuscateData(datum *string) {
	var err error
	if *datum == "" {
		return
	}
	if *datum, err = encryptConfig(*datum); err == nil {
		*datum = hexEncode(xor(string(tossKey()), *datum))
	} else {
		failPrint("crypto: failed to obufscate datum: " + err.Error())
	}
}

// deobfuscateData decodes configuration data.
func deobfuscateData(datum *string) {
	var err error
	if *datum == "" {
		return
	}
	*datum, err = hexDecode(*datum)
	if err != nil {
		println(*datum)
		failPrint("crypto: failed to deobfuscate datum hex: " + err.Error())
		return
	}
	*datum = xor(string(tossKey()), *datum)
	if *datum, err = decryptConfig(*datum); err != nil {
		failPrint("crypto: failed to deobufscate datum: " + err.Error())
	}
}
