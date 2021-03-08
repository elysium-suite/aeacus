package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
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

var delimiter = "|-S#-|"

const (
	FAIL  = "FAIL"
	GREEN = "green"
	RED   = "red"
)

func readTeamID() {
	fileContent, err := readFile(mc.DirPath + "TeamID.txt")
	fileContent = strings.TrimSpace(fileContent)
	if err != nil {
		failPrint("TeamID.txt does not exist!")
		sendNotification("TeamID.txt does not exist!")
		mc.Conn.OverallColor = RED
		mc.Conn.OverallStatus = "Your TeamID files does not exist! Failed to upload scores."
		mc.Connection = false
	} else if fileContent == "" {
		failPrint("TeamID.txt is empty!")
		sendNotification("TeamID.txt is empty!")
		mc.Conn.OverallStatus = RED
		mc.Conn.OverallStatus = "Your TeamID is empty! Failed to upload scores."
		mc.Connection = false
	} else {
		mc.TeamID = fileContent
	}
}

// genChallenge generates a crypto challenge for the CSS endpoint
func genChallenge() (string, error) {
	// Should actually use this for something
	randomHash1 := "71844fd161e20dc78ce6c985b42611cfb11cf196"
	randomHash2 := "e31ad5a009753ef6da499f961edf0ab3a8eb6e5f"
	chalString := hexEncode(xor(randomHash1, randomHash2))
	hasher := sha256.New()
	_, err := hasher.Write([]byte(mc.Config.Password))
	if err != nil {
		return "", err
	}
	key := hexEncode(string(hasher.Sum(nil)))
	return hexEncode(xor(key, chalString)), nil
}

func writeString(stringToWrite *strings.Builder, key, value string) {
	stringToWrite.WriteString(key)
	stringToWrite.WriteString(delimiter)
	stringToWrite.WriteString(value)
	stringToWrite.WriteString(delimiter)
}

func genUpdate() (string, error) {
	var update strings.Builder
	// Write values for score update
	writeString(&update, "team", mc.TeamID)
	writeString(&update, "image", mc.Config.Name)
	writeString(&update, "score", strconv.Itoa(mc.Image.Score))
	chall, err := genChallenge()
	if err != nil {
		return "", err
	}
	writeString(&update, "challenge", chall)
	writeString(&update, "vulns", genVulns())
	writeString(&update, "time", strconv.Itoa(int(time.Now().Unix())))
	infoPrint("Encrypting score update...")
	if err := deobfuscateData(&mc.Config.Password); err != nil {
		errorPrint(err)
	}
	finishedUpdate := hexEncode(encryptString(mc.Config.Password, update.String()))
	if err := obfuscateData(&mc.Config.Password); err != nil {
		errorPrint(err)
	}
	return finishedUpdate, nil
}

func genVulns() string {
	var vulnString strings.Builder

	// Vulns achieved
	vulnString.WriteString(fmt.Sprintf("%d%s", len(mc.Image.Points), delimiter))
	// Total vulns
	vulnString.WriteString(fmt.Sprintf("%d%s", mc.Image.ScoredVulns, delimiter))

	// Build vuln string
	for _, penalty := range mc.Image.Penalties {
		if err := deobfuscateData(&penalty.Message); err != nil {
			errorPrint(err)
		}
		vulnString.WriteString(fmt.Sprintf("%s - N%.0f pts", penalty.Message, math.Abs(float64(penalty.Points))))
		if err := obfuscateData(&penalty.Message); err != nil {
			errorPrint(err)
		}
		vulnString.WriteString(delimiter)
	}

	for _, point := range mc.Image.Points {
		if err := deobfuscateData(&point.Message); err != nil {
			errorPrint(err)
		}
		vulnString.WriteString(fmt.Sprintf("%s - %d pts", point.Message, point.Points))
		if err := obfuscateData(&point.Message); err != nil {
			errorPrint(err)
		}
		vulnString.WriteString(delimiter)
	}

	infoPrint("Encrypting vulnerabilities...")

	deobfuscateData(&mc.Config.Password)
	finishedVulns := hexEncode(encryptString(mc.Config.Password, vulnString.String()))
	obfuscateData(&mc.Config.Password)
	return finishedVulns
}

