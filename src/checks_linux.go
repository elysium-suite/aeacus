package main

import (
	"os/exec"
	"strings"
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
			check.Message = "Password for " + arg1 + " has been changed"
		}
		result, err := passwordChanged(arg1, arg2)
		return err == nil && result
	case "PasswordChangedNot":
		if check.Message == "" {
			check.Message = "Password for " + arg1 + " has not been changed"
		}
		result, err := passwordChanged(arg1, arg2)
		return err == nil && !result
	default:
		failPrint("No check type " + checkType)
	case "PackageVersion":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " is version " + arg2
		}
		result, err := packageVersion(arg1, arg2)
		return err == nil && result
	case "PackageVersionNot":
		if check.Message == "" {
			check.Message = "Package " + arg1 + " is version " + arg2
		}
		result, err := packageVersion(arg1, arg2)
		return err == nil && !result
	case "KernelVersion":
		if check.Message == "" {
			check.Message = "Kernel is version " + arg1
		}
		result, err := kernelVersion(arg1)
		return err == nil && result
	case "KernelVersionNot":
		if check.Message == "" {
			check.Message = "Kernel is not version " + arg1
		}
		result, err := kernelVersion(arg1)
		return err == nil && !result
	case "AutoCheckUpdatesEnabled":
		if check.Message == "" {
			check.Message = "The system automatically checks for updates daily"
		}
		result, err := autoCheckUpdatesEnabled()
		return err == nil && result
	case "AutoCheckUpdatesEnabledNot":
		if check.Message == "" {
			check.Message = "The system does not automatically checks for updates daily"
		}
		result, err := autoCheckUpdatesEnabled()
		return err == nil && !result
	case "OctalPermissionIs":
		if check.Message == "" {
			check.Message = "The octal permissions value of the " + arg1 + " file is " + arg2
		}
		result, err := octalPermissionIs(arg1, arg2)
		return err == nil && result
	case "OctalPermissionIsNot":
		if check.Message == "" {
			check.Message = "The octal permissions value of the " + arg1 + " file is not " + arg2
		}
		result, err := octalPermissionIs(arg1, arg2)
		return err == nil && !result
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

func commandOutput(commandGiven, desiredOutput string) (bool, error) {
	out, err := exec.Command("sh", "-c", commandGiven).Output()
	if err != nil {
		return false, err
	}
	outString := strings.TrimSpace(string(out))
	if outString == desiredOutput {
		return true, nil
	}
	return false, nil
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
	return command("groups " + userName + " | grep -q " + groupName + "")
}

func firewallUp() (bool, error) {
	return command("ufw status | grep -q 'Status: active'")
}

func passwordChanged(user, hash string) (bool, error) {
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

func packageVersion(packageName string, versionNumber string) (bool, error) {
	return command(`dpkg -l | awk '$2=="` + packageName + `" { print $3 }' | grep -q "` + versionNumber + `"`)
}

func kernelVersion(version string) (bool, error) {
	return command("uname -r | grep -q " + version)
}

func autoCheckUpdatesEnabled() (bool, error) {
	return fileContainsRegex("/etc/apt/apt.conf.d/20auto-upgrades", `APT::Periodic::Update-Package-Lists( |)"1";`)
}

func octalPermissionIs(filePath string, permissionToCheck string) (bool, error) {
	return command(`stat -c '%a' ` + filePath + ` | grep -q ` + permissionToCheck)
}