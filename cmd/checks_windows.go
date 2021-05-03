package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	wapi "github.com/iamacarpet/go-win64api"
	"golang.org/x/sys/windows/registry"
)

// processCheck (Windows) will process Windows-specific checks handed to it
// by the processCheckWrapper function.
func processCheck(check *check, checkType, arg1, arg2, arg3 string) bool {
	switch checkType {
	case "UserDetail":
		if check.Message == "" {
			check.Message = "User property " + arg2 + " for " + arg1 + " is equal to \"" + arg3 + "\""
		}
		result, err := userDetail(arg1, arg2, arg3)
		return err == nil && result
	case "UserDetailNot":
		if check.Message == "" {
			check.Message = "User property " + arg2 + " for " + arg1 + " is not equal to \"" + arg3 + "\""
		}
		result, err := userDetail(arg1, arg2, arg3)
		return err == nil && !result
	case "UserRights":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has privilege \"" + arg2 + "\""
		}
		result, err := userRights(arg1, arg2)
		return err == nil && result
	case "UserRightsNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " does not have privilege \"" + arg2 + "\""
		}
		result, err := userRights(arg1, arg2)
		return err == nil && !result
	case "ShareExists":
		if check.Message == "" {
			check.Message = "Share " + arg1 + " exists"
		}
		result, err := shareExists(arg1)
		return err == nil && result
	case "ShareExistsNot":
		if check.Message == "" {
			check.Message = "Share " + arg1 + " doesn't exist"
		}
		result, err := shareExists(arg1)
		return err == nil && !result
	case "ScheduledTaskExists":
		if check.Message == "" {
			check.Message = "Scheduled task " + arg1 + " exists"
		}
		result, err := scheduledTaskExists(arg1)
		return err == nil && result
	case "ScheduledTaskExistsNot":
		if check.Message == "" {
			check.Message = "Scheduled task " + arg1 + " doesn't exist"
		}
		result, err := scheduledTaskExists(arg1)
		return err == nil && !result
	/*
		case "StartupProgramExists":
			if check.Message == "" {
				check.Message = "Startup program " + arg1 + " exists"
			}
			result, err := startupProgramExists(arg1)
			return err == nil && result
		case "StartupProgramExistsNot":
			if check.Message == "" {
				check.Message = "Startup program " + arg1 + " doesn't exist"
			}
			result, err := scheduledTaskExists(arg1)
			return err == nil && !result
	*/
	case "SecurityPolicy":
		if check.Message == "" {
			check.Message = "Security policy option " + arg1 + " is \"" + arg2 + "\""
		}
		result, err := securityPolicy(arg1, arg2)
		return err == nil && result
	case "SecurityPolicyNot":
		if check.Message == "" {
			check.Message = "Security policy option " + arg1 + " is not \"" + arg2 + "\""
		}
		result, err := securityPolicy(arg1, arg2)
		return err == nil && !result
	case "RegistryKey":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " matches \"" + arg2 + "\""
		}
		result, err := registryKey(arg1, arg2, false)
		return err == nil && result
	case "RegistryKeyNot":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " does not match \"" + arg2 + "\""
		}
		result, err := registryKey(arg1, arg2, false)
		return err == nil && !result
	case "RegistryKeyExists":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " exists"
		}
		result, err := registryKey(arg1, arg2, true)
		return err == nil && result
	case "RegistryKeyExistsNot":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " does not exist"
		}
		result, err := registryKey(arg1, arg2, true)
		return err == nil && !result
	case "PasswordChanged":
		if check.Message == "" {
			check.Message = "Password for " + arg1 + " has been changed"
		}
		result, err := passwordChanged(arg1, arg2)
		return err == nil && !result
	case "PasswordChangedNot":
		if check.Message == "" {
			check.Message = "Password for " + arg1 + " has not been changed"
		}
		result, err := passwordChanged(arg1, arg2)
		return err == nil && result
	case "WindowsFeature":
		if check.Message == "" {
			check.Message = arg1 + " feature has been enabled"
		}
		result, err := windowsFeature(arg1)
		return err == nil && result
	case "WindowsFeatureNot":
		if check.Message == "" {
			check.Message = arg1 + " feature has been disabled"
		}
		result, err := windowsFeature(arg1)
		return err == nil && !result
	case "FileOwner":
		if check.Message == "" {
			check.Message = arg1 + " is owned by " + arg2
		}
		result, err := fileOwner(arg1, arg2)
		return err == nil && result
	case "FileOwnerNot":
		if check.Message == "" {
			check.Message = arg1 + " is not owned by " + arg2
		}
		result, err := fileOwner(arg1, arg2)
		return err == nil && !result
	case "ServiceStatus":
		if check.Message == "" {
			check.Message = "The service " + arg1 + " is " + arg2 + " with the startup type set as " + arg3
		}
		result, err := serviceStatus(arg1, arg2, arg3)
		return err == nil && result
	case "ServiceStatusNot":
		if check.Message == "" {
			check.Message = "The service " + arg1 + " is not " + arg2 + " with the startup type not set as " + arg3
		}
		result, err := serviceStatus(arg1, arg2, arg3)
		return err == nil && !result
	default:
		failPrint("No check type " + checkType)
	}
	return false
}

