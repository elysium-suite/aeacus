package main

import (
	"os"

	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
)

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
			check.Message = "File \"" + arg1 + "\" is installed"
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
