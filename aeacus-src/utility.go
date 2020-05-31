package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"os"
)

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

func runningPermsCheck() {
	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}
