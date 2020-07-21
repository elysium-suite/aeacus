package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func readTeamID(mc *metaConfig, id *imageData) {
	fileContent, err := readFile(mc.DirPath + "TeamID.txt")
	if err != nil {
		failPrint("TeamID.txt does not exist!")
		sendNotification(mc, "TeamID.txt does not exist!")
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Your TeamID files does not exist! Failed to upload scores."
		id.Connection = false
	} else if fileContent == "" {
		failPrint("TeamID.txt is empty!")
		sendNotification(mc, "TeamID.txt is empty!")
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Your TeamID is empty! Failed to upload scores."
		id.Connection = false
	} else {
		mc.TeamID = fileContent
	}
}

// genChallenge generates a crypto challenge for the CSS endpoint
func genChallenge(mc *metaConfig) string {
	// I'm aware this is sus, but right now there's no implemented way to generate a key between aeacus and minos on the fly. Stretch goal
	randomHash1 := "71844fd161e20dc78ce6c985b42611cfb11cf196"
	randomHash2 := "e31ad5a009753ef6da499f961edf0ab3a8eb6e5f"
	chalString := hexEncode(xor(randomHash1, randomHash2))
	if mc.Config.Password != "" {
		hasher := sha256.New()
		hasher.Write([]byte(mc.Config.Password))
		key := hexEncode(string(hasher.Sum(nil)))
		return hexEncode(xor(key, chalString))
	}
	return chalString
}

func genVulns(mc *metaConfig, id *imageData) string {
	var vulnString strings.Builder
	// This could be an unprintable character and be more reliable
	delimiter := "|-|"

	// Vulns achieved
	vulnString.WriteString(fmt.Sprintf("%d%s", len(id.Points), delimiter))
	// Total vulns
	vulnString.WriteString(fmt.Sprintf("%d%s", id.ScoredVulns, delimiter))

	// Build vuln string
	for _, penalty := range id.Penalties {
		vulnString.WriteString(fmt.Sprintf("[PENALTY] %s - %.0f pts", penalty.Message, math.Abs(float64(penalty.Points))))
		vulnString.WriteString(delimiter)
	}

	for _, point := range id.Points {
		vulnString.WriteString(fmt.Sprintf("%s - %d pts", point.Message, point.Points))
		vulnString.WriteString(delimiter)
	}

	if mc.Config.Password != "" {
		if verboseEnabled {
			infoPrint("Encrypting vulnerabilities for score report...")
		}
		return hexEncode(encryptString(mc.Config.Password, vulnString.String()))
	}
	return hexEncode(vulnString.String())
}

func reportScore(mc *metaConfig, id *imageData) {
	resp, err := http.PostForm(mc.Config.Remote+"/scores/css/update",
		url.Values{"team": {mc.TeamID},
			"image": {mc.Config.Name},
			"score": {strconv.Itoa(id.Score)},
			// Challenge string: hash of password
			// XORd with some random crap
			"challenge": {genChallenge(mc)},
			// Vulns: Hex encoded list of vulns
			// encrypted if password exists
			"vulns": {genVulns(mc, id)},
			"id":    {"id"}})
	if err != nil {
		failPrint(err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "OK" {
		failPrint("Failed to upload score! Is your TeamID wrong?")
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Failed to upload score! Please ensure that your Team ID is correct."
		id.Connection = false
		sendNotification(mc, "Failed to upload score! Is your Team ID correct?")
		if mc.Config.Local != "yes" {
			if verboseEnabled {
				warnPrint("Local is not set to \"yes\". Clearing scoring data.")
			}
			clearImageData(id)
		}
	}
}

func checkServer(mc *metaConfig, id *imageData) {

	// Internet check (requisite)
	if verboseEnabled {
		infoPrint("Checking for internet connection...")
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := client.Get("http://clients3.google.com/generate_204")

	if err != nil {
		id.ConnStatus[2] = "red"
		id.ConnStatus[3] = "FAIL"
	} else {
		id.ConnStatus[2] = "green"
		id.ConnStatus[3] = "OK"
	}

	// Scoring engine check
	if verboseEnabled {
		infoPrint("Checking for scoring engine connection...")
	}
	resp, err := client.Get(mc.Config.Remote + "/scores/css/status")

	// todo enforce status/time limit
	// grab body or status message from minos
	// if "DESTROY" due to image elapsed time > time_limit,
	// destroy image

	if err != nil {
		id.ConnStatus[4] = "red"
		id.ConnStatus[5] = "FAIL"
	} else {
		if resp.StatusCode == 200 {
			id.ConnStatus[4] = "green"
			id.ConnStatus[5] = "OK"
		} else {
			id.ConnStatus[4] = "red"
			id.ConnStatus[5] = "ERROR"
		}
	}

	// Overall
	if id.ConnStatus[3] == "FAIL" && id.ConnStatus[5] == "OK" {
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Server connection good but no Internet. Assuming you're on an isolated LAN."
		id.Connection = true
	} else if id.ConnStatus[5] == "FAIL" {
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Failure! Can't access remote scoring server."
		failPrint("Can't access remote scoring server!")
		sendNotification(mc, "Score upload failure! Unable to access remote server.")
		id.Connection = false
	} else if id.ConnStatus[4] == "ERROR" {
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Score upload failure. Can't send scores to remote server."
		failPrint("Remote server returned an error for its status!")
		sendNotification(mc, "Score upload failure! Remote server returned an error.")
		id.Connection = false
	} else {
		id.ConnStatus[0] = "green"
		id.ConnStatus[1] = "OK"
		id.Connection = true
	}

	readTeamID(mc, id)
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
