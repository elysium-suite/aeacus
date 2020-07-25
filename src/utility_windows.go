package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gen2brain/beeep"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// readFile (Windows) uses ioutil's ReadFile function and passes the returned
// byte sequence to decodeString.
func readFile(filename string) (string, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return decodeString(string(raw))
}

// decodeString (Windows) attempts to determine the file encoding type
// (typically, UTF-8, UTF-16, or ANSI) and return the appropriately
// encoded string.
func decodeString(fileContent string) (string, error) {
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

	// Make an tranformer that converts MS-Win default to UTF8
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// Decode and print
	decoded, err := ioutil.ReadAll(unicodeReader)
	return string(decoded), err
}

// sendNotification (Windows) employes the beeep library to send notifications
// to the end user.
func sendNotification(mc *metaConfig, messageString string) {
	err := beeep.Notify("Aeacus SE", messageString, mc.DirPath+"assets/logo.png")
	if err != nil {
		failPrint("Notification error: " + err.Error())
	}
}

// rawCmd returns a exec.Command object with the correct PowerShell flags.
//
// rawCmd uses PowerShell's ScriptBlock feature (along with -NoProfile to
// speed things up, as well as some other flags) to run commands on the host
// system and retrieve the return value.
func rawCmd(commandGiven string) *exec.Cmd {
	fmt.Println("[!] Executing a command: ", "powershell.exe", "-NonInteractive", "-NoProfile", "Invoke-Command", "-ScriptBlock", "{ "+commandGiven+" }")
	return exec.Command("powershell.exe", "-NonInteractive", "-NoProfile", "Invoke-Command", "-ScriptBlock", "{ "+commandGiven+" }")
}

// shellCommand (Windows) executes a given command in a PowerShell environment
// and prints an error if one occurred.
func shellCommand(commandGiven string) {
	cmd := rawCmd(commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			if len(commandGiven) > 12 {
				failPrint("Command \"" + commandGiven[:12] + "...\" errored out (code " + err.Error() + ").")
			} else {
				failPrint("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
			}
		}
	}
}

// shellCommand (Windows) executes a given command in a PowerShell environment
// and returns the commands output and its error (if one occurred).
func shellCommandOutput(commandGiven string) (string, error) {
	out, err := rawCmd(commandGiven).Output()
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

func playAudio(wavPath string) {
	commandText := "(New-Object Media.SoundPlayer '" + wavPath + "').PlaySync();"
	shellCommand(commandText)
}

// createFQs is a quality of life function that creates Forensic Question files
// on the Desktop, pre-populated with a template.
func createFQs(mc *metaConfig, numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		shellCommand("echo 'ANSWER:' >> C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		if verboseEnabled {
			infoPrint("Wrote " + fileName + " to Desktop")
		}
	}
}

// adminCheck (Windows) will attempt to open:
//     \\.\PHYSICALDRIVE0
// and will return true if this succeeds, which means the process is running
// as Administrator.
func adminCheck() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func destroyImage(mc *metaConfig) {
	failPrint("Destroying the image!")
	if verboseEnabled {
		warnPrint("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
	} else {
		shellCommand("del /s /q C:\\aeacus")
		if !(mc.Config.NoDestroy == "yes") {
			// nuke registry
			// other destructive commands
			// rm -rf /
			// kill all procceses
			// overwrite system32
			shellCommand("shutdown /r /t 0")
		}
		os.Exit(1)
	}
}

// sidToLocalUser takes an SID as a string and returns a string containing
// the username of the Local User (NTAccount) that it belongs to.
func sidToLocalUser(sid string) string {
	cmdText := "$objSID = New-Object System.Security.Principal.SecurityIdentifier('" + sid + "'); $objUser = $objSID.Translate([System.Security.Principal.NTAccount]); Write-Host $objUser.Value"
	output, err := shellCommandOutput(cmdText)
	if err != nil {
		fmt.Println("yep so err was", err.Error())
	}
	return strings.TrimSpace(output)
}

// localUserToSid takes a username as a string and returns a string containing
// its SID. This is the opposite of sidToLocalUser.
func localUserToSid(userName string) (string, error) {
	return shellCommandOutput("$objUser = New-Object System.Security.Principal.NTAccount('" + userName + "'); $strSID = $objUser.Translate([System.Security.Principal.SecurityIdentifier]); Write-Host $strSID.Value")
}

// getSecedit returns the string value of the secedit.exe command:
//     secedit.exe /export
// which contains security policy options that can't be found in the registry.
func getSecedit() (string, error) {
	return shellCommandOutput("secedit.exe /export /cfg sec.cfg /log NUL; Get-Content sec.cfg; Remove-Item sec.cfg")
}

// getNetUserInfo returns the string output from the command:
//     net user {username}
// in order to get user properties and details.
func getNetUserInfo(userName string) (string, error) {
	return shellCommandOutput("net user " + userName)
}
