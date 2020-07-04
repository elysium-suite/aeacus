package main

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
)

// This Linux processCheck will process Linux-specific checks
// handed to it by the processCheckWrapper function
func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "GuestDisabledLDM":
		if check.Message == "" {
			check.Message = "Guest is disabled"
		}
		result, err := GuestDisabledLDM()
		if err != nil {
			return false
		}
		return result
	case "GuestDisabledLDMNot":
		if check.Message == "" {
			check.Message = "Guest is enabled"
		}
		result, err := GuestDisabledLDM()
		if err != nil {
			return false
		}
		return !result
	default:
		failPrint("No check type " + checkType)
	}
	return false
}

func adminCheck() bool {
	currentUser, err := user.Current()
	uid, _ := strconv.Atoi(currentUser.Uid)
	if err != nil {
		failPrint("Error for checking if running as root.")
		fmt.Println(err)
		return false
	} else if uid != 0 {
		return false
	}
	return true
}

func Command(commandGiven string) (bool, error) {
	cmd := exec.Command("sh", "-c", commandGiven)
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
	// note: dpkg -l does not always return proper error code
	return Command(fmt.Sprintf(`dpkg -s "%s"`, packageName))
}

func ServiceUp(serviceName string) (bool, error) {
	return Command(fmt.Sprintf(`systemctl is-active "%s"`, serviceName))
}

func UserExists(userName string) (bool, error) {
	return Command(fmt.Sprintf(`id -u "%s"`, userName))
}

func UserInGroup(userName string, groupName string) (bool, error) {
	return Command(fmt.Sprintf(`groups "%s" | grep -q "%s"`, userName, groupName))
}

func FirewallUp() (bool, error) {
	return Command("ufw status | grep -q 'Status: active'")
}

func GuestDisabledLDM() (bool, error) {
	result, err := DirContainsRegex("/usr/share/lightdm/lightdm.conf.d/", "allow-guest( |)=( |)false")
	if !result && err == nil {
		result, err = DirContainsRegex("/etc/lightdm/", "allow-guest( |)=( |)false")
	}
	return result, err
}
