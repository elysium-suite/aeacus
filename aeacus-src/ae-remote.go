package main

import (
	"fmt"
	"io/ioutil"
    "time"
	"net/http"
	"net/url"
	"encoding/hex"
	"os"
	"crypto/sha1"
	"runtime"
	"strconv"
)

func readTeamID(mc *metaConfig) {
	if runtime.GOOS == "linux" {
		fileContent, err := readFile("/home/" + mc.Config.User + "/Desktop/TeamID.txt")
		if err != nil {
			failPrint("TeamID.txt does not exist!")
			os.Exit(1)
		}
		if fileContent == "" {
			failPrint("TeamID.txt is empty!")
			os.Exit(1)
		}
		mc.TeamID = fileContent
	} else {
		fileContent, err := readFile("C:\\Users\\" + mc.Config.User + "\\Desktop\\TeamID.txt")
		if err != nil {
			failPrint("TeamID.txt does not exist! (And also your OS is lame).")
			os.Exit(1)
		}
		if fileContent == "" {
			failPrint("TeamID.txt is empty! (And also Windows is lame).")
			os.Exit(1)
		}
		mc.TeamID = fileContent
	}
	// teamid validity checks here
	// todo... what does that even look like?
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
	resp, err := http.PostForm("http://"+mc.Config.Remote+"/scores/css/update",
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
	if string(body) == "FAIL" {
		failPrint("Failed to upload score! Is your TeamID wrong?")
        os.Exit(1)
	}
}

func checkScoring(mc *metaConfig) bool {
	// hit endpoint with check
	return true
}

func checkServer(mc *metaConfig) ([]string, bool) {

	connStatus := make([]string, 6)

	// Internet check (requisite)
	if mc.Cli.Bool("v") {
		infoPrint("Checking for internet connection...")
	}
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		connStatus[2] = "red"
		connStatus[3] = "FAIL"
	} else {
		connStatus[2] = "green"
		connStatus[3] = "OK"
	}

	// Scoring engine check (required)
	if mc.Cli.Bool("v") {
		infoPrint("Checking for scoring engine connection...")
	}
	_, err = http.Get("http://" + mc.Config.Remote)
	if err != nil {
		connStatus[4] = "red"
		connStatus[5] = "FAIL"
	} else {
		if checkScoring(mc) {
			connStatus[4] = "green"
			connStatus[5] = "OK"
		} else {
			connStatus[4] = "yellow"
			connStatus[5] = "ERROR"
		}
	}

	// Overall
	if connStatus[3] == "FAIL" && connStatus[5] == "OK" {
		connStatus[0] = "yellow"
		connStatus[1] = "Server connection good but no Internet. Assuming you're on an isolated LAN."
		return connStatus, true
	} else if connStatus[5] == "FAIL" {
		connStatus[0] = "red"
		connStatus[1] = "Failure! Can't access remote scoring server."
		return connStatus, false
	} else if connStatus[4] == "ERROR" {
		connStatus[0] = "red"
		connStatus[1] = "Score upload failure! Can't send scores to remote server."
		return connStatus, false
	} else {
		connStatus[0] = "green"
		connStatus[1] = "OK"
		return connStatus, true
	}

}
