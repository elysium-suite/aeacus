package main

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	version = "2.0.0"
)

var (
	yesEnabled     bool
	verboseEnabled bool
	debugEnabled   bool
	dirPath        string
	scoringConf    = "scoring.conf"
	scoringData    = "scoring.dat"
)

var (
	timeStart             = time.Now()
	timeWithoutID, _      = time.ParseDuration("0s")
	withoutIDThreshold, _ = time.ParseDuration("30m")
)

const (
	shellCmdLen = 15
)

// determineDirectory sets the dirPath variable based on its environment.
func determineDirectory() error {
	if dirPath == "" {
		if runtime.GOOS == "linux" {
			dirPath = "/opt/aeacus/"
		} else if runtime.GOOS == "windows" {
			dirPath = `C:\aeacus\`
		} else {
			fail("Unknown OS (" + runtime.GOOS + "): you need to specify an aeacus directory")
			return errors.New("unknown OS: " + runtime.GOOS)
		}
	} else if dirPath[len(dirPath)-1] != '\\' && dirPath[len(dirPath)-1] != '/' {
		return errors.New("Your scoring directory must end in a slash: " + dirPath + "/")
	}
	return nil
}

// timeCheck calls destroyImage if the configured EndDate for the image has
// passed. Its purpose is to dissuade or prevent people using an image after
// the round ends.
func timeCheck() {
	if conf.EndDate != "" {
		date, err := time.Parse("2006/01/02 15:04:05 MST", conf.EndDate)
		if err != nil {
			fail("Your EndDate value in the configuration is invalid: " + err.Error())
		} else {
			if time.Now().After(date) {
				destroyImage()
			}
		}
	}
}

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func writeFile(fileName, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0o644)
	if err != nil {
		fail("Error writing file: " + err.Error())
	}
}

// PermsCheck is a convenience function wrapper around
// adminCheck, which prints an error indicating that admin
// permissions are needed.
func permsCheck() {
	if !adminCheck() {
		fail("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
}

// shellCommand executes a given command in a shell environment.
func shellCommand(commandGiven string) error {
	cmd := rawCmd(commandGiven)
	if err := cmd.Run(); err != nil {
		if verboseEnabled {
			if len(commandGiven) > shellCmdLen {
				fail("Command \"" + commandGiven[:shellCmdLen] + "...\" errored out (code " + err.Error() + ").")
			} else {
				fail("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
			}
		}
		return err
	}
	return nil
}

// shellCommandOutput executes a given command in a shell environment and
// returns its output.
func shellCommandOutput(commandGiven string) (string, error) {
	out, err := rawCmd(commandGiven).Output()
	if err != nil {
		if verboseEnabled {
			if len(commandGiven) > shellCmdLen {
				fail("Command \"" + commandGiven[:shellCmdLen] + "...\" errored out (code " + err.Error() + ").")
			} else {
				fail("Command \"" + commandGiven + "\" errored out (code " + err.Error() + ").")
			}
		}
		return "", err
	}
	return string(out), err
}

// assignPoints is used to automatically assign points to checks that don't
// have a hardcoded points value.
func assignPoints() {
	pointlessChecks := []int{}

	for i, check := range conf.Check {
		if check.Points == 0 {
			pointlessChecks = append(pointlessChecks, i)
		} else if check.Points > 0 {
			image.TotalPoints += check.Points
		}
	}

	pointsLeft := 100 - image.TotalPoints
	if pointsLeft <= 0 && len(pointlessChecks) > 0 || len(pointlessChecks) > 100 {
		// If the specified points already value over 100, yet there are checks
		// without points assigned, we assign the default point value of 3
		// (arbitrarily chosen).
		for _, check := range pointlessChecks {
			conf.Check[check].Points = 3
		}
	} else if pointsLeft > 0 && len(pointlessChecks) > 0 {
		pointsEach := pointsLeft / len(pointlessChecks)
		for _, check := range pointlessChecks {
			conf.Check[check].Points = pointsEach
		}
		image.TotalPoints += (pointsEach * len(pointlessChecks))
		if image.TotalPoints < 100 {
			for i := 0; image.TotalPoints < 100; image.TotalPoints++ {
				conf.Check[pointlessChecks[i]].Points++
				i++
				if i > len(pointlessChecks)-1 {
					i = 0
				}
			}
			image.TotalPoints += (100 - image.TotalPoints)
		}
	}

	// Reset TotalPoints, since it was only used as a scratch variable and will
	// be calculated again when the checks are run.
	image.TotalPoints = 0
}

// assignDescriptions is automatically assign descriptions to checks that don't
// have one.
func assignDescriptions() {
	for i, check := range conf.Check {
		var msg string
		if check.Message != "" {
			continue
		}
		for _, cond := range check.Pass {
			if msg != "" {
				newMsg := getDesc(cond)
				if newMsg != "" {
					msg += ", and " + strings.ToLower(string(newMsg[0])) + newMsg[1:]
				}
			} else {
				msg = getDesc(cond)
			}
		}
		for _, cond := range check.PassOverride {
			if msg != "" {
				newMsg := getDesc(cond)
				if newMsg != "" {
					msg += ", OR " + strings.ToLower(string(newMsg[0])) + newMsg[1:]
				}
			} else {
				msg = getDesc(cond)
			}
		}
		if msg == "" {
			msg = "Check passed"
		}
		conf.Check[i].Message = msg
	}
}

func getDesc(c cond) string {
	switch c.Type {
	case "Command":
		return "Command \"" + c.Cmd + "\" passed"
	case "CommandNot":
		return "Command \"" + c.Cmd + "\" failed"
	case "CommandOutput":
		return "Command \"" + c.Cmd + "\" had the output \"" + c.Value + "\""
	case "CommandOutputNot":
		return "Command \"" + c.Cmd + "\" did not have the output \"" + c.Value + "\""
	case "CommandContains":
		return "command \"" + c.Cmd + "\" contained output \"" + c.Value + "\""
	case "CommandContainsNot":
		return "Command \"" + c.Cmd + "\" output did not contain \"" + c.Value + "\""
	case "PathExists":
		return "Path \"" + c.Path + "\" exists"
	case "PathExistsNot":
		return "Path \"" + c.Path + "\" does not exist"
	case "FileContains":
		return "File \"" + c.Path + "\" contains regular expression \"" + c.Value + "\""
	case "FileContainsNot":
		return "File \"" + c.Path + "\" does not contain regular expression \"" + c.Value + "\""
	case "DirContains":
		return "Directory \"" + c.Path + "\" contains expression \"" + c.Value + "\""
	case "DirContainsNot":
		return "Directory \"" + c.Path + "\" does not contain expression \"" + c.Value + "\""
	case "FileEquals":
		return "File \"" + c.Path + "\" matches hash"
	case "FileEqualsNot":
		return "File \"" + c.Path + "\" doesn't match hash"
	case "ProgramInstalled":
		return c.Name + " is installed"
	case "ProgramInstalledNot":
		return c.Name + " has been removed"
	case "ServiceUp":
		return "Service \"" + c.Name + "\" is installed and running"
	case "ServiceUpNot":
		return "Service " + c.Name + " has been stopped"
	case "UserExists":
		return "User " + c.User + " has been added"
	case "UserExistsNot":
		return "User " + c.User + " has been removed"
	case "UserInGroup":
		return "User " + c.User + " is in group \"" + c.Group + "\""
	case "UserInGroupNot":
		return "User " + c.User + " removed or is not in group \"" + c.Group + "\""
	case "FirewallUp":
		return "Firewall has been enabled"
	case "FirewallUpNot":
		return "Firewall has been disabled"
	case "ProgramVersion":
		return c.Name + " is version " + c.Value
	case "ProgramVersionNot":
		return c.Name + " is not version " + c.Value

	// Linux checks
	case "AutoCheckUpdatesEnabled":
		return "The system automatically checks for updates daily"
	case "AutoCheckUpdatesEnabledNot":
		return "The system does not automatically checks for updates daily"
	case "GuestDisabledLDM":
		return "Guest is disabled"
	case "GuestDisabledLDMNot":
		return "Guest is enabled"
	case "KernelVersion":
		return "Kernel is version " + c.Value
	case "KernelVersionNot":
		return "Kernel is not version " + c.Value
	case "PermissionIs":
		return "Permissions of " + c.Path + " are " + c.Value
	case "PermissionIsNot":
		return "Permissions of " + c.Path + " are not " + c.Value

	// Windows checks
	case "BitlockerEnabled":
		return "Bitlocker drive encryption has been enabled"
	case "BitlockerEnabledNot":
		return "Bitlocker drive encryption has been disabled"
	case "FileOwner":
		return c.Path + " is owned by " + c.User
	case "FileOwnerNot":
		return c.Path + " is not owned by " + c.User
	case "PasswordChanged":
		return "Password for " + c.User + " has been changed"
	case "PasswordChangedNot":
		return "Password for " + c.User + " has not been changed"
	case "SecurityPolicy":
		return "Security policy option " + c.Key + " is set to " + c.Value
	case "SecurityPolicyNot":
		return "Security policy option " + c.Key + " is not set to " + c.Value
	case "ServiceStartup":
		return c.Name + " has startup type " + c.Value
	case "ServiceStartupNot":
		return c.Name + " does not have startup type " + c.Value
	case "ScheduledTaskExists":
		return "Scheduled task " + c.Name + " exists"
	case "ScheduledTaskExistsNot":
		return "Scheduled task " + c.Name + " doesn't exist"
	case "ShareExists":
		return "Share " + c.Name + " exists"
	case "ShareExistsNot":
		return "Share " + c.Name + " doesn't exist"
	case "RegistryKey":
		return "Registry key " + c.Key + " matches \"" + c.Value + "\""
	case "RegistryKeyNot":
		return "Registry key " + c.Key + " does not match \"" + c.Value + "\""
	case "RegistryKeyExists":
		return "Registry key " + c.Key + " exists"
	case "RegistryKeyExistsNot":
		return "Registry key " + c.Key + " does not exist"
	case "UserDetail":
		return "User property " + c.Key + " for " + c.User + " is equal to \"" + c.Value + "\""
	case "UserDetailNot":
		return "User property " + c.Key + " for " + c.User + " is not equal to \"" + c.Value + "\""
	case "UserRights":
		return "User or group " + c.User + " has privilege \"" + c.Value + "\""
	case "UserRightsNot":
		return "User or group " + c.User + " does not have privilege \"" + c.Value + "\""
	case "WindowsFeature":
		return c.Name + " feature has been enabled"
	case "WindowsFeatureNot":
		return c.Name + " feature has been disabled"

	default:
		warn("Cannot autogenerate message for check type:", c.Type)
		return ""
	}

}
