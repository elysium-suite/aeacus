package main

import (
	"fmt"
    "time"
	"runtime"
    "strings"
	"strconv"
	"net/url"
	"net/http"
	"io/ioutil"
	"crypto/sha1"
	"encoding/hex"
)

func readTeamID(mc *metaConfig, id *imageData) {
    fileContent := ""
    err := error(nil)
	if runtime.GOOS == "linux" {
		fileContent, err = readFile("/opt/aeacus/misc/TeamID.txt")
	} else {
		fileContent, err = readFile("C:\\Users\\" + mc.Config.User + "\\Desktop\\TeamID.txt")
	}
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
    	mc.TeamID = fileContent
    }
}

// TODO
func getAuthToken() {
	fmt.Println("init connection")
}

func genChallenge() string {
    baseString := "71844fd169e20dc88ce6f985b42611cfb31cf196"
	genTime := time.Now()
	genTimeHash := "e31ab5a0097531f6d8499f761edf0ab3a8eb6e5f"
	hasher := sha1.New()
	hasher.Write([]byte(genTime.Format("2006/01/02 15:04")))
	genTimeHash = hex.EncodeToString(hasher.Sum(nil))
    chalString := xor(baseString, genTimeHash)
    return chalString
}

func reportScore(mc *metaConfig, id *imageData) {
	resp, err := http.PostForm(mc.Config.Remote+"/scores/css/update",
		url.Values{"team": {mc.TeamID},
			"image":     {mc.Config.Name},
			"score":     {strconv.Itoa(id.Score)},
			"challenge": {genChallenge()},
            "id": {"id"}})
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
        if strings.ToLower(mc.Config.Local) != "yes" {
            if mc.Cli.Bool("v") {

                warnPrint("Local is not set to \"yes\". Clearing scoring data.")
            }
            clearImageData(id)
        }
	}
}

func checkScoring(mc *metaConfig) bool {
	// hit endpoint with check
    // get status?
	return true
}

func checkServer(mc *metaConfig, id *imageData) {

	// Internet check (requisite)
	if mc.Cli.Bool("v") {
		infoPrint("Checking for internet connection...")
	}
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		id.ConnStatus[2] = "red"
		id.ConnStatus[3] = "FAIL"
	} else {
		id.ConnStatus[2] = "green"
		id.ConnStatus[3] = "OK"
	}

	// Scoring engine check (required)
	if mc.Cli.Bool("v") {
		infoPrint("Checking for scoring engine connection...")
	}
	_, err = http.Get(mc.Config.Remote)
	if err != nil {
		id.ConnStatus[4] = "red"
		id.ConnStatus[5] = "FAIL"
	} else {
		if checkScoring(mc) {
			id.ConnStatus[4] = "green"
			id.ConnStatus[5] = "OK"
		} else {
			id.ConnStatus[4] = "yellow"
			id.ConnStatus[5] = "ERROR"
		}
	}

	// Overall
	if id.ConnStatus[3] == "FAIL" && id.ConnStatus[5] == "OK" {
		id.ConnStatus[0] = "yellow"
		id.ConnStatus[1] = "Server connection good but no Internet. Assuming you're on an isolated LAN."
        id.Connection = true
	} else if id.ConnStatus[5] == "FAIL" {
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Failure! Can't access remote scoring server."
        failPrint("Can't access remote scoring server!")
		id.Connection = false
	} else if id.ConnStatus[4] == "ERROR" {
		id.ConnStatus[0] = "red"
		id.ConnStatus[1] = "Score upload failure! Can't send scores to remote server."
        failPrint("Remote server returned an error for its status!")
        id.Connection = false
	} else {
		id.ConnStatus[0] = "green"
		id.ConnStatus[1] = "OK"
        id.Connection = true
	}

	readTeamID(mc, id)
}
