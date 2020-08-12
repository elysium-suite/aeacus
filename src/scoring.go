package main

import (
	"fmt"
	"strconv"
	"sync"
)

var points = make(map[int]scoreItem)

func scoreImage() {
	// Ensure checks aren't blank, and grab TeamID.
	checkConfigData()

	// If local is enabled, we want to:
	//    1. Score checks
	//    2. Check if server is up (if remote)
	//    3. If connection, report score
	//    4. Generate report
	if mc.Config.Local {
		scoreChecks()
		if mc.Config.Remote != "" {
			checkServer()
			if mc.Connection {
				reportScore()
			}
		}
		genReport(mc.Image)

		// If local is disabled, we want to:
		//    1. Check if server is up
		//    2. If no connection, generate report with err text
		//    3. If connection, score checks
		//    4. Report the score
		//    5. If reporting failed, show error, wipe scoring data
		//    6. Generate report
	} else {
		checkServer()
		if !mc.Connection {
			if verboseEnabled {
				warnPrint("Connection failed-- generating blank report.")
			}
			genReport(mc.Image)
			return
		}
		scoreChecks()
		err := reportScore()
		if err != nil {
			mc.Image = imageData{}
			if verboseEnabled {
				warnPrint("Local is disabled, scoring data removed.")
			}
		}
		genReport(mc.Image)
	}

	// Check if points increased/decreased
	prevPoints, err := readFile(mc.DirPath + "previous.txt")

	// Write previous.txt before playing sound, in case execution is
	// interrupted while playing it.
	writeFile(mc.DirPath+"previous.txt", strconv.Itoa(mc.Image.Score))

	if err == nil {
		prevScore, err := strconv.Atoi(prevPoints)
		if err != nil {
			failPrint("Don't mess with previous.txt!")
		} else {
			if prevScore < mc.Image.Score {
				sendNotification("You gained points!")
				playAudio(mc.DirPath + "assets/gain.wav")
			} else if prevScore > mc.Image.Score {
				sendNotification("You lost points!")
				playAudio(mc.DirPath + "assets/alarm.wav")
			}
		}
	} else {
		warnPrint("Reading from previous.txt failed. This is probably fine.")
	}

}

// checkConfigData performs preliminary checks on the configuration data,
// and reads the TeamID file.
func checkConfigData() {
	if len(mc.Config.Check) == 0 {
		mc.Conn.OverallColor = "red"
		mc.Conn.OverallStatus = "There were no checks found in the configuration."
	} else {
		// For none-remote local connections
		mc.Conn.OverallColor = "green"
		mc.Conn.OverallStatus = "OK"
	}
	readTeamID()
}

// scoreChecks runs through every check configured and runs them concurrently.
func scoreChecks() {
	mc.Image = imageData{}
	assignPoints()

	points = make(map[int]scoreItem)

	if mc.DirPath == linuxDir {
		var wg sync.WaitGroup
		var m sync.Mutex

		for index, check := range mc.Config.Check {
			wg.Add(1)
			go scoreCheck(index, check, &wg, &m)
		}

		wg.Wait()
	} else {
		for index, check := range mc.Config.Check {
			scoreCheckBlocking(index, check)
		}
	}

	// Order checks
	for index, check := range mc.Config.Check {
		if check.Points >= 0 {
			if point, ok := points[index]; ok {
				mc.Image.Points = append(mc.Image.Points, point)
			}
		} else {
			if point, ok := points[index]; ok {
				mc.Image.Penalties = append(mc.Image.Penalties, point)
			}
		}
	}

	if verboseEnabled {
		infoPrint("Finished running all checks.")
	}

	if verboseEnabled {
		infoPrint(fmt.Sprintf("Score: %d", mc.Image.Score))
	}
}

// scoreCheck will go through each condition inside a check, and determine
// whether or not the check passes. It does this concurrently.
func scoreCheck(index int, check check, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	status := true

	// If a fail condition passes, the check fails, no other checks required.
	if len(check.Fail) > 0 {
		status = checkFails(&check)
		if !status {
			return
		}
	}

	// If a PassOverride succeeds, that overrides the Pass checks
	passOverrideStatus := false
	if len(check.PassOverride) > 0 {
		passOverrideStatus = checkPassOverrides(&check)
		status = passOverrideStatus
	}

	if !passOverrideStatus && len(check.Pass) > 0 {
		status = checkPass(&check)
	}

	if status {
		if check.Points >= 0 {
			if verboseEnabled {
				passPrint(fmt.Sprintf("Check passed: %s - %d pts", check.Message, check.Points))
			}
			m.Lock()
			points[index] = scoreItem{check.Message, check.Points}
			mc.Image.Contribs += check.Points
			m.Unlock()
		} else {
			if verboseEnabled {
				failPrint(fmt.Sprintf("Penalty triggered: %s - %d pts", check.Message, check.Points))
			}
			m.Lock()
			points[index] = scoreItem{check.Message, check.Points}
			mc.Image.Detracts += check.Points
			m.Unlock()
		}
		mc.Image.Score += check.Points
	}
}

