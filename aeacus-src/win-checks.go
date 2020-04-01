package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func adminCheckW() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func scoreW(mc *metaConfig, id *imageData) {
	id.Score = 0
	id.ScoredVulns = 0
	id.TotalPoints = 0
	id.Points = []scoreItem{}
	pointlessChecks := []int{}

	for i, check := range mc.Config.Check {
		if check.Points == 0 {
			pointlessChecks = append(pointlessChecks, i)
			id.ScoredVulns += 1
		} else if check.Points > 0 {
			id.TotalPoints += check.Points
			id.ScoredVulns += 1
		}
	}

	pointsLeft := 100 - id.TotalPoints
	if pointsLeft > 0 && len(pointlessChecks) > 0 {
		pointsEach := pointsLeft / len(pointlessChecks)
		for _, check := range pointlessChecks {
			mc.Config.Check[check].Points = pointsEach
		}
		id.TotalPoints += (pointsEach * len(pointlessChecks))
		if id.TotalPoints != 100 {
			mc.Config.Check[pointlessChecks[0]].Points += (100 - id.TotalPoints)
			id.TotalPoints += (100 - id.TotalPoints)
		}
	}

	for _, check := range mc.Config.Check {
		status := false
		failStatus := false
		for _, condition := range check.Pass {
			status = processCheckW(mc, &check, condition.Type, condition.Arg1, condition.Arg2)
			if status {
				break
			}
		}
		for _, condition := range check.Fail {
			failStatus = processCheckW(mc, &check, condition.Type, condition.Arg1, condition.Arg2)
			if failStatus {
				status = false
				break
			}
		}
		if check.Points >= 0 {
			if status {
				if mc.Cli.Bool("v") {
					passPrint("")
					fmt.Printf("Check passed: %s - %d pts\n", check.Message, check.Points)
				}
				id.Points = append(id.Points, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Contribs += check.Points
			}
		} else {
			if status {
				if mc.Cli.Bool("v") {
					failPrint("")
					fmt.Printf("Penalty triggered: %s - %d pts\n", check.Message, check.Points)
				}
				id.Penalties = append(id.Penalties, scoreItem{check.Message, check.Points})
				id.Score += check.Points
				id.Detracts += check.Points
			}
		}
	}
	if mc.Cli.Bool("v") {
		infoPrint("")
		fmt.Printf("Score: %d\n", id.Score)
	}
}

func processCheckW(mc *metaConfig, check *check, checkType string, arg1 string, arg2 string) bool {
	switch checkType {
	case "Command":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" passed"
		}
		result, err := CommandW(arg1)
		if err != nil {
			return false
		}
		return result
	case "CommandNot":
		if check.Message == "" {
			check.Message = "Command \"" + arg1 + "\" failed"
		}
		result, err := CommandW(arg1)
		if err != nil {
			return false
		}
		return !result
	case "FileExists":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" exists"
		}
		result, err := FileExistsW(arg1)
		if err != nil {
			return false
		}
		return result
	case "FileExistsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not exist"
		}
		result, err := FileExistsW(arg1)
		if err != nil {
			return false
		}
		return !result
	case "FileContains":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains \"" + arg2 + "\""
		}
		result, err := FileContainsW(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileContainsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain + \"" + arg2 + "\""
		}
		result, err := FileContainsW(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "FileContainsRegex":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" contains expression \"" + arg2 + "\""
		}
		result, err := FileContainsRegexW(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileContainsRegexNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" does not contain expression \"" + arg2 + "\""
		}
		result, err := FileContainsRegexW(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "FileEquals":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" matches hash"
		}
		result, err := FileEqualsW(arg1, arg2)
		if err != nil {
			return false
		}
		return result
	case "FileEqualsNot":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" doesn't match hash"
		}
		result, err := FileEqualsW(arg1, arg2)
		if err != nil {
			return false
		}
		return !result
	case "PackageInstalled":
		if check.Message == "" {
			check.Message = "File \"" + arg1 + "\" is installed"
		}
		result, err := PackageInstalledW(arg1)
		if err != nil {
			return false
		}
		return result
	case "PackageInstalledNot":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " has been removed"
		}
		result, err := PackageInstalledW(arg1)
		if err != nil {
			return false
		}
		return !result
	case "UserExists":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been added"
		}
		result, err := UserExistsW(arg1)
		if err != nil {
			return false
		}
		return result
	case "UserExistsNot":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been removed"
		}
		result, err := UserExistsW(arg1)
		if err != nil {
			return false
		}
		return !result
	default:
		if mc.Cli.Bool("v") {
			failPrint("No check type " + checkType)
		}
	}
	return false
}

/////////////////////
// CHECK FUNCTIONS //
/////////////////////

func CommandW(commandGiven string) (bool, error) {
	cmd := exec.Command("sh", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
}

func FileExistsW(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err), nil
}

func FileContainsW(fileName string, searchString string) (bool, error) {
	fileContent, err := readFile(fileName)
	return strings.Contains(fileContent, searchString), err
}

func FileContainsRegexW(fileName string, expressionString string) (bool, error) {
	fileContent, _ := readFile(fileName)
	matched, err := regexp.Match(expressionString, []byte(fileContent))
	return matched, err
}

func FileEqualsW(fileName string, fileHash string) (bool, error) {
	fileContent, err := readFile(fileName)
	if err != nil {
		return false, err
	}
	hasher := sha1.New()
	hasher.Write([]byte(fileContent))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash == fileHash, err
}

func PackageInstalledW(packageName string) (bool, error) {
	// not super happy with the command implementation
	// could just keylog sh or replace dpkg binary or something
	// should use golang dpkg library if it existed and was good
	result, err := CommandW(fmt.Sprintf("dpkg -l %s", packageName))
	return result, err
}

func UserExistsW(userName string) (bool, error) {
	// see above comment
	result, err := CommandW(fmt.Sprintf("id -u %s", userName))
	return result, err
}
