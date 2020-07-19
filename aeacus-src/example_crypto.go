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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

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

	if mc.Cli.Bool("v") {
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

// encryptString takes a password and a plaintext and returns an encrypted byte
// sequence (as a string). It uses AES-GCM with a 12-byte IV (as is
// recommended). The IV is prefixed to the string.
//
// This function is used in aeacus to encrypt reported vulnerability data to
// the remote scoring endpoint (ex. minos).
func encryptString(password, plaintext string) string {

	// Create a sha256sum hash of the password provided.
	hasher := sha256.New()
	hasher.Write([]byte(password))
	key := hasher.Sum(nil)

	// Pad plaintext to be a 16-byte block.
	paddingArray := make([]byte, (aes.BlockSize - len(plaintext)%aes.BlockSize))
	for char := range paddingArray {
		paddingArray[char] = 0x20 // Padding with space character.
	}
	plaintext = plaintext + string(paddingArray)
	if len(plaintext)%aes.BlockSize != 0 {
		panic("Plaintext is not a multiple of block size!")
	}

	// Create cipher block with key.
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Generate nonce.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Create NewGCM cipher.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Encrypt and seal plaintext.
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	ciphertext = []byte(fmt.Sprintf("%s%s", nonce, ciphertext))

	return string(ciphertext)
}

// decryptString takes a password and a ciphertext and returns a decrypted
// byte sequence (as a string). The function uses typical AES-GCM.
func decryptString(password, ciphertext string) string {

	// Create a sha256sum hash of the password provided.
	hasher := sha256.New()
	hasher.Write([]byte(password))
	key := hasher.Sum(nil)

	// Grab the IV from the first 12 bytes of the file.
	iv := []byte(ciphertext[:12])
	ciphertext = ciphertext[12:]

	// Create the AES block object.
	block, err := aes.NewCipher(key)
	if err != nil {
		failPrint(err.Error())
		return ""
	}

	// Create the AES-GCM cipher with the generated block.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		failPrint(err.Error())
		return ""
	}

	// Decrypt (and check validity, since it's GCM) of ciphertext.
	plaintext, err := aesgcm.Open(nil, iv, []byte(ciphertext), nil)
	if err != nil {
		failPrint(err.Error())
		return ""
	}

	return string(plaintext)
}
