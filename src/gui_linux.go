package main

func launchIDPrompt() {
	teamID, err := shellCommandOutput(`
		#!/bin/bash
		teamid=$(
			zenity --entry \
			--title="TeamID" \
			--text="Enter your TeamID:"
		)
		echo $teamid
	`)
	if err == nil {
		writeFile(mc.DirPath+"TeamID.txt", teamID)
	} else {
		failPrint("Error saving TeamID!")
		sendNotification("Error saving TeamID!")
	}
}

func launchConfigGui() {
	warnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := shellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		infoPrint("Configuration successfully written to scoring.conf!")
	}
}
