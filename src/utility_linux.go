package main

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

// readFile (Linux) wraps ioutil's ReadFile function.
func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	return string(fileContent), err
}

// decodeString (linux) strictly does nothing, however it's here
// for compatibility with Windows ANSI/UNICODE/etc.
func decodeString(fileContent string) (string, error) {
	return fileContent, nil
}

// sendNotification sends a notification to the end user.
func sendNotification(mc *metaConfig, messageString string) {
	shellCommand(`l_display=":$(ls /tmp/.X11-unix/* | sed 's#/tmp/.X11-unix/X##' | head -n 1)"; l_user=$(who | grep '('$display')' | awk '{print $1}' | head -n 1); if [ -z "$l_user" ]; then l_user="` + mc.Config.User + `"; fi; l_uid=$(id -u $l_user); sudo -u $l_user DISPLAY=$l_display DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$l_uid/bus notify-send -i /opt/aeacus/assets/logo.png "Aeacus SE" "` + messageString + `"`)
}

// createFQs is a quality of life function that creates Forensic Question files
// on the Desktop, pre-populated with a template.
func createFQs(mc *metaConfig, numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > /home/" + mc.Config.User + "/Desktop/" + fileName)
		shellCommand("echo 'ANSWER:' >> /home/" + mc.Config.User + "/Desktop/" + fileName)
		if verboseEnabled {
			infoPrint("Wrote " + fileName + " to Desktop")
		}
	}
}

// shellCommand executes a given command in a sh environment, and prints an
// error if one occurred.
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

// shellCommandOutput executes a given command in a sh environment, and
// returns its output and error (if one occurred).
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

// playAudio plays a .wav file with the given path.
func playAudio(wavPath string) {
	commandText := "aplay " + wavPath
	shellCommand(commandText)
}

// hashFileMD5 generates the MD5 Hash of a file with the given path.
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

func adminCheck() bool {
	currentUser, err := user.Current()
	uid, _ := strconv.Atoi(currentUser.Uid)
	if err != nil {
		failPrint("Error for checking if running as root: " + err.Error())
		return false
	} else if uid != 0 {
		return false
	}
	return true
}

// destroyImage removes the aeacus directory (to stop scoring) and optionally
// can destroy the entire machine.
func destroyImage(mc *metaConfig) {
	failPrint("Destroying the image!")
	if verboseEnabled {
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
