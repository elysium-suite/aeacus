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
func sendNotification(messageString string) {
	if conf.User == "" {
		fail("User not specified in configuration, can't send notification.")
	} else {
		err := shellCommand(`
			user="` + conf.User + `"
			uid="$(id -u $user)" # Ubuntu >= 18
			if [ -e /run/user/$uid/bus ]; then
			    display="unix:path=/run/user/$uid/bus"
			else # Ubuntu <= 16
			    display="unix:abstract=$(cat /run/user/$uid/dbus-session | cut -d '=' -f3)"
			fi
			sudo -u $user DISPLAY=:0 DBUS_SESSION_BUS_ADDRESS=$display notify-send -i ` + dirPath + `assets/img/logo.png "Aeacus SE" "` + messageString + `"`)
		if err != nil {
			fail("Sending notification failed. Is the user in the configuration correct, and are they logged in to a desktop environment?")
		}
	}
}

func checkTrace() {
	result, err := cond{
		Path:  "/proc/self/status",
		Value: `^TracerPid:\s+0$`,
	}.FileContains()

	// If there was an error reading the file, the user may be restricting access to /proc for the phocus binary
	// through tools such as AppArmor. In this case, the engine should error out.
	if !result || err {
		fail("Try harder instead of ptracing the engine, please.")
		os.Exit(1)
	}
}

// createFQs is a quality of life function that creates Forensic Question files
// on the Desktop, pre-populated with a template.
func CreateFQs(numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > /home/" + conf.User + "/Desktop/" + fileName)
		shellCommand("echo 'ANSWER:' >> /home/" + conf.User + "/Desktop/" + fileName)
		info("Wrote " + fileName + " to Desktop")
	}
}

// rawCmd returns a exec.Command object for Linux shell commands.
func rawCmd(commandGiven string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", commandGiven)
}

// playAudio plays a .wav file with the given path.
func playAudio(wavPath string) {
	info("Playing audio:", wavPath)
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
		fail("Error for checking if running as root: " + err.Error())
		return false
	} else if uid != 0 {
		return false
	}
	return true
}

// destroyImage removes the aeacus directory (to stop scoring) and optionally
// can destroy the entire machine.
func destroyImage() {
	// TODO, ensure this doesn't implode
	fail("Destroying the image is temporarily cancelled.")
	os.Exit(1)
	/*
		fail("Destroying the image!")
		if verboseEnabled {
			warn("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
		} else {
			shellCommand("rm -rf " + dirPath)
			if conf.Destroy {
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
	*/
}

func getInfo(infoType string) {
	warn("Info gathering is not supported for Linux-- there's always a better, easier command-line tool.")
}
