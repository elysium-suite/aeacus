package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func readTeamID(mc *metaConfig, id *imageData) {
	fileContent := ""
	err := error(nil)
	fileContent, err = readFile(mc.DirPath + "misc/TeamID.txt")
	if err != nil {
		failPrint("TeamID.txt does not exist!")
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Your TeamID files does not exist! Failed to upload scores."
		id.Connection = false
	} else if fileContent == "" {
		failPrint("TeamID.txt is empty!")
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Your TeamID is empty! Failed to upload scores."
		id.Connection = false
	} else {
		// teamid validity checks here
		// todo... what does that even look like?
		// should there be a standard format? :thinking;
		mc.TeamID = fileContent
	}
}

// genChallenge generates a crypto challenge for the CSS endpoint
func genChallenge(mc *metaConfig) string {
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
		if mc.Cli.Bool("v") {
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
		failPrint("error occured :()")
		fmt.Println(err)
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
			if mc.Cli.Bool("v") {
				warnPrint("Local is not set to \"yes\". Clearing scoring data.")
			}
			clearImageData(id)
		}
	}
}

func checkServer(mc *metaConfig, id *imageData) {

	// Internet check (requisite)
	if mc.Cli.Bool("v") {
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
	if mc.Cli.Bool("v") {
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
