package cmd

func LaunchIDPrompt() {
	err := shellCommand(`zenity --entry --title="TeamID" --text "Enter your TeamID:" --width=200 >/opt/aeacus/TeamID.txt`)
	if err != nil {
		sendNotification("Error saving TeamID!")
	}
}

func LaunchConfigGui() {
	warnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := shellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		infoPrint("Configuration successfully written to scoring.conf!")
	}
}
