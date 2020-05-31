package main

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"unicode/utf8"
)

func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	return string(fileContent), err
}

func tryDecodeString(fileContent string) (string, error) {
	// For compatability with Windows ANSI/UNICODE/etcetc
	// and if Linux ever decides to use weird encoding
	return fileContent, nil
}

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

func playAudio(wavPath string) {
	commandText := "aplay " + wavPath
	shellCommand(commandText)
}

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

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func verifyBinary(binName string) bool {
	// function is untested
	// TODO
	binPath, err := shellCommandOutput("which " + binName)
	if err != nil {
		return false
	}
	binPkg := "dpkg -S" + binPath + "cut -d':' -f1"
	binPath = trimFirstRune(binPath)
	binPkg, err = shellCommandOutput(binPkg)
	if err != nil {
		return false
	}
	binHash := "grep /var/lib/dpkg/info/" + binPkg + ".md5sums | grep " + binName + " cut -d' ' -f1"
	binHash, err = shellCommandOutput(binHash)
	if err != nil {
		return false
	}
	binHashExpected, err := hashFileMD5(binPath)
	if err == nil && binHash == binHashExpected {
		return true
	}
	return false
}

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

func destroyImage(mc *metaConfig) {
	failPrint("Destroying the image!")
	if mc.Cli.Bool("v") {
		warnPrint("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
	} else {
		// destroying ideas: start a classic rm -rf in the background, delete /etc/passwd, start a rm -rf in foreground for /etc/, then /bin/, then /home/, then kill -9 all processes, then rm -rf foreground everything else
	}
}
