package main

import (
	"github.com/gen2brain/beeep"
)

func sendNotification(mc *metaConfig, messageString string) {
	// does NOT work for linux --> run as root, doesnt send notify
	// to all users on the system
	err := beeep.Notify("Aeacus SE", messageString, mc.DirPath+"web/assets/logo.png")
	if err != nil {
		failPrint("Notification error: " + err.Error())
	}
}
