package main

import (
	"github.com/gen2brain/beeep"
	"os"
)

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
