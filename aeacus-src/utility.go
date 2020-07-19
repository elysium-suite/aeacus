package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"os"
	"time"
)

var aeacusVersion = "1.1.1"

func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func sendNotification(mc *metaConfig, messageString string) {
	// does NOT work for linux --> run as root, doesnt send notify
	// to all users on the system
	err := beeep.Notify("Aeacus SE", messageString, mc.DirPath+"web/assets/logo.png")
	if err != nil {
		failPrint("Notification error: " + err.Error())
	}
}

func timeCheck(mc *metaConfig) {
	if mc.Config.EndDate != "" {
		endDate, err := time.Parse("2006/01/02 15:04:05 MST", mc.Config.EndDate)
		if err != nil {
			failPrint("Your EndDate value in the configuration is invalid.")
		} else {
			if time.Now().After(endDate) {
				destroyImage(mc)
			}
		}
	}
}

func runningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}