func reportScore() error {
	update, err := genUpdate()
	if err != nil {
		failPrint(err.Error())
		return err
	}
	resp, err := http.PostForm(mc.Config.Remote+"/update",
		url.Values{"update": {update}})
	if err != nil {
		failPrint(err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		mc.Conn.OverallColor = RED
		mc.Conn.OverallStatus = "Failed to upload score! Please ensure that your Team ID is correct."
		mc.Connection = false
		failPrint("Failed to upload score!")
		sendNotification("Failed to upload score!")
		return errors.New("Non-200 response from remote scoring endpoint")
	}
	return nil
}

func checkServer() {
	// Internet check (requisite)
	infoPrint("Checking for internet connection...")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := client.Get("http://example.org")

	if err != nil {
		mc.Conn.NetColor = RED
		mc.Conn.NetStatus = FAIL
	} else {
		mc.Conn.NetColor = GREEN
		mc.Conn.NetStatus = "OK"
	}

	// Scoring engine check
	infoPrint("Checking for scoring engine connection...")
	resp, err := client.Get(mc.Config.Remote + "/status/" + mc.TeamID + "/" + mc.Config.Name)
	// todo enforce status/time limit
	// grab body or status message from minos
	// if "DESTROY" due to image elapsed time > time_limit,
	// destroy image

	if err != nil {
		mc.Conn.ServerColor = RED
		mc.Conn.ServerStatus = FAIL
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			failPrint("Error reading Status body.")
			mc.Conn.ServerColor = RED
			mc.Conn.ServerStatus = FAIL
		} else {
			handleStatus(string(body))
			if resp.StatusCode == 200 {
				mc.Conn.ServerColor = GREEN
				mc.Conn.ServerStatus = "OK"
			} else {
				mc.Conn.ServerColor = RED
				mc.Conn.ServerStatus = "ERROR"
			}
		}
	}

	// Overall
	if mc.Conn.NetStatus == FAIL && mc.Conn.ServerStatus == "OK" {
		timeStart = time.Now()
		mc.Conn.OverallColor = "goldenrod"
		mc.Conn.OverallStatus = "Server connection good but no Internet. Assuming you're on an isolated LAN."
		mc.Connection = true
	} else if mc.Conn.ServerStatus == FAIL {
		timeStart = time.Now()
		mc.Conn.OverallColor = RED
		mc.Conn.OverallStatus = "Failure! Can't access remote scoring server."
		failPrint("Can't access remote scoring server!")
		sendNotification("Score upload failure! Unable to access remote server.")
		mc.Connection = false
	} else if mc.Conn.ServerStatus == "ERROR" {
		timeWithoutID = time.Since(timeStart)
		if !mc.Config.NoDestroy && timeWithoutID > withoutIDThreshold {
			failPrint("Destroying the image! Too long without inputting valid ID.")
			// destroyImage()
		}
		mc.Conn.OverallColor = RED
		mc.Conn.OverallStatus = "Scoring engine rejected your TeamID!"
		failPrint("Remote server returned an error for its status! Your ID is probably wrong.")
		sendNotification("Status check failed, TeamID incorrect!")
		mc.Connection = false
	} else {
		timeStart = time.Now()
		mc.Conn.OverallColor = GREEN
		mc.Conn.OverallStatus = "OK"
		mc.Connection = true
	}
}

func handleStatus(status string) {
	var statusStruct statusRes
	if err := json.Unmarshal([]byte(status), &statusStruct); err != nil {
		failPrint("Failed to parse JSON response (status): " + err.Error())
	}

	switch statusStruct.Status {
	case "DIE":
		failPrint("Destroying image! Server has told me to die.")
		// destroyImage()
	case "GIMMESHELL":
		if !mc.Config.DisableShell && !mc.ShellActive {
			go connectWs()
		}
	}
}

// encryptString takes a password and a plaintext and returns an encrypted byte
// sequence (as a string). It uses AES-GCM with a 12-byte IV (as is
// recommended). The IV is prefixed to the string.
//
// This function is used in aeacus to encrypt reported vulnerability data to
// the remote scoring endpoint (ex. minos).
func encryptString(password, plainText string) string {
	// Create a sha256sum hash of the password provided.
	hasher := sha256.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		errorPrint(err)
	}
	key := hasher.Sum(nil)

	// Pad plainText to be a 16-byte block.
	paddingArray := make([]byte, (aes.BlockSize - len(plainText)%aes.BlockSize))
	for char := range paddingArray {
		paddingArray[char] = 0x20 // Padding with space character.
	}
	plainText = plainText + string(paddingArray)
	if len(plainText)%aes.BlockSize != 0 {
		panic("plainText is not a multiple of block size!")
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

	// Encrypt and seal plainText.
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plainText), nil)
	ciphertext = []byte(fmt.Sprintf("%s%s", nonce, ciphertext))

	return string(ciphertext)
}

// decryptString takes a password and a ciphertext and returns a decrypted
// byte sequence (as a string). The function uses typical AES-GCM.
func decryptString(password, ciphertext string) string {
	// Create a sha256sum hash of the password provided.
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(password)); err != nil {
		errorPrint(err)
	}
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
	plainText, err := aesgcm.Open(nil, iv, []byte(ciphertext), nil)
	if err != nil {
		failPrint(err.Error())
		return ""
	}

	return strings.TrimSpace(string(plainText))
}
