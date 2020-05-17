package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// adminCheck will return true if process is being run as Administrator
func adminCheck() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// This Windows processCheck will process Windows-specific checks
// handed to it by the processCheckWrapper function
func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "RegistryKey":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " matches \"" + arg2 + "\""
		}
		result, err := RegistryKey(arg1, arg2, false)
		if err != nil {
			return false
		}
		return result
	case "RegistryKeyNot":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " does not match \"" + arg2 + "\""
		}
		result, err := RegistryKey(arg1, arg2, false)
		if err != nil {
			return false
		}
		return !result
	case "RegistryKeyExists":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " exists"
		}
		result, err := RegistryKey(arg1, arg2, true)
		if err != nil {
			return false
		}
		return result
	case "RegistryKeyExistsNot":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " does not exist"
		}
		result, err := RegistryKey(arg1, arg2, true)
		if err != nil {
			return false
		}
		return !result
	case "UserRights":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has privilege \"" + arg2 + "\""
		}
		result, err := UserRights(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "UserRightsNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " does not have privilege \"" + arg2 + "\""
		}
		result, err := UserRights(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "SecurityPolicy":
		if check.Message == "" {
			check.Message = "Security policy option " + arg1 + " is \"" + arg2 + "\""
		}
		result, err := SecurityPolicy(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "SecurityPolicyNot":
		if check.Message == "" {
			check.Message = "Security policy option " + arg1 + " is not \"" + arg2 + "\""
		}
		result, err := SecurityPolicy(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	default:
		failPrint("No check type " + checkType)
	}
	return false
}

func Command(commandGiven string) (bool, error) {
	cmd := exec.Command("powershell.exe", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
}

func PackageInstalled(packageName string) (bool, error) {
	packageList := getPackages()
	for _, p := range packageList {
		if p == packageName {
			return true, nil
		}
	}
	return false, nil
}

func ServiceUp(serviceName string) (bool, error) {
	return Command("(!((Get-Service -Name '" + serviceName + "').Status -eq 'Running')) { Throw 'Service is stopped' }")
}

func UserExists(userName string) (bool, error) {
	// eventually going to not use powershell for everything
	// but until then...
	return Command(fmt.Sprintf("Get-LocalUser %s", userName))
}

func FirewallUp() (bool, error) {
	fwProfiles := []string{"Domain", "Public", "Private"}
	for _, profile := range fwProfiles {
		// This is kind of jank and kind of slow
		cmdText := fmt.Sprintf("if (!((Get-NetFirewallProfile -Name '%s').Enabled -eq 'True')) { Throw 'Firewall profile is disabled' }", profile)
		result, err := Command(cmdText)
		if result == false || err != nil {
			return result, err
		}
	}
	return true, nil
}

func UserRights(userOrGroup string, privilege string) (bool, error) {
	// todo consider /mergedpolicy when windows domain is active?
	// domain support is untested, it should be easy to add a domain
	// flag in the config though. then just make sure you're not getting
	// invalid local policies instead of gpo
	seceditOutput, err := getSecedit()
	if err != nil {
		return false, err
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + privilege + ".*$")
	privilegeString := string(re.Find([]byte(seceditOutput)))
	if privilegeString == "" {
		return false, errors.New("Invalid privilege")
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

func SecurityPolicy(keyName string, keyValue string) (bool, error) {
	var desiredString string
	if regKey, ok := secpolToKey[keyName]; ok {
		return RegistryKey(regKey, keyValue, false)
	} else {
		// Yes, this is jank, but is there a better way? Probably
		output, err := getSecedit()
		if err != nil {
			return false, err
		}
		if keyName == "NewAdministratorName" || keyName == "NewGuestName" {
			// These two are strings, not numbers, so they have ""
			desiredString = fmt.Sprintf("%s = \"%s\"", keyName, keyValue)
		} else {
			desiredString = fmt.Sprintf("%s = %s", keyName, keyValue)
		}
		return strings.Contains(output, desiredString), err
	}
}

func RegistryKey(keyName string, keyValue string, existCheck bool) (bool, error) {

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
			failPrint("Registry opening key failed: " + err.Error())
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
		// for RegistryKey or RegistryKeyNot, so we return an error
		if existCheck {
			return false, nil
		} else {
			failPrint("Registry opening key failed: " + err.Error())
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
