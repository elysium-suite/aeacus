package main

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
)

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

func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "MagicLinuxOnlyCheck":
		if check.Message == "" {
			check.Message = "User " + arg1 + " has been removed"
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
	return Command(fmt.Sprintf("dpkg -l %s", packageName))
}

func ServiceUp(serviceName string) (bool, error) {
	return Command("systemctl is-active " + serviceName)
}

func UserExists(userName string) (bool, error) {
	// see above comment
	return Command("id -u " + userName)
}

func UserInGroup(userName string, groupName string) (bool, error) {
	return Command("groups " + userName + " | grep -q " + groupName)
}

func FirewallUp() (bool, error) {
	return FileExists("/proc/net/ip_tables_names")
}

// For when DirContainsRegex() is ready...
// func GuestDisabled() (bool, error) {
// 	return DirContainsRegex("/usr/share/lightdm/lightdm.conf.d/", "allow-guest*=*false")
// }
// Also, more restrictive regex would be nice...
