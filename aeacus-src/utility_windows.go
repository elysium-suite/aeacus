package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Similar to ioutil.ReadFile() but decodes UTF-16.  Useful when
// reading data from MS-Windows systems that generate UTF-16BE files,
// but will do the right thing if other BOMs are found.
func readFile(filename string) (string, error) {
	// Read the file into a []byte
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return tryDecodeString(string(raw))

}

func tryDecodeString(fileContent string) (string, error) {
	// If contains ~>40% null bytes, we're gonna assume its Unicode
	raw := []byte(fileContent)
	index := bytes.IndexByte(raw, 0)
	if index >= 0 {
		nullCount := 0
		for _, byteChar := range raw {
			if byteChar == 0 {
				nullCount++
			}
		}
		percentNull := float32(nullCount) / float32(len(raw))
		if percentNull < 0.40 {
			return string(raw), nil
		}
	} else {
		return string(raw), nil
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	return string(decoded), err
}

func shellCommand(commandGiven string) {
	cmd := exec.Command("powershell.exe", "-NonInteractive", "-NoProfile", "Invoke-Command", "-ScriptBlock", "{ "+commandGiven+" }")
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
	out, err := exec.Command("powershell.exe", "-NonInteractive", "-NoProfile", "Invoke-Command", "-ScriptBlock", "{ "+commandGiven+" }").Output()
	if err != nil {
		if len(commandGiven) > 9 {
			failPrint("Command \"" + commandGiven[:9] + "...\" errored out (code " + err.Error() + ").")
		} else {
			failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), err
}

func createFQs(mc *metaConfig, numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		shellCommand("echo 'ANSWER:' >> C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		if mc.Cli.Bool("v") {
			infoPrint("Wrote " + fileName + " to Desktop")
		}
	}
}

func playAudio(wavPath string) {
	commandText := "(New-Object Media.SoundPlayer '" + wavPath + "').PlaySync();"
	shellCommand(commandText)
}

func destroyImage(mc *metaConfig) {
	failPrint("Destroying the image!")
	if mc.Cli.Bool("v") {
		warnPrint("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
	} else {
		// ideas for destroying windows
		// nuke registry
		// rm -rf /
		// kill all procceses
		// overwrite system32

	}
}

// sidToLocalUser takes an SID as a string and returns a string containing
// the username of the Local User (NTAccount) that it belongs to
func sidToLocalUser(sid string) string {
	cmdText := "$objSID = New-Object System.Security.Principal.SecurityIdentifier('" + sid + "'); $objUser = $objSID.Translate([System.Security.Principal.NTAccount]); Write-Host $objUser.Value"
	output, err := shellCommandOutput(cmdText)
	if err != nil {
		fmt.Println("yep so err was", err.Error())
	}
	return strings.TrimSpace(output)
}

// localUserToSid takes a username as a string and returns a string containing
// its SID. This is the opposite of sidToLocalUser
func localUserToSid(userName string) (string, error) {
	return shellCommandOutput(fmt.Sprintf("$objUser = New-Object System.Security.Principal.NTAccount('%s'); $strSID = $objUser.Translate([System.Security.Principal.SecurityIdentifier]); Write-Host $strSID.Value", userName))
}

// getSecedit returns the string value of the secedit.exe /export command
// which contains security policy options that can't be found in the registry
func getSecedit() (string, error) {
	return shellCommandOutput("secedit.exe /export /cfg sec.cfg /log NUL; Get-Content sec.cfg; Remove-Item sec.cfg")
}

// getNetUserInfo returns the string output from the command `net user {username}` in order to get user properties and details
func getNetUserInfo(userName string) (string, error) {
	return shellCommandOutput("net user " + userName)
}

// parseCmdOutput takes Windows CMD output of keys in the form `Key Value`, `Key = Value,Value,Value`, and `Key = "Value"` and returns a string map of values and keys
// should really implement this for standardized command output processing
func parseCmdOutput(inputStr string) []string {
	valuePairs := []string{}
	// split inputstr on whitespace
	// parsing loop for each line
	// trimspace every field
	// if equal sign, split on that
	// if comma, split on commas
	// if quotes, remove those
	// else no equal sign
	// assign first to the remainder
	return valuePairs
}
