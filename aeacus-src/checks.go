package main

// This file contains checks that are identical for both Linux and Windows.
// If a checkType does not match one specified, it is handed off to
// processCheck for the OS-specific checks

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// processCheckWrapper takes the data from a check in the config
// and runs the correct function with the correct parameters
func processCheckWrapper(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "Command":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" passed"
		}
		result, err := Command(arg1)
		if err != nil {
			return false
		}
		return result
	case "CommandNot":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" failed"
		}
		result, err := Command(arg1)
		if err != nil {
			return false
		}
		return !result
	case "FileExists":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" exists"
		}
		result, err := FileExists(arg1)
		if err != nil {
			return false
		}
		return result
	case "FileExistsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not exist"
		}
		result, err := FileExists(arg1)
		if err != nil {
			return false
		}
		return !result
	case "FileContains":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains \"" + arg2 + "\""
		}
		result, err := FileContains(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileContainsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain \"" + arg2 + "\""
		}
		result, err := FileContains(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "FileContainsRegex":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains expression \"" + arg2 + "\""
		}
		result, err := FileContainsRegex(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileContainsRegexNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain expression \"" + arg2 + "\""
		}
		result, err := FileContainsRegex(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "DirContainsRegex":
		if check.Message == "" {
			check.Message = "Directory \"" + arg1 + "\" contains expression \"" + arg2 + "\""
		}
		result, err := DirContainsRegex(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "DirContainsRegexNot":
		if check.Message == "" {
			check.Message = "Directory \"" + arg1 + "\" does not contain expression \"" + arg2 + "\""
		}
		result, err := DirContainsRegex(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "FileEquals":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" matches hash"
		}
		result, err := FileEquals(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileEqualsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" doesn't match hash"
		}
		result, err := FileEquals(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "PackageInstalled":
		if check.Message == "" {
			check.Message = "Package \"" + arg1 + "\" is installed"
		}
		result, err := PackageInstalled(arg1)
		if err != nil {
			return false
		}
		return result
	case "PackageInstalledNot":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " has been removed"
		}
		result, err := PackageInstalled(arg1)
		if err != nil {
			return false
		}
		return !result
	case "ServiceUp":
		if check.Message == "" {
			check.Message = "Service \"" + arg1 + "\" is installed and running"
		}
		result, err := ServiceUp(arg1)
		if err != nil {
			return false
		}
		return result
	case "ServiceUpNot":
		if check.Message == "" {
			check.Message = "Service " + arg1 + " has been stopped"
		}
		result, err := ServiceUp(arg1)
		if err != nil {
			return false
		}
		return !result
	case "UserExists":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been added"
		}
		result, err := UserExists(arg1)
		if err != nil {
			return false
		}
		return result
	case "UserExistsNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been removed"
		}
		result, err := UserExists(arg1)
		if err != nil {
			return false
		}
		return !result
	case "UserInGroup":
		if check.Message == "" {
			check.Message = "User " + arg1 + " is in group \"" + arg2 + "\""
		}
		result, err := UserInGroup(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "UserInGroupNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " is not in group \"" + arg2 + "\""
		}
		result, err := UserInGroup(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "FirewallUp":
		if check.Message == "" {
			check.Message = "Firewall has been enabled"
		}
		result, err := FirewallUp()
		if err != nil {
			return false
		}
		return result
	case "FirewallUpNot":
		if check.Message == "" {
			// Who is ever going to use this?
			// Maybe as a penalty?
			check.Message = "Firewall has been disabled"
		}
		result, err := FirewallUp()
		if err != nil {
			return false
		}
		return !result
	default:
		return processCheck(check, checkType, arg1, arg2, arg3)
	}
}

func FileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err), nil
}

func FileContains(fileName string, searchString string) (bool, error) {
	fileContent, err := readFile(fileName)
	return strings.Contains(strings.TrimSpace(fileContent), searchString), err
}

func FileContainsRegex(fileName string, expressionString string) (bool, error) {
	fileContent, err := readFile(fileName)
	if err != nil {
		return false, err
	}
	matched, err := regexp.Match(expressionString, []byte(fileContent))
	if err != nil {
		failPrint("There's an error with your regular expression for FileContainsRegex: " + err.Error())
	}
	return matched, err
}

// DirContainsRegex returns true if any file in the directory matches the regular expression provided
func DirContainsRegex(dirName string, expressionString string) (bool, error) {
	result, err := FileExists(dirName)
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
		result, err := FileContainsRegex(file, expressionString)
		if err != nil {
			return false, err
		}
		if result {
			return result, nil
		}
	}
	return false, nil
}

// FileEquals calculates the SHA1 sum of a file and compares it
// with the hash provided in the check
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
