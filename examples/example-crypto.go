package main

import (
	"io"
    "os"
	"fmt"

    "bufio"
    "bytes"

	"crypto/aes"
	"crypto/rand"
    "crypto/cipher"
	"crypto/sha256"

    "compress/zlib"
)

func writeCryptoConfig(mc *metaConfig) string {

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

	randomHash := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
    key := xor(randomHash, "ThatsASecureStringLol")

    // zlib compress
    var encryptedFile bytes.Buffer
    writer := zlib.NewWriter(&encryptedFile)
    writer.Write(configBuffer)
    writer.Close()

    // apply xor key
    return xor(key, encryptedFile.String())

}

func readCryptoConfig(mc *metaConfig) string {
    dataFile, err := readFile(mc.DirPath + "scoring.dat")
	if err != nil {
		failPrint("Data file not found.")
		os.Exit(1)
	}

	randomHash := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
    key := xor(randomHash, "ThatsASecureStringLol")

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

func encryptString(password string, plaintext string) string {

	hasher := sha256.New()
	hasher.Write([]byte(password))
    key := hasher.Sum(nil)

    // Pad plaintext to be a 16-byte block
    paddingArray := make([]byte, (aes.BlockSize - len(plaintext)%aes.BlockSize))
    for char := range paddingArray {
        paddingArray[char] = 0x20
    }
    plaintext = plaintext + string(paddingArray)
	if len(plaintext)%aes.BlockSize != 0 {
		panic("Plaintext is not a multiple of block size!")
	}

    // Create cipher block with key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

    // Generate nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

    // Create NewGCM cipher
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

    // Encrypt and seal plaintext
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
    ciphertext = []byte(fmt.Sprintf("%s%s", nonce, ciphertext))

    return string(ciphertext)
}

func decryptString(password string, ciphertext string) string {

	hasher := sha256.New()
	hasher.Write([]byte(password))
    key := hasher.Sum(nil)

    iv := []byte(ciphertext[:12])
    ciphertext = ciphertext[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		failPrint(err.Error())
        return ""
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		failPrint(err.Error())
        return ""
	}

	plaintext, err := aesgcm.Open(nil, iv, []byte(ciphertext), nil)
	if err != nil {
		failPrint(err.Error())
        return ""
	}

    return string(plaintext)

}
