package main

import (
	"fmt"
	"os/exec"
)

// processCheck (Linux) will process Linux-specific checks
// handed to it by the processCheckWrapper function
func processCheck(check *check, checkType string, arg1 string, arg2 string, arg3 string) bool {
	switch checkType {
	case "GuestDisabledLDM":
		if check.Message == "" {
			check.Message = "Guest is disabled"
		}
		result, err := guestDisabledLDM()
		return err == nil && result
	case "GuestDisabledLDMNot":
		if check.Message == "" {
			check.Message = "Guest is enabled"
		}
		result, err := guestDisabledLDM()
		return err == nil && !result
	case "PasswordChanged":
		if check.Message == "" {
			check.Message = "Insecure password has been changed"
		}
	default:
		failPrint("No check type " + checkType)
	}
	return false
}

func command(commandGiven string) (bool, error) {
	cmd := exec.Command("sh", "-c", commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
}

func packageInstalled(packageName string) (bool, error) {
	return command("dpkg -s " + packageName)
}

func serviceUp(serviceName string) (bool, error) {
	return command("systemctl is-active " + serviceName)
}

func userExists(userName string) (bool, error) {
	return command("id -u " + userName)
}

func userInGroup(userName string, groupName string) (bool, error) {
	return command(fmt.Sprintf(`groups "%s" | grep -q "%s"`, userName, groupName))
}

func firewallUp() (bool, error) {
	return command("ufw status | grep -q 'Status: active'")
}

func passwordChanged(hash string) (bool, error) {
	res, err := fileContains("/etc/shadow", hash)
	return !res, err
}

func guestDisabledLDM() (bool, error) {
	result, err := dirContainsRegex("/usr/share/lightdm/lightdm.conf.d/", "allow-guest( |)=( |)false")
	if !result && err == nil {
		result, err = dirContainsRegex("/etc/lightdm/", "allow-guest( |)=( |)false")
	}
	return result, err
}
