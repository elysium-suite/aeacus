package main

import (
	"fmt"
	"os/exec"
)

func shellCommand(commandGiven string) {
	cmd := exec.Command("powershell.exe", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
		}
	}
}

func sendNotification(userName string, notifyText string) {
	fmt.Println("not supported yet oopsies")
	fmt.Printf("tried to send notification as user %s with text %s", userName, notifyText)
}

func playAudio(wavPath string) {
	commandText := "(New-Object Media.SoundPlayer '" + wavPath + "').PlaySync();"
	shellCommand(commandText)
}


func destroyImage() {
	fmt.Println("cant do that yet. not supported on windows. enjoy ur undestryoed imaeg")
}
