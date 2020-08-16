package main

import (
	"os/exec"
	"strings"
)

// processCheck (Linux) will process Linux-specific checks
// handed to it by the processCheckWrapper function
func processCheck(check *check, checkType, arg1, arg2, arg3 string) bool {
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
	case "PermissionIs":
		if check.Message == "" {
			check.Message = "The permissions of " + arg1 + " are " + arg2
		}
		result, err := permissionIs(arg1, arg2)
		return err == nil && result
	case "PermissionIsNot":
		if check.Message == "" {
			check.Message = "The permissions of " + arg1 + " are not " + arg2
		}
		result, err := permissionIs(arg1, arg2)
		return err == nil && !result
	case "FirefoxPrefIs":
		if check.Message == "" {
			check.Message = "Firefox preference " + arg1 + " is set to " + arg2
		}
		result, err := firefoxSetting(arg1, arg2)
		return err == nil && result
	case "FirefoxPrefIsNot":
		if check.Message == "" {
			check.Message = "Firefox preference " + arg1 + " is not set to " + arg2
		}
		result, err := firefoxSetting(arg1, arg2)
		return err == nil && !result
	}
	return false
}

func command(commandGiven string) (bool, error) {
	cmd := rawCmd(commandGiven)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
	}
	return true, nil
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

func packageInstalled(packageName string) (bool, error) {
	return command("dpkg -s " + packageName)
}

func serviceUp(serviceName string) (bool, error) {
	return command("systemctl is-active " + serviceName)
}

func userExists(userName string) (bool, error) {
	return command("id -u " + userName)
}

func userInGroup(userName, groupName string) (bool, error) {
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

func packageVersion(packageName, versionNumber string) (bool, error) {
	return command(`dpkg -l | awk '$2=="` + packageName + `" { print $3 }' | grep -q "` + versionNumber + `"`)
}

func kernelVersion(version string) (bool, error) {
	return command("uname -r | grep -q " + version)
}

func autoCheckUpdatesEnabled() (bool, error) {
	return fileContainsRegex("/etc/apt/apt.conf.d/20auto-upgrades", `APT::Periodic::Update-Package-Lists( |)"1";`)
}

func permissionIs(filePath, permissionToCheck string) (bool, error) {
	return command(`stat -c '%a' ` + filePath + ` | grep -q ` + permissionToCheck)
}

func firefoxSetting(param, value string) (bool, error) {
	dirs := []string{
		"/usr/lib/firefox/defaults/pref/",
		"/usr/lib/firefox/",
		"/etc/firefox/",
		"/home/" + mc.Config.User + "/.mozilla",
	}
	prefStyle := []string{
		"lockPref",
		"pref",
		"user_pref",
	}

	res := false
	var finalErr error
	for _, el1 := range dirs {
		for _, el2 := range prefStyle {
			res2, err := dirContainsRegex(el1, el2+`("`+param+`", `+value+`);`)
			res = res || res2
			if err != nil {
				finalErr = err // err handling janky, returns last err, plsfix ~safin
			}
		}
	}

	return res, finalErr
}
