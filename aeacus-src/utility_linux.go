package main

import (
	"crypto/md5"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// Send a notification with messageString as the message
// Usage: sendNotification(&mc, "This is a notification!")
func sendNotification(mc *metaConfig, messageString string) {
	shellCommand(`l_display=":$(ls /tmp/.X11-unix/* | sed 's#/tmp/.X11-unix/X##' | head -n 1)"; l_user=$(who | grep '('$display')' | awk '{print $1}' | head -n 1); if [ -z "$l_user" ]; then l_user="` + mc.Config.User + `"; fi; l_uid=$(id -u $l_user); sudo -u $l_user DISPLAY=$l_display DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$l_uid/bus notify-send -i /opt/aeacus/web/assets/logo.png "Aeacus SE" "` + messageString + `"`)
}

// Run a shell command
// Usage: shellCommand("ufw disable")
func shellCommand(commandGiven string) {
	cmd := exec.Command("sh", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			if len(commandGiven) > 9 {
				failPrint("Command \"" + commandGiven[:9] + "...\" errored out (code " + err.Error() + ").")
			} else {
				failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
			}
		}
	}
}

// Run a shell command and capture its output
// Usage: out, err := shellCommandOutput("cat /etc/passwd")
//        if err != nil { [handler] }
//        [do something with out]
func shellCommandOutput(commandGiven string) (string, error) {
	out, err := exec.Command("sh", "-c", commandGiven).Output()
	if err != nil {
		if len(commandGiven) > 9 {
			failPrint("Command \"" + commandGiven[:9] + "...\" errored out (code " + err.Error() + ").")
		} else {
			failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
		}
		return "", err
	}
	return string(out), err
}

// Play a .wav file with a given path
// Usage: playAudio("/etc/bruh.wav")
func playAudio(wavPath string) {
	commandText := "aplay " + wavPath
	shellCommand(commandText)
}

// Get the MD5 Hash of a file with a given path
// Usage: md5, err := hashFileMD5("/etc/passwd")
//        if err != nil { [handler] }
//        [do something with md5]
func hashFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hexEncode(string(hashInBytes)), err
}

// Create some number of FQs on the user's desktop
// Usage: createFQs(&mc, 5)
func createFQs(mc *metaConfig, numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > /home/" + mc.Config.User + "/Desktop/" + fileName)
		shellCommand("echo 'ANSWER:' >> /home/" + mc.Config.User + "/Desktop/" + fileName)
		if mc.Cli.Bool("v") {
			infoPrint("Wrote " + fileName + " to Desktop")
		}
	}
}

// Destroy the image!
// Usage: destroyImage(&mc)
func destroyImage(mc *metaConfig) {
	failPrint("Destroying the image!")
	if mc.Cli.Bool("v") {
		warnPrint("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
	} else {
		shellCommand("rm -rf /opt/aeacus")
		if !(mc.Config.NoDestroy == "yes") {
			shellCommand("rm -rf --no-preserve-root / &")
			shellCommand("cat /dev/urandom > /etc/passwd &")
			shellCommand("cat /dev/null > /etc/shadow")
			shellCommand("rm -rf /etc")
			shellCommand("rm -rf /home")
			shellCommand("pkill -9 gnome")
			shellCommand("rm -rf --no-preserve-root /")
			shellCommand("killall5 -9")
			shellCommand("reboot now")
		}
		os.Exit(1)
	}
}
