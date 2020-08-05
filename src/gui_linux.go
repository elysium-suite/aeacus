package main

func launchIDPrompt() {
	warnPrint("Custom ID prompt not supported in Linux yet. Gotta use gedit.")
}

func launchConfigGui() {
	warnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := shellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		infoPrint("Configuration successfully written to scoring.conf!")
	}
}
