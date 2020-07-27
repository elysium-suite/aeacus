package main

import (
	"fmt"
	"strconv"
)

func scoreImage(mc *metaConfig, id *imageData) {
	// Check connection and configuration
	if mc.Config.Remote != "" && mc.Config.Local != "yes" {
		checkServer(mc, id)
		if !id.Connection {
			genReport(mc, id)
			return
		}
	}

	scoreChecks(mc, id)
	genReport(mc, id)

	if mc.Config.Remote != "" && mc.Config.Local != "yes" {
		reportScore(mc, id)
	}

	// Check if points increased/decreased
	prevPoints, err := readFile(mc.DirPath + "/previous.txt")
	if err == nil {
		prevScore, _ := strconv.Atoi(prevPoints)
		if prevScore < id.Score {
			sendNotification(mc, "You gained points!")
			playAudio(mc.DirPath + "assets/gain.wav")
		} else if prevScore > id.Score {
			sendNotification(mc, "You lost points!")
			playAudio(mc.DirPath + "assets/alarm.wav")
		}
	}

	writeFile(mc.DirPath+"/previous.txt", strconv.Itoa(id.Score))

	if mc.Config.Remote != "" && mc.Config.Local == "yes" {
		checkServer(mc, id)
		if id.Connection {
			reportScore(mc, id)
		}
	}
}

func scoreChecks(mc *metaConfig, id *imageData) {
	clearImageData(id)
	pointlessChecks := []int{}

	for i, check := range mc.Config.Check {
		if check.Points == 0 {
			pointlessChecks = append(pointlessChecks, i)
			id.ScoredVulns++
		} else if check.Points > 0 {
			id.TotalPoints += check.Points
			id.ScoredVulns++
		}
	}

	pointsLeft := 100 - id.TotalPoints
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
		id.TotalPoints += (pointsEach * len(pointlessChecks))
		if id.TotalPoints < 100 {
			for i := 0; id.TotalPoints < 100; id.TotalPoints++ {
				mc.Config.Check[pointlessChecks[i]].Points++
				i++
				if i > len(pointlessChecks)-1 {
					i = 0
				}
			}
			id.TotalPoints += (100 - id.TotalPoints)
		}
	}

	for _, check := range mc.Config.Check {
		status := true
		for _, condition := range check.Pass {
			status = processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
			if status {
				break
			}
		}
		for _, condition := range check.Fail {
			failStatus := processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
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
				id.Points = append(id.Points, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Contribs += check.Points
			}
		} else {
			if status {
				if verboseEnabled {
					failPrint("")
					fmt.Printf("Penalty triggered: %s - %d pts\n", check.Message, check.Points)
				}
				id.Penalties = append(id.Penalties, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Detracts += check.Points
			}
		}
	}
	if verboseEnabled {
		infoPrint(fmt.Sprintf("Score: %d", id.Score))
	}
}
