package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// processCheck (Windows) will process Windows-specific checks handed to it
// by the processCheckWrapper function.
func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
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
		result, err := PasswordChanged(arg1, arg2)
		return err == nil && result
	case "PasswordChangedNot":
		if check.Message == "" {
			check.Message = "Password for " + arg1 + " has not been changed"
		}
		result, err := PasswordChanged(arg1, arg2)
		return err == nil && !result
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
	case "FirefoxPrefIs":
		if check.Message == "" {
			check.Message = "Firefox preference " + arg1 + " is set to " + arg2
		}
		result, err := firefoxSetting(arg1, arg2)
		return err == nil && result
	case "FirefoxPrefIsNot":
		if check.Message == "" {
			check.Message = "Firefox preference " + arg1 + " is not set to " + arg2
		}
		result, err := firefoxSetting(arg1, arg2)
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

func commandOutput(commandGiven, desiredOutput string) (bool, error) {
	out, err := rawCmd(commandGiven).Output()
	if err != nil {
		return false, err
	}
	outString := strings.TrimSpace(string(out))
	if outString == desiredOutput {
		return true, nil
	}
	return false, nil
}

func packageInstalled(packageName string) (bool, error) {
	packageList, err := getPackages()
	if err != nil {
		return false, err
	}
	for _, p := range packageList {
		if p == packageName {
			return true, nil
		}
	}
	return false, nil
}

func serviceUp(serviceName string) (bool, error) {
	return commandOutput("(Get-Service -Name '"+serviceName+"').Status", "Running")
}

func PasswordChanged(user, date string) (bool, error) {
	return command(`Get-LocalUser " + user + " | select PasswordLastSet | Select-String "` + date + `"`)
}

func windowsFeature(feature string) (bool, error) {
	return commandOutput("(Get-WindowsOptionalFeature -FeatureName "+feature+" -Online).State", "Enabled")
}

func fileOwner(filePath, owner string) (bool, error) {
	return commandOutput("(Get-Acl "+filePath+").Owner", owner)
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

func userInGroup(userName string, groupName string) (bool, error) {
	userInfo, err := getNetUserInfo(userName)
	if err != nil {
		return false, err
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*Group.*$")
	detailString := strings.TrimSpace(string(re.Find([]byte(userInfo))))
	if detailString == "" {
		// This is likely because an invalid user was tested.
		failPrint("Group check output empty-- please ensure you entered a valid user.")
		return false, errors.New("Error parsing net user output for Group")
	}
	return strings.Contains(detailString, groupName), nil
}

func firewallUp() (bool, error) {
	fwProfiles := []string{"Domain", "Public", "Private"}
	for _, profile := range fwProfiles {
		// This is kind of jank and kind of slow
		cmdText := "(Get-NetFirewallProfile -Name '" + profile + "').Enabled"
		result, err := commandOutput(cmdText, "True")
		if result == false || err != nil {
			return result, err
		}
	}
	return true, nil
}

func userDetail(userName string, detailName string, detailValue string) (bool, error) {
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

func userRights(userOrGroup string, privilege string) (bool, error) {
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
		return false, errors.New("Invalid privilege")
	}
	if strings.Contains(privilegeString, userOrGroup) {
		// Sometimes, Windows just puts their user or group name instead of the SID. Real cool
		return true, nil
	}
	privStringSplit := strings.Split(privilegeString, " ")
	if len(privStringSplit) != 3 {
		return false, errors.New("Error splitting privilege")
	}
	privStringSplit = strings.Split(string(privStringSplit[2]), ",")
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
	return true, nil
}

func securityPolicy(keyName string, keyValue string) (bool, error) {
	var desiredString string
	if regKey, ok := secpolToKey[keyName]; ok {
		return registryKey(regKey, keyValue, false)
	} else {
		output, err := getSecedit()
		if err != nil {
			return false, err
		}
		if keyName == "NewAdministratorName" || keyName == "NewGuestName" {
			// These two are strings, not numbers, so they have ""
			desiredString = keyName + " = " + keyValue
		} else if keyName == "MinimumPasswordAge" ||
			keyName == "MinimumPasswordAge" ||
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
				if strings.Contains(output, desiredString) {
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
				if strings.Contains(output, desiredString) {
					return true, nil
				}
			}
		} else {
			desiredString = keyName + " = " + keyValue
		}
		return strings.Contains(output, desiredString), nil
	}
}

func registryKey(keyName string, keyValue string, existCheck bool) (bool, error) {
	// Break down input
	registryArgs := regexp.MustCompile("[\\\\]+").Split(keyName, -1)
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
		} else {
			failPrint("Unknown registry hive: " + registryHiveText)
			return false, errors.New("Unknown registry hive" + registryHiveText)
		}
	}

	// Actually get the key
	k, err := registry.OpenKey(registryHive, keyPath, registry.QUERY_VALUE)
	if err != nil {
		if existCheck {
			return false, nil
		} else {
			warnPrint("Registry opening key failed (and that's probably fine): " + err.Error())
			return false, err
		}
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
		} else {
			warnPrint("Registry opening key failed (and that's probably fine): " + err.Error())
			return false, err
		}
	} else {
		if existCheck {
			return true, nil
		}
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
		failPrint("Unknown registry type: " + string(valType))
	}

	// fmt.Printf("Registry value: %s, keyvalue %s\n", registryValue, keyValue)
	if registryValue == keyValue {
		return true, err
	}
	return false, err
}

func firefoxSetting(param, value string) (bool, error) {
	res := false
	var err error
	// Check firefox install dir cus that may change where settings are located
	bit64, _ := pathExists(`C:\Program Files\Mozilla Firefox`)
	bit32, _ := pathExists(`C:\Program Files (x86)\Mozilla Firefox`)
	if bit64 {
		check, err := dirContainsRegex(`C:\Program Files\Mozilla Firefox\defaults\pref`, `pref("general.config.filename"`)
		if err != nil {
			return res, err
		}

		if check {
			res, err = dirContainsRegex(`C:\Program Files\Mozilla Firefox`, `("`+param+`",`+value+`)`)
		} else {
			res, err = dirContainsRegex(`C:\Users\`+mc.Config.User+`\AppData\Roaming\Mozilla\Firefox\Profiles`, `("`+param+`",`+value+`)`)
		}

	} else if bit32 {
		check, err := dirContainsRegex(`C:\Program Files (x86)\Mozilla Firefox\defaults\pref`, `pref("general.config.filename"`)
		if err != nil {
			return res, err
		}

		if check {
			res, err = dirContainsRegex(`C:\Program Files (x86)\Mozilla Firefox`, `("`+param+`",`+value+`)`)
		} else {
			res, err = dirContainsRegex(`C:\Users\`+mc.Config.User+`\AppData\Roaming\Mozilla\Firefox\Profiles`, `("`+param+`",`+value+`)`)
		}

	} else {
		err = errors.New("Firefox was not detected on the system")
	}

	return res, err
}
