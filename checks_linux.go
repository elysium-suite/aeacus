package main

import (
	"errors"
	"os"
	"strings"
	"syscall"
)

func (c cond) AutoCheckUpdatesEnabled() (bool, error) {
	return cond{
		Path:  "/etc/apt/apt.conf.d/",
		Value: `(?i)^\s*APT::Periodic::Update-Package-Lists\s+"1"\s*;\s*$`,
	}.DirContains()
}

// Command checks if a given shell command ran successfully (that is, did not
// return or raise any errors).
func (c cond) Command() (bool, error) {
	c.requireArgs("Cmd")
	if c.Cmd == "" {
		fail("Missing command for", c.Type)
	}
	err := shellCommand(c.Cmd)
	if err != nil {
		// This check does not return errors, since it is based on successful
		// execution. If any errors occurred, it means that the check failed,
		// not errored out.
		//
		// It would be an error if failure to execute the command resulted in
		// an inability to meaningfully score the check (e.g., if the uname
		// syscall failed for KernelVersion).
		return false, nil
	}
	return true, nil
}

func (c cond) FirewallUp() (bool, error) {
	return cond{
		Path:  "/etc/ufw/ufw.conf",
		Value: `^\s*ENABLED=yes\s*$`,
	}.FileContains()
}

func (c cond) GuestDisabledLDM() (bool, error) {
	guestStr := `\s*allow-guest\s*=\s*false`
	result, err := cond{
		Path:  "/usr/share/lightdm/lightdm.conf.d/",
		Value: guestStr,
	}.DirContains()
	if !result {
		return cond{
			Path:  "/etc/lightdm/",
			Value: guestStr,
		}.DirContains()
	}
	return result, err
}

func (c cond) KernelVersion() (bool, error) {
	c.requireArgs("Value")
	utsname := syscall.Utsname{}
	err := syscall.Uname(&utsname)
	releaseUint := []byte{}
	for i := 0; i < 65; i++ {
		if utsname.Release[i] == 0 {
			break
		}
		releaseUint = append(releaseUint, uint8(utsname.Release[i]))
	}
	debug("System uname value is", string(releaseUint), "and our value is", c.Value)
	return string(releaseUint) == c.Value, err
}

func (c cond) PasswordChanged() (bool, error) {
	c.requireArgs("User", "Value")
	fileContent, err := readFile("/etc/shadow")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(fileContent, "\n") {
		if strings.Contains(line, c.User+":") {
			if strings.Contains(line, c.User+":"+c.Value) {
				debug("Exact value found in /etc/shadow for user", c.User+":", line)
				return false, nil
			}
			debug("Differing value found in /etc/shadow for user", c.User+":", line)
			return true, nil
		}
	}
	return false, errors.New("user not found")
}

func (c cond) PermissionIs() (bool, error) {
	c.requireArgs("Path", "Value")
	f, err := os.Stat(c.Path)
	if err != nil {
		return false, err
	}

	fileMode := f.Mode()
	modeBytes := []byte(fileMode.String())
	if len(modeBytes) != 10 {
		fail("System permission string is wrong length:", string(modeBytes))
		return false, errors.New("Invalid system permission string")
	}

	// Permission string includes suid/sgid as the special bit (MSB), while
	// GNU coreutils replaces the executable bit, which we need to emulate.
	if fileMode&os.ModeSetuid != 0 {
		modeBytes[0] = '-'
		modeBytes[3] = 's'
	}
	if fileMode&os.ModeSetgid != 0 {
		modeBytes[0] = '-'
		modeBytes[6] = 's'
	}

	c.Value = strings.TrimSpace(c.Value)
	if len(c.Value) == 9 {
		c.Value = "-" + c.Value
	} else if len(c.Value) != 10 {
		fail("Your permission string is the wrong length (should be 9 or 10 characters):", c.Value)
		return false, errors.New("Invalid user permission string")
	}

	for i := 0; i < len(c.Value); i++ {
		if c.Value[i] == '?' {
			continue
		}
		if c.Value[i] != modeBytes[i] {
			return false, nil
		}
	}
	return true, nil
}

func (c cond) ProgramInstalled() (bool, error) {
	c.requireArgs("Name")
	return cond{
		Cmd: "dpkg -s " + c.Name,
	}.Command()
}

func (c cond) ProgramVersion() (bool, error) {
	c.requireArgs("Name", "Value")
	return cond{
		Cmd:   `dpkg -s ` + c.Name + ` | grep Version | cut -d" " -f2`,
		Value: c.Value,
	}.CommandOutput()
}

func (c cond) ServiceUp() (bool, error) {
	// TODO: detect and use other init systems
	c.requireArgs("Name")
	return cond{
		Cmd: "systemctl is-active " + c.Name,
	}.Command()
}

func (c cond) UserExists() (bool, error) {
	c.requireArgs("User")
	return cond{
		Path:  "/etc/passwd",
		Value: "^" + c.User + ":",
	}.FileContains()
}

func (c cond) UserInGroup() (bool, error) {
	c.requireArgs("User", "Group")
	return cond{
		Path:  "/etc/group",
		Value: c.Group + `[0-9a-zA-Z,:\s+]+` + c.User,
	}.FileContains()
}