// scoreCheckBlocking will run checks non-concurrently.
func scoreCheckBlocking(index int, check check) {
	status := true

	// If a fail condition passes, the check fails, no other checks required.
	if len(check.Fail) > 0 {
		status = checkFails(&check)
		if !status {
			return
		}
	}

	// If a PassOverride succeeds, that overrides the Pass checks
	passOverrideStatus := false
	if len(check.PassOverride) > 0 {
		passOverrideStatus = checkPassOverrides(&check)
		status = passOverrideStatus
	}

	if !passOverrideStatus && len(check.Pass) > 0 {
		status = checkPass(&check)
	}

	if status {
		if check.Points >= 0 {
			if verboseEnabled {
				passPrint(fmt.Sprintf("Check passed: %s - %d pts", check.Message, check.Points))
			}
			points[index] = scoreItem{check.Message, check.Points}
			mc.Image.Contribs += check.Points
		} else {
			if verboseEnabled {
				failPrint(fmt.Sprintf("Penalty triggered: %s - %d pts", check.Message, check.Points))
			}
			points[index] = scoreItem{check.Message, check.Points}
			mc.Image.Detracts += check.Points
		}
		mc.Image.Score += check.Points
	}
}

func checkFails(check *check) bool {
	for _, condition := range check.Fail {
		failStatus := processCheckWrapper(check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
		if debugEnabled {
			infoPrint(fmt.Sprint("Result of fail check was ", failStatus))
		}
		if failStatus {
			return true
		}
	}
	return false
}

func checkPassOverrides(check *check) bool {
	for _, condition := range check.PassOverride {
		status := processCheckWrapper(check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
		if debugEnabled {
			infoPrint(fmt.Sprint("Result of pass override was ", status))
		}
		if status {
			return true
		}
	}
	return false
}

func checkPass(check *check) bool {
	status := true
	passStatus := []bool{}
	for i, condition := range check.Pass {
		passItemStatus := processCheckWrapper(check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
		passStatus = append(passStatus, passItemStatus)
		if debugEnabled {
			infoPrint(fmt.Sprint("Result of component pass check was ", passStatus[i]))
		}
	}

	// For multiple pass conditions, will only be true if ALL of them are
	for _, result := range passStatus {
		status = status && result
		if !status {
			break
		}
	}
	if debugEnabled {
		infoPrint(fmt.Sprint("Result of all pass check was ", status))
	}
	return status
}

// assignPoints is used to automatically assign points to checks that don't
// have a hardcoded points value.
func assignPoints() {
	pointlessChecks := []int{}

	for i, check := range mc.Config.Check {
		if check.Points == 0 {
			pointlessChecks = append(pointlessChecks, i)
			mc.Image.ScoredVulns++
		} else if check.Points > 0 {
			mc.Image.TotalPoints += check.Points
			mc.Image.ScoredVulns++
		}
	}

	pointsLeft := 100 - mc.Image.TotalPoints
	if pointsLeft <= 0 && len(pointlessChecks) > 0 || len(pointlessChecks) > 100 {
		// If the specified points already value over 100, yet there are
		// checks without points assigned, we assign the default point value
		// of 3 (arbitrarily chosen).
		for _, check := range pointlessChecks {
			mc.Config.Check[check].Points = 3
		}
	} else if pointsLeft > 0 && len(pointlessChecks) > 0 {
		pointsEach := pointsLeft / len(pointlessChecks)
		for _, check := range pointlessChecks {
			mc.Config.Check[check].Points = pointsEach
		}
		mc.Image.TotalPoints += (pointsEach * len(pointlessChecks))
		if mc.Image.TotalPoints < 100 {
			for i := 0; mc.Image.TotalPoints < 100; mc.Image.TotalPoints++ {
				mc.Config.Check[pointlessChecks[i]].Points++
				i++
				if i > len(pointlessChecks)-1 {
					i = 0
				}
			}
			mc.Image.TotalPoints += (100 - mc.Image.TotalPoints)
		}
	}
}
