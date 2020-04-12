package main

import (
	"fmt"
	"strconv"
)

func scoreImage(mc *metaConfig, id *imageData) {

	// Check connection and configuration
	if mc.Config.Remote != "" {
		checkServer(mc, id)
		if !id.Connection {
			genReport(mc, id)
			return
		}
	}

	scoreChecks(mc, id)
	if mc.Config.Remote != "" {
		reportScore(mc, id)
	}
	genReport(mc, id)

	// Check if points increased/decreased
	prevPoints, err := readFile(mc.DirPath + "web/assets/previous.txt")
	if err == nil {
		prevScore, _ := strconv.Atoi(prevPoints)
		if prevScore < id.Score {
			sendNotification(mc.Config.User, "You gained points!")
			playAudio(mc.DirPath + "web/assets/gain.wav")
		} else if prevScore > id.Score {
			sendNotification(mc.Config.User, "You lost points!")
			playAudio(mc.DirPath + "web/assets/alarm.wav")
		}
	} else {
		fmt.Println(err)
	}

	writeFile(mc.DirPath+"web/assets/previous.txt", strconv.Itoa(id.Score))
}

func scoreChecks(mc *metaConfig, id *imageData) {

	clearImageData(id)
	pointlessChecks := []int{}

	for i, check := range mc.Config.Check {
		if check.Points == 0 {
			pointlessChecks = append(pointlessChecks, i)
			id.ScoredVulns += 1
		} else if check.Points > 0 {
			id.TotalPoints += check.Points
			id.ScoredVulns += 1
		}
	}

	pointsLeft := 100 - id.TotalPoints
	if pointsLeft > 0 && len(pointlessChecks) > 0 {
		pointsEach := pointsLeft / len(pointlessChecks)
		for _, check := range pointlessChecks {
			mc.Config.Check[check].Points = pointsEach
		}
		id.TotalPoints += (pointsEach * len(pointlessChecks))
		if id.TotalPoints != 100 {
			mc.Config.Check[pointlessChecks[0]].Points += (100 - id.TotalPoints)
			id.TotalPoints += (100 - id.TotalPoints)
		}
	}

	for _, check := range mc.Config.Check {
		status := false
		failStatus := false
		for _, condition := range check.Pass {
			status = processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
			if status {
				break
			}
		}
		for _, condition := range check.Fail {
			failStatus = processCheckWrapper(&check, condition.Type, condition.Arg1, condition.Arg2, condition.Arg3)
			if failStatus {
				status = false
				break
			}
		}
		if check.Points >= 0 {
			if status {
				if mc.Cli.Bool("v") {
					passPrint("")
					fmt.Printf("Check passed: %s - %d pts\n", check.Message, check.Points)
				}
				id.Points = append(id.Points, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Contribs += check.Points
			}
		} else {
			if status {
				if mc.Cli.Bool("v") {
					failPrint("")
					fmt.Printf("Penalty triggered: %s - %d pts\n", check.Message, check.Points)
				}
				id.Penalties = append(id.Penalties, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Detracts += check.Points
			}
		}
	}
	if mc.Cli.Bool("v") {
		infoPrint("")
		fmt.Printf("Score: %d\n", id.Score)
	}
}

func clearImageData(id *imageData) {
	id.Score = 0
	id.ScoredVulns = 0
	id.TotalPoints = 0
	id.Contribs = 0
	id.Detracts = 0
	id.Points = []scoreItem{}
	id.Penalties = []scoreItem{}
}
