package cmd

func launchIDPrompt() {
	// zenity prompt ree?
}

func launchConfigGui() {
	warnPrint("The script doesn't currently have the ability to add multiple check or fail conditions-- you must still do these manually.")
	_, err := shellCommandOutput("bash ./misc/gui_linux.sh")
	if err == nil {
		infoPrint("Configuration successfully written to scoring.conf!")
	}
}
