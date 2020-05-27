package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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

func createFQs(mc *metaConfig) {
	var numFQ int
	printerPrompt("How many FQs do you want to create? ")
	fmt.Scanln(&numFQ)

	for i := 1; i <= numFQ; i++ {
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

func destroyImage() {
	// ideas for destroying windows
	// nuke registry
	// rm -rf /
	// kill all procceses
	// overwrite system32
	fmt.Println("cant do that yet. not supported on windows. enjoy ur undestryoed imaeg")
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
