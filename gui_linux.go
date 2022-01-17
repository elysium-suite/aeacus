package main

func launchIDPrompt() {
	err := shellCommand(`zenity --title "Team ID Prompt" --text "Enter your Team ID below!" --entry > ` + dirPath + "TeamID.txt")
	if err != nil {
		fail("Error running ID prompt command: " + err.Error())
	}
}
