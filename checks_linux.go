package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func (c cond) AutoCheckUpdatesEnabled() (bool, error) {
	return cond{
		Path:  "/etc/apt/apt.conf.d/",
		Value: `APT::Periodic::Update-Package-Lists\s+"1"`,
	}.DirContains()
}

// Command checks if a given shell command ran successfully (that is, did not
// return or raise any errors).
func (c cond) Command() (bool, error) {
	err := shellCommand(c.Cmd)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c cond) FirewallUp() (bool, error) {
	return cond{
		Path:  "/etc/ufw/ufw.conf",
		Value: `\s*ENABLED\s*=\s*yes`,
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
	utsname := syscall.Utsname{}
	err := syscall.Uname(&utsname)
	releaseUint := []byte{}
	for i := 0; i < 65; i++ {
		if utsname.Release[i] == 0 {
			break
		}
		releaseUint = append(releaseUint, uint8(utsname.Release[i]))
	}
	return string(releaseUint) == c.Version, err
}

func (c cond) PasswordChanged() (bool, error) {
	if c.Hash == "" {
		fail("Missing hash for PasswordChanged.")
		return false, errors.New("missing hash")
	}

	fileContent, err := readFile("/etc/shadow")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(fileContent, "\n") {
		if strings.Contains(line, c.User+":") {
			if strings.Contains(line, c.User+":"+c.Hash) {
				return false, nil
			}
			return true, nil
		}
	}
	return false, errors.New("user not found")
}

func (c cond) PermissionIs() (bool, error) {
	f, err := os.Stat(c.Path)
	if err != nil {
		return false, err
	}

	perm := fmt.Sprint(f.Mode().Perm())
	if len(perm) != 10 {
		fail("System permission string is wrong length:", perm)
		return false, errors.New("Invalid system permission string")
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
		if c.Value[i] != perm[i] {
			return false, nil
		}
	}
	return true, nil
}

func (c cond) ProgramInstalled() (bool, error) {
	return cond{
		Cmd: "dpkg -s " + c.Program,
	}.Command()
}

func (c cond) ProgramVersion() (bool, error) {
	commandGiven := `dpkg -s ` + c.Program + ` | grep Version | cut -d" " -f2`
	out, err := rawCmd(commandGiven).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == c.Version, nil
}

func (c cond) ServiceUp() (bool, error) {
	// TODO: detect and use other init systems
	ret, err := cond{
		Cmd: "systemctl is-active " + c.Name,
	}.Command()
	return ret, err
}

func (c cond) UserExists() (bool, error) {
	return cond{
		Path:  "/etc/passwd",
		Value: "^" + c.User + ":",
	}.FileContains()
}

func (c cond) UserInGroup() (bool, error) {
	return cond{
		Path:  "/etc/group",
		Value: c.Group + `[0-9a-zA-Z,:\s+]+` + c.User,
	}.FileContains()
}
