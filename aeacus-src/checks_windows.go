package main

import (
	"os"
    "fmt"
    "errors"
	"regexp"
    "strings"
    "os/exec"
    "strconv"
    "encoding/binary"

    "golang.org/x/sys/windows/registry"
)

func adminCheck() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "RegistryKey":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " matches \"" + arg2 + "\""
		}
		result, err := RegistryKey(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "RegistryKeyNot":
		if check.Message == "" {
			check.Message = "Registry key " + arg1 + " does not match \"" + arg2 + "\""
		}
		result, err := RegistryKey(arg1, arg2)
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
	// not super happy with the command implementation
	// could just keylog sh or replace dpkg binary or something
	// should use golang dpkg library if it existed and was good
	result, err := Command(fmt.Sprintf("dpkg -l %s", packageName))
	return result, err
}

func UserExists(userName string) (bool, error) {
	// see above comment
	result, err := Command(fmt.Sprintf("id -u %s", userName))
	return result, err
}

func RegistryKey(keyName string, keyValue string) (bool, error) {

    // Break down input
	registryArgs := regexp.MustCompile("[\\\\]+").Split(keyName, -1)
    registryHiveText := registryArgs[0]
	keyPath := fmt.Sprintf(strings.Join(registryArgs[1:len(registryArgs)-1], "\\"))
	keyLoc := registryArgs[len(registryArgs)-1]

    var registryHive registry.Key
    switch registryHiveText {
    case "HKEY_CLASSES_ROOT":
        registryHive = registry.CLASSES_ROOT
    case "HKEY_CURRENT_USER":
        registryHive = registry.CURRENT_USER
    case "HKEY_LOCAL_MACHINE":
        registryHive = registry.LOCAL_MACHINE
    case "HKEY_USERS":
        registryHive = registry.USERS
    case "HKEY_CURRENT_CONFIG":
        registryHive = registry.CURRENT_CONFIG
    default:
        failPrint("Unknown registry hive: " +  registryHiveText)
        return false, errors.New("Unkown registry hive")
    }

    // Actually get the key
    k, err := registry.OpenKey(registryHive, keyPath, registry.QUERY_VALUE)
	if err != nil {
        fmt.Println("errored out 1", err.Error())
		return false, err
	}
	defer k.Close()

    // Fetch registry value
    registrySlice := make([]byte, 256)
	regLength, valType, err := k.GetValue(keyLoc, registrySlice)
    registrySlice = registrySlice[:regLength]
	fmt.Printf("Retrieved registry value was %d (length %d, type %d)\n", registrySlice, regLength, valType)

    // Determine value type to convert to string
    var registryValue string
    switch valType {
    case 1:  // SZ
        registryValue, _, err = k.GetStringValue(keyLoc)
    case 2:  // EXPAND_SZ
        registryValue, _, err = k.GetStringValue(keyLoc)
    case 3:  // BINARY
        failPrint("Binary registry format not yet supported.")
    case 4: // DWORD
        registryValue = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(registrySlice)), 10)
    default:
        failPrint("Unknown registry type: " + string(valType))
    }
	if err != nil {
        fmt.Println("Registry error:", err.Error())
        return false, err
	}

    fmt.Printf("Registry value: %s, keyvalue %s\n", registryValue, keyValue)
    if registryValue == keyValue {
    	return true, err
    }
    return false, err
}
