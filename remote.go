package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Use non-ASCII bytes as a delimiter.
var delimiter = string(byte(255)) + string(byte(222))

const (
	FAIL  = "FAIL"
	GREEN = "green"
	RED   = "red"
)

func readTeamID() {
	fileContent, err := readFile(dirPath + "TeamID.txt")
	fileContent = strings.TrimSpace(fileContent)
	if err != nil {
		teamID = ""
		if conf.Remote != "" {
			fail("TeamID.txt does not exist!")
			conn.OverallColor = RED
			conn.OverallStatus = "Your TeamID files does not exist! Failed to score image."
			conn.Status = false
		} else {
			warn("TeamID.txt does not exist! This image is local only, so we will continue.")
		}
		sendNotification("TeamID.txt does not exist!")
	} else if fileContent == "" {
		teamID = ""
		fail("TeamID.txt is empty!")
		sendNotification("TeamID.txt is empty!")
		if conf.Remote != "" {
			conn.OverallStatus = RED
			conn.OverallStatus = "Your TeamID is empty! Failed to score image."
			conn.Status = false
		}
	} else {
		teamID = fileContent
	}
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
	writeString(&update, "team", teamID)
	writeString(&update, "image", conf.Name)
	writeString(&update, "score", strconv.Itoa(image.Score))
	writeString(&update, "vulns", genVulns())
	writeString(&update, "time", strconv.Itoa(int(time.Now().Unix())))
	info("Encrypting score update...")
	if err := deobfuscateData(&conf.Password); err != nil {
		fail(err)
		return "", err
	}
	finishedUpdate := hexEncode(encryptString(conf.Password, update.String()))
	if err := obfuscateData(&conf.Password); err != nil {
		fail(err)
		return "", err
	}
	return finishedUpdate, nil
}

func genVulns() string {
	var vulnString strings.Builder

	// Vulns achieved
	vulnString.WriteString(fmt.Sprintf("%d%s", len(image.Points), delimiter))
	// Total vulns
	vulnString.WriteString(fmt.Sprintf("%d%s", image.ScoredVulns, delimiter))

	// Build vuln string
	for _, penalty := range image.Penalties {
		if err := deobfuscateData(&penalty.Message); err != nil {
			fail(err)
		}
		vulnString.WriteString(fmt.Sprintf("%s - N%.0f pts", penalty.Message, math.Abs(float64(penalty.Points))))
		if err := obfuscateData(&penalty.Message); err != nil {
			fail(err)
		}
		vulnString.WriteString(delimiter)
	}

	for _, point := range image.Points {
		if err := deobfuscateData(&point.Message); err != nil {
			fail(err)
		}
		vulnString.WriteString(fmt.Sprintf("%s - %d pts", point.Message, point.Points))
		if err := obfuscateData(&point.Message); err != nil {
			fail(err)
		}
		vulnString.WriteString(delimiter)
	}

	info("Encrypting vulnerabilities...")

	deobfuscateData(&conf.Password)
	finishedVulns := hexEncode(encryptString(conf.Password, vulnString.String()))
	obfuscateData(&conf.Password)
	return finishedVulns
}

func reportScore() error {
	update, err := genUpdate()
	if err != nil {
		fail(err.Error())
		return err
	}
	resp, err := http.PostForm(conf.Remote+"/update",
		url.Values{"update": {update}})
	if err != nil {
		fail(err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		conn.OverallColor = RED
		conn.OverallStatus = "Failed to upload score! Please ensure that your Team ID is correct."
		conn.Status = false
		fail("Failed to upload score!")
		sendNotification("Failed to upload score!")
		return errors.New("Non-200 response from remote scoring endpoint")
	}
	return nil
}

func checkServer() {
	// Internet check (requisite)
	info("Checking for internet connection...")

	// Poor example.org :(
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	_, err := client.Get("http://example.org")

	if err != nil {
		conn.NetColor = RED
		conn.NetStatus = FAIL
	} else {
		conn.NetColor = GREEN
		conn.NetStatus = "OK"
	}

	// Scoring engine check
	info("Checking for scoring engine connection...")
	resp, err := client.Get(conf.Remote + "/status/" + teamID + "/" + conf.Name)

	if err != nil {
		conn.ServerColor = RED
		conn.ServerStatus = FAIL
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fail("Error reading Status body.")
			conn.ServerColor = RED
			conn.ServerStatus = FAIL
		} else {
			if resp.StatusCode == 200 {
				conn.ServerColor = GREEN
				conn.ServerStatus = "OK"
			} else {
				conn.ServerColor = RED
				conn.ServerStatus = "ERROR"
			}
			handleStatus(string(body))
		}
	}

	// Overall
	if conn.NetStatus == FAIL && conn.ServerStatus == "OK" {
		timeStart = time.Now()
		conn.OverallColor = "goldenrod"
		conn.OverallStatus = "Server connection good but no Internet. Assuming you're on an isolated LAN."
		conn.Status = true
	} else if conn.ServerStatus == FAIL {
		timeStart = time.Now()
		conn.OverallColor = RED
		conn.OverallStatus = "Failure! Can't access remote scoring server."
		fail("Can't access remote scoring server!")
		sendNotification("Score upload failure! Unable to access remote server.")
		conn.Status = false
	} else if conn.ServerStatus == "ERROR" {
		timeWithoutID = time.Since(timeStart)
		conn.OverallColor = RED
		conn.OverallStatus = "Scoring engine rejected your TeamID!"
		fail("Remote server returned an error for its status! Your ID is probably wrong.")
		sendNotification("Status check failed, TeamID incorrect!")
		conn.Status = false
	} else if conn.ServerStatus == "DISABLED" {
		conn.OverallColor = RED
		conn.OverallStatus = "Remote scoring server is no longer accepting scores."
		fail("Remote scoring server is no longer accepting scores.")
		sendNotification("Remote scoring server is no longer accepting scores.")
		conn.Status = false
	} else {
		timeStart = time.Now()
		conn.OverallColor = GREEN
		conn.OverallStatus = "OK"
		conn.Status = true
	}
}

func handleStatus(status string) {
	var statusStruct statusRes
	if err := json.Unmarshal([]byte(status), &statusStruct); err != nil {
		fail("Failed to parse JSON response (status): " + err.Error())
	}
	switch statusStruct.Status {
	case "DISABLED":
		conn.ServerStatus = "DISABLED"
	}
}
