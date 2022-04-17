package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"

	"github.com/DataDog/datadog-agent/pkg/util/winutil"
	"github.com/gen2brain/beeep"
	wapi "github.com/iamacarpet/go-win64api"
	"github.com/iamacarpet/go-win64api/shared"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	kernel32DLL   = windows.NewLazyDLL("Kernel32.dll")
	debuggerCheck = kernel32DLL.NewProc("IsDebuggerPresent")
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

func checkTrace() {
	result, _, _ := debuggerCheck.Call()
	if int(result) != 0 {
		failPrint("Reversing is cool, but we would appreciate if you practiced your skills in an environment that was less destructive to other peoples' experiences.")
		os.Exit(1)
	}
}

// sendNotification (Windows) employes the beeep library to send notifications
// to the end user.
func sendNotification(messageString string) {
	err := beeep.Notify("Aeacus SE", messageString, mc.DirPath+"assets/img/logo.png")
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
	cmdInput := "powershell.exe -NonInteractive -NoProfile Invoke-Command -ScriptBlock { " + commandGiven + " }"
	debugPrint("rawCmd input: " + cmdInput)
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

// CreateFQs is a quality of life function that creates Forensic Question files
// on the Desktop, pre-populated with a template.
func CreateFQs(numFqs int) {
	for i := 1; i <= numFqs; i++ {
		fileName := "'Forensic Question " + strconv.Itoa(i) + ".txt'"
		shellCommand("echo 'QUESTION:' > C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		shellCommand("echo 'ANSWER:' >> C:\\Users\\" + mc.Config.User + "\\Desktop\\" + fileName)
		infoPrint("Wrote " + fileName + " to Desktop")
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

func destroyImage() {
	failPrint("Destroying the image!")
	if verboseEnabled {
		warnPrint("Since you're running this in verbose mode, I assume you're a developer who messed something up. You've been spared from image deletion but please be careful.")
	} else {
		shellCommand("del /s /q C:\\aeacus")
		if !mc.Config.NoDestroy {
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
	output, _ := shellCommandOutput(cmdText)
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

// getPrograms returns a list of currently installed Programs and
// their versions.
func getPrograms() ([]string, error) {
	softwareList := []string{}
	sw, err := wapi.InstalledSoftwareList()
	if err != nil {
		failPrint("Couldn't get programs: " + err.Error())
		return softwareList, err
	}
	for _, s := range sw {
		softwareList = append(softwareList, s.Name()+" - version "+s.DisplayVersion)
	}
	return softwareList, nil
}

// getProgram returns the Software struct of program data from a name.
// The first Program that contains the substring passed as the programName
// is returned.
func getProgram(programName string) (shared.Software, error) {
	prog := shared.Software{}
	sw, err := wapi.InstalledSoftwareList()
	if err != nil {
		failPrint("Couldn't get programs: " + err.Error())
	}
	for _, s := range sw {
		if strings.Contains(s.Name(), programName) {
			return s, nil
		}
	}
	return prog, errors.New("program not found")
}

func getLocalUsers() ([]shared.LocalUser, error) {
	ul, err := wapi.ListLocalUsers()
	if err != nil {
		failPrint("Couldn't get local users: " + err.Error())
	}
	return ul, err
}

func getLocalAdmins() ([]shared.LocalUser, error) {
	ul, err := wapi.ListLocalUsers()
	if err != nil {
		failPrint("Couldn't get local users: " + err.Error())
	}
	var admins []shared.LocalUser
	for _, user := range ul {
		if user.IsAdmin {
			admins = append(admins, user)
		}
	}
	return admins, err
}

func getLocalUser(userName string) (shared.LocalUser, error) {
	userList, err := getLocalUsers()
	if err != nil {
		return shared.LocalUser{}, err
	}
	for _, user := range userList {
		if user.Username == userName {
			return user, nil
		}
	}
	return shared.LocalUser{}, nil
}

func getLocalServiceStatus(serviceName string) (shared.Service, error) {
	serviceDataList, err := wapi.GetServices()
	var serviceStatusData shared.Service
	if err != nil {
		failPrint("Couldn't get local service: " + err.Error())
		return serviceStatusData, err
	}
	for _, v := range serviceDataList {
		if v.SCName == serviceName {
			return v, nil
		}
	}
	failPrint(`Specified service '` + serviceName + `' was not found on the system`)
	return serviceStatusData, err
}

func getFileAccess(mask uint32) string {
	var ret []string

	if mask == 0x1f01ff {
		ret = append(ret, "FullControl")
	} else if mask == 0x1301bf {
		ret = append(ret, "Modify")
	} else {
		if mask&windows.FILE_WRITE_ATTRIBUTES != 0 {
			ret = append(ret, "Write")
		}
		if mask&windows.FILE_READ_ATTRIBUTES != 0 {
			if mask&0x000020 != 0 {
				ret = append(ret, "ReadAndExecute")
			} else {
				ret = append(ret, "Read")
			}
		}
	}
	return strings.Join(ret, ", ")
}

func getFileRights(filePath, username string) (map[string]string, error) {

	userSID, _, _, err := windows.LookupSID("", username)

	if err != nil {
		return nil, errors.Wrapf(err, "acl: failed to lookup sid for user '%s'", username)
	}

	ret := map[string]string{}

	var fileDACL *winutil.Acl
	if err := winutil.GetNamedSecurityInfo(filePath,
		winutil.SE_FILE_OBJECT,
		winutil.DACL_SECURITY_INFORMATION,
		nil,
		nil,
		&fileDACL,
		nil,
		nil); err != nil {
		return nil, errors.Wrapf(err, "acl: failed to get security info for '%s' ", filePath)
	}

	var aclSizeInfo winutil.AclSizeInformation
	if err := winutil.GetAclInformation(fileDACL, &aclSizeInfo, winutil.AclSizeInformationEnum); err != nil {
		return nil, errors.Wrapf(err, "acl: failed to get acl size info for '%s' ", filePath)
	}

	for i := uint32(0); i < aclSizeInfo.AceCount; i++ {
		var pACE *winutil.AccessAllowedAce
		if err := winutil.GetAce(fileDACL, i, &pACE); err != nil {
			return nil, errors.Wrapf(err, "acl: failed to get acl for '%s' ", filePath)
		}
		// update
		sid := (*windows.SID)(unsafe.Pointer(&pACE.SidStart))
		if strings.EqualFold(userSID.String(), sid.String()) {
			ret["identityreference"] = username
			if pACE.AceType == winutil.ACCESS_DENIED_ACE_TYPE {
				ret["accesscontroltype"] = "Deny"
			}
			if pACE.AceType == winutil.ACCESS_ALLOWED_ACE_TYPE {
				ret["accesscontroltype"] = "Allow"
			}
			ret["filesystemrights"] = getFileAccess(pACE.AccessMask)
		}
	}
	return ret, nil
}
