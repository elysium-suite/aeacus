package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"unicode/utf8"
)

func shellCommand(commandGiven string) {
	cmd := exec.Command("sh", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			failPrint("Command \"" + commandGiven + "\" errored out (" + err.Error() + ").")
		}
	}
}

func shellCommandOutput(commandGiven string) (string, error) {
	out, err := exec.Command("sh", "-c", commandGiven).Output()
	if err != nil {
		failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
		return "", err
	}
	return string(out), err
}

func sendNotification(userName string, notifyText string) {
	commandText := "/sbin/runuser -l " + userName + " -c  '/usr/bin/notify-send -i /opt/aeacus/web/assets/logo.png \"Aeacus Scoring System\" \"" + notifyText + "\"'"
	shellCommand(commandText)
}

func playAudio(wavPath string) {
	commandText := "aplay " + wavPath
	shellCommand(commandText)
}

func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func verifyBinary(binPath string) bool {
	commandText := "which " + binPath
	path, err := shellCommandOutput(commandText)
	binPkg := "dpkg -S" + path + "cut -d':' -f1"
	thing4 = trimFirstRune(path)
	thing, err2 := shellCommandOutput(binPkg)
	thing3 := "grep /var/lib/dpkg/info/" + thing + ".md5sums | grep " + thing4 + " cut -d' ' -f1"
	thing2 := shellCommandOutput(thing3)
	thing5 := hash_file_md5(binPath)
	if thing2 == thing5 {
		return true //the binary is ok
	}
	return false //the binary is not ok
	// sorry if the var names make no sense i wrote this at 1 am
}

func destroyImage() {
	warnPrint("Destroying the image! (jk for now. that's dangerous)")
}
