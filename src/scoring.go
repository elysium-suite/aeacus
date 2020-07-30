package main

import (
	"fmt"
	"strconv"
)

func scoreImage() {
	checkConfigData()

	if mc.Config.Local {
		scoreChecks()
		if mc.Config.Remote != "" {
			checkServer()
			if mc.Connection {
				reportScore()
			}
		}
		genReport(mc.Image)

	} else {
		if mc.Config.Remote != "" {
			checkServer()
			if !mc.Connection {
				genReport(mc.Image)
				return
			}
		}
		scoreChecks()
		err := reportScore()
		if err != nil {
			return
		}
		genReport(mc.Image)
	}

	// Check if points increased/decreased
	prevPoints, err := readFile(mc.DirPath + "/previous.txt")
	if err == nil {
		prevScore, _ := strconv.Atoi(prevPoints)
		if prevScore < mc.Image.Score {
			sendNotification("You gained points!")
			playAudio(mc.DirPath + "assets/gain.wav")
		} else if prevScore > mc.Image.Score {
			sendNotification("You lost points!")
			playAudio(mc.DirPath + "assets/alarm.wav")
		}
	} else {
		warnPrint("Reading from previous.txt failed.")
	}

	writeFile(mc.DirPath+"/previous.txt", strconv.Itoa(mc.Image.Score))
}

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

func scoreChecks() {
	mc.Image = imageData{}
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
	if pointsLeft < 0 && len(pointlessChecks) > 0 {
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

	for _, check := range mc.Config.Check {
		status := true
		for _, condition := range check.Pass {
			status = processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
			if reverseEnabled {
				status = !status
			}
			if status {
				break
			}
		}
		for _, condition := range check.Fail {
			failStatus := processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
			if reverseEnabled {
				failStatus = !failStatus
			}
			if failStatus {
				status = false
				break
			}
		}
		if check.Points >= 0 {
			if status {
				if verboseEnabled {
					passPrint("")
					fmt.Printf("Check passed: %s - %d pts\n", check.Message, check.Points)
				}
				mc.Image.Points = append(mc.Image.Points, scoreItem{check.Message, check.Points})
				mc.Image.Score += check.Points
				mc.Image.Contribs += check.Points
			}
		} else {
			if status {
				if verboseEnabled {
					failPrint("")
					fmt.Printf("Penalty triggered: %s - %d pts\n", check.Message, check.Points)
				}
				mc.Image.Penalties = append(mc.Image.Penalties, scoreItem{check.Message, check.Points})
				mc.Image.Score += check.Points
				mc.Image.Detracts += check.Points
			}
		}
	}
	if verboseEnabled {
		infoPrint(fmt.Sprintf("Score: %d", mc.Image.Score))
	}
}
