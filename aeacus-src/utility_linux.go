package main

import (
	"os/exec"
)

func shellCommand(commandGiven string) {
	cmd := exec.Command("sh", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			failPrint("Command \"" + commandGiven + "\" errored out (" + err.Error() + ").")
		}
	}
}

func sendNotification(userName string, notifyText string) {
	commandText := "/sbin/runuser -l " + userName + " -c  '/usr/bin/notify-send -i /opt/aeacus/web/assets/logo.png \"Aeacus Scoring System\" \"" + notifyText + "\"'"
	shellCommand(commandText)
}

func destroyImage() {
	warnPrint("Destroying the image! (jk for now. that's dangerous)")
}
