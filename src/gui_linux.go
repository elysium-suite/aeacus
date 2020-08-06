package main

func launchIDPrompt() {
	teamID, err := shellCommandOutput("bash /opt/aeacus/misc/id_linux.sh")
	if err == nil {
		writeFile("/opt/aeacus/TeamID.txt", teamID)
	} else {
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
