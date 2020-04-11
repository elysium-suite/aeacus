package main

import (
	"os"
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
			check.Message = "Registry key " + arg1 + arg2 + " matches " + arg3
		}
		result, err := UserExists(arg1)
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

func FileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err), nil
}

func FileContains(fileName string, searchString string) (bool, error) {
	fileContent, err := readFile(fileName)
	return strings.Contains(fileContent, searchString), err
}

func FileContainsRegex(fileName string, expressionString string) (bool, error) {
	fileContent, _ := readFile(fileName)
	matched, err := regexp.Match(expressionString, []byte(fileContent))
	return matched, err
}

func FileEquals(fileName string, fileHash string) (bool, error) {
	fileContent, err := readFile(fileName)
	if err != nil {
		return false, err
	}
	hasher := sha1.New()
	hasher.Write([]byte(fileContent))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash == fileHash, err
}

func PackageInstalledL(packageName string) (bool, error) {
	// not super happy with the command implementation
	// could just keylog sh or replace dpkg binary or something
	// should use golang dpkg library if it existed and was good
	result, err := CommandL(fmt.Sprintf("dpkg -l %s", packageName))
	return result, err
}

func PackageInstalledW(packageName string) (bool, error) {
	// not super happy with the command implementation
	// could just keylog sh or replace dpkg binary or something
	// should use golang dpkg library if it existed and was good
	result, err := CommandL(fmt.Sprintf("dpkg -l %s", packageName))
	return result, err
}

func UserExistsL(userName string) (bool, error) {
	// see above comment
	result, err := CommandL(fmt.Sprintf("id -u %s", userName))
	return result, err
}

func UserExistsW(userName string) (bool, error) {
	// see above comment
	result, err := CommandL(fmt.Sprintf("id -u %s", userName))
	return result, err
}

func RegistryKey(keyName string, keyValue string) (bool, error) {
	registryArgs := regexp.MustCompile("[\\s]+").Split(keyName, -1)
	keyPath := registryArgs[:len(registryArgs)-1]
	keyLoc := registryArgs[len(registryArgs)]
	fmt.Printf("PATH %s getting KEY %s", keyPath, keyLoc)
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()

	s, _, err := k.GetStringValue(keyLoc)
	if err != nil {
		return false, err
	}
	fmt.Printf("retreievd reg value was %s", s)
	return true, err
}