func command(commandGiven string) (bool, error) {
	cmd := rawCmd(commandGiven + "; if (!($?)) { Throw 'Error' }")
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
}

func programInstalled(programName string) (bool, error) {
	programList, err := getPrograms()
	if err != nil {
		return false, err
	}
	for _, p := range programList {
		if strings.Contains(p, programName) {
			return true, nil
		}
	}
	return false, nil
}

func programVersion(programName, versionNum, compareMode string) (bool, error) {
	prog, err := getProgram(programName)
	if err != nil {
		return false, err
	}
	switch compareMode {
	case "eq":
		if prog.DisplayVersion == versionNum {
			return true, nil
		}
	case "gt":
		if prog.DisplayVersion > versionNum {
			return true, nil
		}
	case "ge":
		if prog.DisplayVersion >= versionNum {
			return true, nil
		}
	}
	return false, nil
}

func serviceUp(serviceName string) (bool, error) {
	serviceStatus, err := getLocalServiceStatus(serviceName)
	return serviceStatus.IsRunning, err
}

func serviceStatus(serviceName, wantedStatus, startupType string) (bool, error) {
	status, err := getLocalServiceStatus(serviceName)
	var boolWantedStatus bool
	if err != nil {
		return false, err
	}
	switch wantedStatus = strings.ToLower(wantedStatus); wantedStatus {
	case "running":
		boolWantedStatus = true
	case "stopped":
		boolWantedStatus = false
	default:
		errMessage := "Unknown status type found for " + serviceName
		failPrint(errMessage)
		return false, errors.New(errMessage)
	}
	if status.IsRunning == boolWantedStatus {
		serviceKey := `HKLM\SYSTEM\CurrentControlSet\Services\` + serviceName + `\Start`
		var wantedStartupTypeNumber string
		switch startupType = strings.ToLower(startupType); startupType {
		case "automatic":
			wantedStartupTypeNumber = "2"
		case "manual":
			wantedStartupTypeNumber = "3"
		case "disabled":
			wantedStartupTypeNumber = "4"
		default:
			failPrint("Unknown startup type found for " + serviceName)
			return false, errors.New("Unknown status type found for " + serviceName)
		}
		check, err := registryKey(serviceKey, wantedStartupTypeNumber, false)
		if err != nil {
			return false, err
		}
		if check {
			return true, nil
		}
	}
	return false, err
}

func passwordChanged(user, date string) (bool, error) {
	changed, _ := commandOutput(`(Get-LocalUser ` + user + ` | select PasswordLastSet).PasswordLastSet -replace "n",", " -replace "r",", "`)
	return changed >= date, nil
}

func windowsFeature(feature string) (bool, error) {
	state, _ := commandOutput("(Get-WindowsOptionalFeature -FeatureName " + feature + " -Online).State")
	return state == "Enabled", nil
}

func fileOwner(filePath, owner string) (bool, error) {
	theowner, _ := commandOutput("(Get-Acl " + filePath + ").Owner")
	return theowner == owner, nil
}

func userExists(userName string) (bool, error) {
	user, err := getLocalUser(userName)
	if err != nil {
		return false, err
	}
	if user.Username == "" {
		return false, nil
	}
	return true, nil
}

func userInGroup(userName, groupName string) (bool, error) {
	users, err := wapi.LocalGroupGetMembers(groupName)
	if err != nil {
		// Error is returned if group is empty.
		return false, nil
	}
	for _, user := range users {
		justName := strings.Split(user.Name, `\`)[1]
		if userName == user.Name || userName == justName {
			return true, nil
		}
	}
	return false, nil
}

func firewallUp() (bool, error) {
	fwProfiles := []string{"Domain", "Public", "Private"}
	for _, profile := range fwProfiles {
		// This is kind of jank and kind of slow
		cmdText := "(Get-NetFirewallProfile -Name '" + profile + "').Enabled"
		result, err := commandOutput(cmdText)
		if result != "True" || err != nil {
			return false, err
		}
	}
	return true, nil
}

func userDetail(userName, detailName, detailValue string) (bool, error) {
	detailValue = strings.TrimSpace(detailValue)
	lookingFor := false
	if strings.ToLower(detailValue) == "yes" {
		lookingFor = true
	}
	user, err := getLocalUser(userName)
	if err != nil {
		return false, err
	}
	switch detailName {
	case "FullName":
		if user.FullName == detailValue {
			return true, nil
		}
	case "IsEnabled":
		return user.IsEnabled == lookingFor, nil
	case "IsLocked":
		return user.IsLocked == lookingFor, nil
	case "IsAdmin":
		return user.IsAdmin == lookingFor, nil
	case "PasswordNeverExpires":
		return user.PasswordNeverExpires == lookingFor, nil
	default:
		failPrint("detailName (" + detailName + ") passed to userDetail is invalid.")
		return false, errors.New("Invalid detailName")
	}
	return false, nil
}

func userRights(userOrGroup, privilege string) (bool, error) {
	// todo consider /mergedpolicy when windows domain is active?
	// domain support is untested, it should be easy to add a domain
	// flag in the config though. then just make sure you're not getting
	// invalid local policies instead of gpo

	seceditOutput, err := getSecedit()
	// TODO: only get section of users -- this can also falsely score correct for other secedit fields (like LegalNoticeText)
	if err != nil {
		return false, err
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + privilege + ".*$")
	privilegeString := string(re.Find([]byte(seceditOutput)))
	if privilegeString == "" {
		return false, nil
	}
	if strings.Contains(privilegeString, userOrGroup) {
		// Sometimes, Windows just puts their user or group name instead of the SID. Real cool
		return true, nil
	}
	privStringSplit := strings.Split(privilegeString, " ")
	if len(privStringSplit) != 3 {
		return false, errors.New("Error splitting privilege")
	}
	privStringSplit = strings.Split(privStringSplit[2], ",")
	for _, sidValue := range privStringSplit {
		sidValue = strings.TrimSpace(sidValue)
		userForSid := strings.Split(sidToLocalUser(sidValue[1:]), "\\")
		userSid := strings.TrimSpace(userForSid[0])
		if len(userForSid) == 2 {
			userSid = strings.TrimSpace(userForSid[1])
		}
		if userSid == userOrGroup {
			return true, nil
		}
	}
	return false, err
}

func shareExists(shareName string) (bool, error) {
	return command("Get-SmbShare -Name '" + shareName + "'")
}

func scheduledTaskExists(taskName string) (bool, error) {
	return command("Get-ScheduledTask -TaskName '" + taskName + "'")
}

func startupProgramExists(progName string) (bool, error) {
	// need to work out the implementation on this one too...
	// multiple startup locations
	// rot
	return true, nil
}

func securityPolicy(keyName, keyValue string) (bool, error) {
	var desiredString string
	if regKey, ok := secpolToKey[keyName]; ok {
		return registryKey(regKey, keyValue, false)
	}
	seceditOutput, err := getSecedit()
	if err != nil {
		return false, err
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + keyName + ".*$")
	output := strings.TrimSpace(string(re.Find([]byte(seceditOutput))))
	if output == "" {
		return false, errors.New("securitypolicy item not found")
	}
	if keyName == "NewAdministratorName" || keyName == "NewGuestName" {
		// These two are strings, not numbers, so they have ""
		desiredString = keyName + " = " + keyValue
	} else if keyName == "MinimumPasswordAge" ||
		keyName == "MinimumPasswordLength" ||
		keyName == "LockoutDuration" ||
		keyName == "ResetLockoutCount" {
		// Fields where the arg should be X or higher (up to 999)
		intKeyValue, err := strconv.Atoi(keyValue)
		if err != nil {
			failPrint(keyValue + " is not a valid integer for SecurityPolicy check")
			return false, errors.New("Invalid keyValue")
		}
		for c := intKeyValue; c <= 999; c++ {
			desiredString = keyName + " = " + strconv.Itoa(c)
			if output == desiredString {
				return true, err
			}
		}
	} else if keyName == "MaximumPasswordAge" || keyName == "LockoutBadCount" {
		// Fields where arg should be X or lower but NOT 0
		intKeyValue, err := strconv.Atoi(keyValue)
		if err != nil {
			failPrint(keyValue + " is not a valid integer for SecurityPolicy check")
			return false, errors.New("Invalid keyValue")
		}
		for c := intKeyValue; c > 0; c-- {
			desiredString = keyName + " = " + strconv.Itoa(c)
			if output == desiredString {
				return true, nil
			}
		}
	} else {
		desiredString = keyName + " = " + keyValue
	}
	return output == desiredString, nil
}

func registryKey(keyName, keyValue string, existCheck bool) (bool, error) {
	// Break down input
	registryArgs := regexp.MustCompile(`[\\]+`).Split(keyName, -1)
	registryHiveText := registryArgs[0]
	keyPath := fmt.Sprintf(strings.Join(registryArgs[1:len(registryArgs)-1], "\\")) // idk??
	keyLoc := registryArgs[len(registryArgs)-1]
	// fmt.Printf("REGISTRY: getting keypath %s from %s\n", keyPath, registryHiveText)

	var registryHive registry.Key
	switch registryHiveText {
	case "HKEY_CLASSES_ROOT", "HKCR":
		registryHive = registry.CLASSES_ROOT
	case "HKEY_CURRENT_USER", "HKCU":
		registryHive = registry.CURRENT_USER
	case "HKEY_LOCAL_MACHINE", "HKLM", "MACHINE":
		registryHive = registry.LOCAL_MACHINE
	case "HKEY_USERS", "HKU":
		registryHive = registry.USERS
	case "HKEY_CURRENT_CONFIG", "HKCC":
		registryHive = registry.CURRENT_CONFIG
	case "SOFTWARE":
		registryHive = registry.LOCAL_MACHINE
		keyPath = "SOFTWARE\\" + keyPath
	default:
		if existCheck {
			return false, nil
		}
		failPrint("Unknown registry hive: " + registryHiveText)
		return false, errors.New("Unknown registry hive" + registryHiveText)
	}

	// Actually get the key
	k, err := registry.OpenKey(registryHive, keyPath, registry.QUERY_VALUE)
	if err != nil {
		if existCheck {
			return false, nil
		}
		warnPrint("Registry opening key failed (and that's probably fine): " + err.Error())
		return false, err
	}
	defer k.Close()

	// Fetch registry value
	registrySlice := make([]byte, 256)
	regLength, valType, err := k.GetValue(keyLoc, registrySlice)
	if err != nil {
		// Error is probably about the key not existing.
		// This is fine, some keys are not defined until the setting
		// is explicitly set. However, the check should not pass
		// for RegistryKey or RegistryKeyNot, so we return an error.
		if existCheck {
			return false, nil
		}
		warnPrint("Registry opening key failed (and that's probably fine): " + err.Error())
		return false, err
	}
	if existCheck {
		return true, nil
	}

	registrySlice = registrySlice[:regLength]
	// fmt.Printf("Retrieved registry value was %d (length %d, type %d)\n", registrySlice, regLength, valType)

	// Determine value type to convert to string
	var registryValue string
	switch valType {
	case 1: // SZ
		registryValue, _, err = k.GetStringValue(keyLoc)
	case 2: // EXPAND_SZ
		registryValue, _, err = k.GetStringValue(keyLoc)
	case 3: // BINARY
		failPrint("Binary registry format not yet supported.")
	case 4: // DWORD
		registryValue = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(registrySlice)), 10)
	default:
		failPrint("Unknown registry type: " + fmt.Sprint(valType))
	}

	// fmt.Printf("Registry value: %s, keyvalue %s\n", registryValue, keyValue)
	if registryValue == keyValue {
		return true, err
	}
	return false, err
}
