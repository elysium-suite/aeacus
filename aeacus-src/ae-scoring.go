package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func scoreImage(mc *metaConfig, id *imageData) {
	// Check connection and configuration
	connStatus := []string{"green", "OK", "green", "OK", "green", "OK"}
	if mc.Config.Remote != "" {
		readTeamID(mc)
		connStatus, connection := checkServer(mc)
		if !connection {
			failPrint("Can't access remote scoring server!")
			genReport(mc, id, connStatus)
			os.Exit(1)
		}
	}

	// Score image
	if runtime.GOOS == "linux" {
		scoreL(mc, id)
	} else {
		scoreW(mc, id)
	}

	// Check if points increased/decreased
	prevPoints, err := readFile(mc.DirPath + "web/assets/previous.txt")
	if err == nil {
		prevScore, _ := strconv.Atoi(prevPoints)
		if prevScore < id.Score {
			sendNotification(mc.Config.User, "You gained points!")
			// TODO play gain noise
		} else if prevScore > id.Score {
			sendNotification(mc.Config.User, "You lost points!")
			// TODO play loss noise
		}
	} else {
		fmt.Println(err)
	}
	if mc.Config.Remote != "" {
		reportScore(mc, id)
	}
	writeFile(mc.DirPath+"web/assets/previous.txt", strconv.Itoa(id.Score))
	genReport(mc, id, connStatus)
}

func destroyImage() {
	// destroy the image if outside time range
	if runtime.GOOS == "linux" {
		fmt.Println("oops dfestroying thel inux system")
	} else {
		fmt.Println("cant do that yet. not supported on windows. enjoy ur undestryoed imaeg")
	}
}
