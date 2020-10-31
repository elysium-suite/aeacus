// checks.go contains checks that are identical for both Linux and Windows.
// If a checkType does not match one specified, it is handed off to
// processCheck for the OS-specific checks.

package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// processCheckWrapper takes the data from a check in the config
// and runs the correct function with the correct parameters
func processCheckWrapper(check *check, checkType, arg1, arg2, arg3 string) bool {
	if err := deobfuscateData(&checkType); err != nil {
		// wat
	}
	if err := deobfuscateData(&arg1); err != nil {
		// wat
	}
	if err := deobfuscateData(&arg2); err != nil {
		// wat
	}
	if err := deobfuscateData(&arg3); err != nil {
		// wat
	}
	switch checkType {
	case "Command":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" passed"
		}
		result, err := command(arg1)
		return err == nil && result
	case "CommandNot":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" failed"
		}
		result, err := command(arg1)
		return err == nil && !result
	case "CommandOutput":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" had the output \"" + arg2 + "\""
		}
		result, err := commandOutput(arg1, arg2)
		return err == nil && result
	case "CommandOutputNot":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" did not have the output \"" + arg2 + "\""
		}
		result, err := commandOutput(arg1, arg2)
		return err == nil && !result
	case "CommandContains":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" contained output \"" + arg2 + "\""
		}
		result, err := commandContains(arg1, arg2)
		return err == nil && result
	case "CommandContainsNot":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" output did not contain \"" + arg2 + "\""
		}
		result, err := commandContains(arg1, arg2)
		return err == nil && !result
	case "FileExists":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" exists"
		}
		result, err := pathExists(arg1)
		return err == nil && result
	case "FileExistsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not exist"
		}
		result, err := pathExists(arg1)
		return err == nil && !result
	case "FileContains":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains \"" + arg2 + "\""
		}
		result, err := fileContains(arg1, arg2)
		return err == nil && result
	case "FileContainsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain \"" + arg2 + "\""
		}
		result, err := fileContains(arg1, arg2)
		return err == nil && !result
	case "FileContainsRegex":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains expression \"" + arg2 + "\""
		}
		result, err := fileContainsRegex(arg1, arg2)
		return err == nil && result
	case "FileContainsRegexNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain expression \"" + arg2 + "\""
		}
		result, err := fileContainsRegex(arg1, arg2)
		return err == nil && !result
	case "DirContainsRegex":
		if check.Message == "" {
			check.Message = "Directory \"" + arg1 + "\" contains expression \"" + arg2 + "\""
		}
		result, err := dirContainsRegex(arg1, arg2)
		return err == nil && result
	case "DirContainsRegexNot":
		if check.Message == "" {
			check.Message = "Directory \"" + arg1 + "\" does not contain expression \"" + arg2 + "\""
		}
		result, err := dirContainsRegex(arg1, arg2)
		return err == nil && !result
	case "FileEquals":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" matches hash"
		}
		result, err := fileEquals(arg1, arg2)
		return err == nil && result
	case "FileEqualsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" doesn't match hash"
		}
		result, err := fileEquals(arg1, arg2)
		return err == nil && !result
	case "PackageInstalled":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " is installed"
		}
		result, err := packageInstalled(arg1)
		return err == nil && result
	case "PackageInstalledNot":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " has been removed"
		}
		result, err := packageInstalled(arg1)
		return err == nil && !result
	case "ServiceUp":
		if check.Message == "" {
			check.Message = "Service \"" + arg1 + "\" is installed and running"
		}
		result, err := serviceUp(arg1)
		return err == nil && result
	case "ServiceUpNot":
		if check.Message == "" {
			check.Message = "Service " + arg1 + " has been stopped"
		}
		result, err := serviceUp(arg1)
		return err == nil && !result
	case "UserExists":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been added"
		}
		result, err := userExists(arg1)
		return err == nil && result
	case "UserExistsNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been removed"
		}
		result, err := userExists(arg1)
		return err == nil && !result
	case "UserInGroup":
		if check.Message == "" {
			check.Message = "User " + arg1 + " is in group \"" + arg2 + "\""
		}
		result, err := userInGroup(arg1, arg2)
		return err == nil && result
	case "UserInGroupNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " removed or is not in group \"" + arg2 + "\""
		}
		result, err := userInGroup(arg1, arg2)
		return err == nil && !result
	case "FirewallUp":
		if check.Message == "" {
			check.Message = "Firewall has been enabled"
		}
		result, err := firewallUp()
		return err == nil && result
	case "FirewallUpNot":
		if check.Message == "" {
			check.Message = "Firewall has been disabled"
		}
		result, err := firewallUp()
		return err == nil && !result
	default:
		return processCheck(check, checkType, arg1, arg2, arg3)
	}
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

func commandContains(commandGiven, desiredContains string) (bool, error) {
	out, err := rawCmd(commandGiven).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	outString := strings.TrimSpace(string(out))
	if strings.Contains(outString, desiredContains) {
		return true, nil
	}
	return false, nil
}

// pathExists is a wrapper around os.Stat and os.IsNotExist, and determines
// whether a file or folder exists.
func pathExists(pathName string) (bool, error) {
	_, err := os.Stat(pathName)
	return !os.IsNotExist(err), nil // TODO is not not IsNotExist instead of nil
}

// fileContains searches for a given searchString in the provided fileName.
func fileContains(fileName, searchString string) (bool, error) {
	fileContent, err := readFile(fileName)
	return strings.Contains(strings.TrimSpace(fileContent), searchString), err
}

// fillContainsRegex determines whether a file contains a given regular
// expression.
//
// Newlines in regex may not work as expected, especially on Windows. It's
// best to not use these (ex. ^ and $).
func fileContainsRegex(fileName, expressionString string) (bool, error) {
	fileContent, err := readFile(fileName)
	if err != nil {
		return false, err
	}
	matched := false
	found := false
	for _, line := range strings.Split(fileContent, "\n") {
		found, err = regexp.Match(expressionString, []byte(line))
		if found {
			matched = found
		}
	}
	if err != nil {
		failPrint("There's an error with your regular expression for fileContainsRegex: " + err.Error())
	}
	return matched, err
}

// dirContainsRegex returns true if any file in the directory matches the regular expression provided
func dirContainsRegex(dirName, expressionString string) (bool, error) {
	result, err := pathExists(dirName)
	if err != nil || !result {
		return false, errors.New("DirContainsRegex: file does not exist")
	}

	var files []string
	err = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		if len(files) > 10000 {
			failPrint("Recursive indexing has exceeded limit, erroring out.")
			return errors.New("Indexed too many files in recursive search")
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	for _, file := range files {
		result, err := fileContainsRegex(file, expressionString)
		if os.IsPermission(err) {
			return false, err
		}

		if result {
			return result, nil
		}
	}
	return false, nil
}

// fileEquals calculates the SHA1 sum of a file and compares it
// with the hash provided in the check.
func fileEquals(fileName, fileHash string) (bool, error) {
	fileContent, err := readFile(fileName)
	if err != nil {
		return false, err
	}
	hasher := sha1.New()
	_, err = hasher.Write([]byte(fileContent))
	if err != nil {
		return false, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash == fileHash, nil
}
