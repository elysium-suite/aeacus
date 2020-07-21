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
	default:
		failPrint("No check type " + checkType)
	}
	return false
}

func command(commandGiven string) (bool, error) {
	cmd := rawCmd(commandGiven + "; if (!($?)) { Throw 'Error' } }")
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
}

func packageInstalled(packageName string) (bool, error) {
	packageList := getPackages()
	for _, p := range packageList {
		if p == packageName {
			return true, nil
		}
	}
	return false, nil
}

func serviceUp(serviceName string) (bool, error) {
	return command(fmt.Sprintf("if (!((Get-Service -Name '%s').Status -eq 'Running')) { Throw 'Service is stopped' }", serviceName))
}

func userExists(userName string) (bool, error) {
	// eventually going to not use powershell for everything
	// but until then...
	return command(fmt.Sprintf("Get-LocalUser '%s'", userName))
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
		cmdText := fmt.Sprintf("if (!((Get-NetFirewallProfile -Name '%s').Enabled -eq 'True')) { Throw 'Firewall profile is disabled' }", profile)
		result, err := command(cmdText)
		if result == false || err != nil {
			return result, err
		}
	}
	return true, nil
}

func userDetail(userName string, detailName string, detailValue string) (bool, error) {
	if userName == "" || detailName == "" || detailValue == "" {
		failPrint("Invalid parameters to UserDetail check")
		return false, errors.New("Invalid parameters")
	}
	userInfo, err := getNetUserInfo(userName)
	if err != nil {
		return false, err
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + detailName + ".*$")
	detailString := string(re.Find([]byte(userInfo)))
	if detailString == "" {
		return false, errors.New("Invalid user detail name")
	}
	re = regexp.MustCompile("[\\s]+")
	userInfoSlice := re.Split(detailString, -1)
	if len(userInfoSlice) < 2 {
		failPrint("Error splitting user detail string into two or more parts")
		return false, errors.New("Error splitting detail string")
	}
	indexValue := len(re.Split(detailName, -1)) + 1
	if indexValue >= len(userInfoSlice) {
		failPrint("Error in calculating index for detailValue")
		return false, errors.New("Error splitting detailName")
	}
	userDetailValue := strings.TrimSpace(userInfoSlice[indexValue])
	//fmt.Println("is", detailValue, "equal to", userDetailValue)
	if detailValue == userDetailValue {
		return true, nil
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
	return command(fmt.Sprintf("Get-SmbShare -Name '%s'", shareName))
}

func scheduledTaskExists(taskName string) (bool, error) {
	return command(fmt.Sprintf("Get-ScheduledTask -TaskName '%s'", taskName))
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
			desiredString = fmt.Sprintf("%s = \"%s\"", keyName, keyValue)
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
				desiredString = fmt.Sprintf("%s = %d", keyName, c)
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
				desiredString = fmt.Sprintf("%s = %d", keyName, c)
				if strings.Contains(output, desiredString) {
					return true, err
				}
			}
		} else {
			desiredString = fmt.Sprintf("%s = %s", keyName, keyValue)
		}
		return strings.Contains(output, desiredString), err
	}
}

func registryKey(keyName string, keyValue string, existCheck bool) (bool, error) {

	// Break down input
	registryArgs := regexp.MustCompile("[\\\\]+").Split(keyName, -1)
	registryHiveText := registryArgs[0]
	keyPath := fmt.Sprintf(strings.Join(registryArgs[1:len(registryArgs)-1], "\\"))
	keyLoc := registryArgs[len(registryArgs)-1]
	//fmt.Printf("REGISTRY: getting keypath %s from %s\n", keyPath, registryHiveText)

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
	//fmt.Printf("Retrieved registry value was %d (length %d, type %d)\n", registrySlice, regLength, valType)

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

	//fmt.Printf("Registry value: %s, keyvalue %s\n", registryValue, keyValue)
	if registryValue == keyValue {
		return true, err
	}
	return false, err
}
