package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/safinsingh/aeaconf2"
)

var (
	teamID string
	conf   = &config{}
	image  = &imageData{}
	conn   = &connData{}

	// checkCount keeps track of the current check being scored, and is used
	// for identifying which check caused a given error or warning.
	checkCount int
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
	Hints       []hintItem
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
	Index   int
	Message string
	Points  int
}

// hintItem is the scoring report representation of a hint, which can contain
// multiple messages.
type hintItem struct {
	Index    int
	Messages []string
	Points   int
}

// config is a representation of the TOML configuration typically
// specific in scoring.conf.
type config struct {
	DisableRemoteEncryption bool
	Local                   bool
	Shell                   bool
	EndDate                 string
	Name                    string
	OS                      string
	Password                string
	Remote                  string
	Title                   string
	User                    string
	Version                 string
	MaxPoints               int
	Checks                  []*aeaconf2.Check
}

// statusRes is to parse a JSON response from the remote server.
type statusRes struct {
	Status string `json:"status"`
}

// readScoringData is a convenience function around readData and decodeString,
// which parses the encrypted scoring configuration file.
func readScoringData() error {
	info("Decrypting data from " + dirPath + scoringData + "...")

	// Read in the encrypted configuration file
	dataFile, err := readFile(dirPath + scoringData)
	if err != nil {
		return err
	} else if dataFile == "" {
		return errors.New("Scoring data is empty!")
	}

	decryptedData, err := decryptConfig(dataFile)
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
	if timeCheck() {
		log.Fatal("Image is running outside of the specified end date.")
	}
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
			warn("Reporting image score failed, and local is disabled. Score data removed.")
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
	for i, check := range conf.Check {
		// checkCount is the same as in printConfig, 1-based count
		checkCount = i + 1
		scoreCheck(check)
	}
	checkCount = 0
	info(fmt.Sprintf("Score: %d", image.Score))
}

// scoreCheck will go through each condition inside a check, and determine
// whether or not the check passes.
func scoreCheck(check check) {
	status := false
	failed := false

	// Create hint var in case any checks have hints
	hint := hintItem{
		Index:  checkCount,
		Points: check.Points,
	}

	// If a Fail condition passes, the check fails, no other checks required.
	if len(check.Fail) > 0 {
		failed = checkOr(check.Fail, &hint)
	}

	// If a PassOverride succeeds, that overrides the Pass checks
	if !failed && len(check.PassOverride) > 0 {
		status = checkOr(check.PassOverride, &hint)
	}

	// Finally, we check the normal Pass checks
	if !failed && !status && len(check.Pass) > 0 {
		status = checkAnd(check.Pass, &hint)
	}

	if status {
		if check.Points >= 0 {
			if verboseEnabled {
				deobfuscateData(&check.Message)
				pass(fmt.Sprintf("Check passed: %s - %d pts", check.Message, check.Points))
				obfuscateData(&check.Message)
			}
			image.Points = append(image.Points, scoreItem{checkCount, check.Message, check.Points})
			image.Contribs += check.Points
		} else {
			if verboseEnabled {
				deobfuscateData(&check.Message)
				fail(fmt.Sprintf("Penalty triggered: %s - %d pts", check.Message, check.Points))
				obfuscateData(&check.Message)
			}
			image.Penalties = append(image.Penalties, scoreItem{checkCount, check.Message, check.Points})
			image.Detracts += check.Points
		}
		image.Score += check.Points
	} else {
		// If there is a check-wide hint, add to start of hint messages.
		if check.Hint != "" {
			hints := []string{check.Hint}
			hints = append(hints, hint.Messages...)
			hint.Messages = hints
		}

		// If the check failed, and there are hints, see if we should display them.
		// All hints triggered (based on which conditions ran) are displayed in sequential order.
		if len(hint.Messages) > 0 {
			image.Hints = append(image.Hints, hint)
		}
	}

	// If check is not a penalty, add to total
	if check.Points >= 0 {
		image.ScoredVulns++
		image.TotalPoints += check.Points
	}
}

// checkOr runs a set of conditions and returns true if any of them pass.
// It is a logical "OR".
func checkOr(conds []cond, hint *hintItem) bool {
	for _, cond := range conds {
		if runCheck(cond) {
			return true
		}
		if cond.Hint != "" {
			hint.Messages = append(hint.Messages, cond.Hint)
		}
	}
	return false
}

// checkAnd runs a set of conditions and returns true iff ALL of them pass.
// It is a logical "AND".
func checkAnd(conds []cond, hint *hintItem) bool {
	for _, cond := range conds {
		if !runCheck(cond) {
			if cond.Hint != "" {
				hint.Messages = append(hint.Messages, cond.Hint)
			}
			return false
		}
	}
	return true
}
