package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	teamID string
	conf   = &config{}
	image  = &imageData{}
	conn   = &connData{}
)

// imageData is the current scoring data for the image. It is able to be
// wiped, removed, etc, on each run without affecting anything else.
type imageData struct {
	Contribs    int
	Detracts    int
	Score       int
	ScoredVulns int
	TotalPoints int
	Penalties   []scoreItem
	Points      []scoreItem
}

// connData represents the current connectivity state of the image to the
// internet and the scoring server.
type connData struct {
	Status        bool
	OverallColor  string
	OverallStatus string
	NetColor      string
	NetStatus     string
	ServerColor   string
	ServerStatus  string
}

// scoreItem is the scoring report representation of a check, containing only
// the message and points associated with it.
type scoreItem struct {
	Message string
	Points  int
}

// config is a representation of the TOML configuration typically
// specific in scoring.conf.
type config struct {
	Local    bool
	Destroy  bool
	EndDate  string
	Name     string
	OS       string
	Password string
	Remote   string
	Title    string
	User     string
	Version  string
	Check    []check
}

// statusRes is to parse a JSON response from the remote server.
type statusRes struct {
	Status string `json:"status"`
}

// ReadScoringData is a convenience function around readData and decodeString,
// which parses the encrypted scoring configuration file.
func readScoringData() error {
	info("Decrypting data from " + dirPath + scoringData + "...")
	decryptedData, err := readData()
	if err != nil {
		fail("Error reading in scoring data: " + err.Error())
		return err
	} else if decryptedData == "" {
		fail("Scoring data is empty! Is the file corrupted?")
		return errors.New("Scoring data is empty!")
	} else {
		info("Data decryption successful!")
	}
	parseConfig(decryptedData)
	return nil
}

// ScoreImage is the main function for scoring the image.
func scoreImage() {
	checkTrace()
	timeCheck()
	info("Scoring image...")

	// Ensure checks aren't blank, and grab TeamID.
	checkConfigData()

	// If local is enabled, we want to:
	//    1. Score checks
	//    2. Check if server is up (if remote)
	//    3. If connection, report score
	//    4. Generate report
	if conf.Local {
		scoreChecks()
		if conf.Remote != "" {
			checkServer()
			if conn.Status {
				err := reportScore()
				if err != nil {
					fail(err)
				}
			}
		}
		genReport(image)
	} else {
		// If local is disabled, we want to:
		//    1. Check if server is up
		//    2. If no connection, generate report with err text
		//    3. If connection, score checks
		//    4. Report the score
		//    5. If reporting failed, show error, wipe scoring data
		//    6. Generate report
		checkServer()
		if !conn.Status {
			warn("Connection failed-- generating blank report.")
			genReport(image)
			return
		}
		scoreChecks()
		err := reportScore()
		if err != nil {
			image = &imageData{}
			warn("Local is disabled, scoring data removed.")
		}
		genReport(image)
	}

	// Check if points increased/decreased.
	prevPoints, err := readFile(dirPath + "assets/previous.txt")
	if err == nil {
		prevScore, err := strconv.Atoi(prevPoints)
		if err != nil {
			fail("Don't mess with previous.txt! It only helps us know when to play sound and send notifications.")
		} else {
			if prevScore < image.Score {
				sendNotification("You gained points!")
				playAudio(dirPath + "assets/wav/gain.wav")
			} else if prevScore > image.Score {
				sendNotification("You lost points!")
				playAudio(dirPath + "assets/wav/alarm.wav")
			}
		}
	} else if os.IsExist(err) {
		fail("Reading from previous.txt failed!")
	}

	// Write previous.txt from current round.
	writeFile(dirPath+"assets/previous.txt", strconv.Itoa(image.Score))

	verboseEnabled = true
	info("image is", image)
	info("conf is", conf)

	// Remove imageData for next scoring round.
	image = &imageData{}
}

// checkConfigData performs preliminary checks on the configuration data, reads
// in the TeamID, and autogenerates missing values.
func checkConfigData() {
	if len(conf.Check) == 0 {
		conn.OverallColor = "red"
		conn.OverallStatus = "There were no checks found in the configuration."
	} else {
		// For none-remote local connections
		conn.OverallColor = "green"
		conn.OverallStatus = "OK"
		conn.Status = true
	}

	readTeamID()
}

// scoreChecks runs through every check configured.
func scoreChecks() {
	for _, check := range conf.Check {
		scoreCheck(check)
	}
	info(fmt.Sprintf("Score: %d", image.Score))
}

// scoreCheck will go through each condition inside a check, and determine
// whether or not the check passes.
func scoreCheck(check check) {
	status := true

	// If a fail condition passes, the check fails, no other checks required.
	if len(check.Fail) > 0 {
		status = checkFails(&check)
		if status {
			return
		}
	}

	// If a PassOverride succeeds, that overrides the Pass checks
	passOverrideStatus := false
	if len(check.PassOverride) > 0 {
		passOverrideStatus = checkPassOverrides(&check)
		status = passOverrideStatus
	}

	// Finally, we check the normal pass checks.
	if !passOverrideStatus && len(check.Pass) > 0 {
		status = checkPass(&check)
	}

	if status {
		if check.Points >= 0 {
			if verboseEnabled {
				deobfuscateData(&check.Message)
				pass(fmt.Sprintf("Check passed: %s - %d pts", check.Message, check.Points))
				obfuscateData(&check.Message)
			}
			image.Points = append(image.Points, scoreItem{check.Message, check.Points})
			image.Contribs += check.Points
		} else {
			if verboseEnabled {
				deobfuscateData(&check.Message)
				fail(fmt.Sprintf("Penalty triggered: %s - %d pts", check.Message, check.Points))
				obfuscateData(&check.Message)
			}
			image.Penalties = append(image.Penalties, scoreItem{check.Message, check.Points})
			image.Detracts += check.Points
		}
		image.Score += check.Points
	}

	image.ScoredVulns++
	image.TotalPoints += check.Points
}

func checkFails(check *check) bool {
	for _, cond := range check.Fail {
		failStatus := runCheck(cond)
		if failStatus {
			return true
		}
	}
	return false
}

func checkPassOverrides(check *check) bool {
	for _, cond := range check.PassOverride {
		status := runCheck(cond)
		if status {
			return true
		}
	}
	return false
}

func checkPass(check *check) bool {
	status := true
	passStatus := []bool{}
	for _, cond := range check.Pass {
		passItemStatus := runCheck(cond)
		passStatus = append(passStatus, passItemStatus)
	}

	// For multiple pass conditions, will only be true if ALL of them are
	for _, result := range passStatus {
		status = status && result
		if !status {
			break
		}
	}
	return status
}
