package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	wapi "github.com/iamacarpet/go-win64api"
	wapiShared "github.com/iamacarpet/go-win64api/shared"
)

func (c cond) BitlockerEnabled() (bool, error) {
	status, err := wapi.GetBitLockerConversionStatusForDrive("C:")
	if err == nil {
		if status.ConversionStatus == wapiShared.FULLY_ENCRYPTED ||
			status.ConversionStatus == wapiShared.ENCRYPTION_IN_PROGRESS {
			return true, nil
		}
	}
	return false, nil
}

func (c cond) FileOwner() (bool, error) {
	c.requireArgs("Path", "Name")
	owner, err := shellCommandOutput("(Get-Acl " + c.Path + ").Owner")
	owner = strings.TrimSpace(owner)
	return owner == c.Name, err
}

func (c cond) FirewallUp() (bool, error) {
	fwProfilesInt := []int{wapi.NET_FW_PROFILE2_DOMAIN, wapi.NET_FW_PROFILE2_PRIVATE, wapi.NET_FW_PROFILE2_PUBLIC}
	for profile := range fwProfilesInt {
		profileResult, err := wapi.FirewallIsEnabled(int32(profile))
		if err != nil {
			return false, err
		} else if !profileResult {
			return false, nil
		}
	}
	return true, nil
}

func (c cond) FirewallDefaultBehavior() (bool, error) {
	c.requireArgs("Name", "Key", "Value")
	var profile int32
	var behavior int32
	switch strings.ToLower(c.Name) {
	case "domain":
		profile = wapi.NET_FW_PROFILE2_DOMAIN
	case "private":
		profile = wapi.NET_FW_PROFILE2_PRIVATE
	case "public":
		profile = wapi.NET_FW_PROFILE2_PUBLIC
	default:
		fail("Unknown firewall profile: '" + c.Name + "'")
		return false, errors.New("Unknown firewall profile: " + c.Name)
	}
	switch strings.ToLower(c.Value) {
	case "allow":
		behavior = wapi.NET_FW_ACTION_ALLOW
	case "block":
		behavior = wapi.NET_FW_ACTION_BLOCK
	default:
		fail("Unknown firewall action: '" + c.Value + "'")
		return false, errors.New("Unknown firewall action: " + c.Value)
	}
	switch strings.ToLower(c.Key) {
	case "inbound":
		action, err := wapi.FirewallGetDefaultInboundAction(profile)
		return action == behavior, err
	case "outbound":
		action, err := wapi.FirewallGetDefaultOutboundAction(profile)
		return action == behavior, err
	default:
		fail("Unknown firewall direction: '" + c.Key + "'")
		return false, errors.New("Unknown firewall direction: " + c.Key)
	}
}

// PasswordChanged checks if the password for a given user was changed more
// recently than specified. The date format output by this command is:
//
//	Monday, January 2, 2006 3:04:05 PM
//
// Which somehow manages to defy every common date format. Thanks, Windows.
func (c cond) PasswordChanged() (bool, error) {
	c.requireArgs("User", "After")
	timeStr := "Monday, January 2, 2006 3:04:05 PM"
	configDate, err := time.Parse(timeStr, strings.TrimSpace(c.After))
	if err != nil {
		return false, err
	}
	changed, err := shellCommandOutput(`(Get-LocalUser ` + c.User + `).PasswordLastSet`)
	if err != nil {
		return false, err
	}
	changeDate, err := time.Parse(timeStr, strings.TrimSpace(changed))
	if err != nil {
		return false, err
	}
	return changeDate.After(configDate), nil
}

func (c cond) PermissionIs() (bool, error) {
	c.requireArgs("Path", "Name", "Value")
	permissions, err := getFileRights(c.Path, c.Name)
	if err != nil {
		return false, err
	}
	rights := permissions["filesystemrights"]
	access := permissions["accesscontroltype"]
	return strings.Contains(rights, c.Value) && !strings.EqualFold(access, "Deny"), nil
}

func (c cond) ProgramInstalled() (bool, error) {
	c.requireArgs("Name")
	programList, err := getPrograms()
	if err != nil {
		return false, err
	}
	for _, p := range programList {
		if strings.Contains(p, c.Name) {
			return true, nil
		}
	}
	return false, nil
}

func (c cond) ProgramVersion() (bool, error) {
	c.requireArgs("Name", "Value")
	prog, err := getProgram(c.Name)
	if err != nil {
		return false, err
	}
	return prog.DisplayVersion == c.Value, nil
}

func (c cond) RegistryKey() (bool, error) {
	c.requireArgs("Key", "Value")
	registryArgs := regexp.MustCompile(`\\+`).Split(c.Key, -1)
	if len(registryArgs) < 2 {
		fail("Invalid key for RegistryKey. Did you supply 'key'?")
		return false, errors.New("invalid registry key path: " + c.Key)
	}
	registryHiveText := registryArgs[0]
	keyPath := fmt.Sprintf(strings.Join(registryArgs[1:len(registryArgs)-1], `\`))
	keyLoc := registryArgs[len(registryArgs)-1]

	var registryHive registry.Key
	switch registryHiveText {
	case "HKEY_CLASSES_ROOT", "HKCR":
		registryHive = registry.CLASSES_ROOT
	case "HKEY_CURRENT_USER", "HKCU":
		registryHive = registry.CURRENT_USER
	case "HKEY_LOCAL_MACHINE", "HKLM", "MACHINE":
		registryHive = registry.LOCAL_MACHINE
	case "HKEY_USERS", "HKU":
		registryHive = registry.USERS
	case "HKEY_CURRENT_CONFIG", "HKCC":
		registryHive = registry.CURRENT_CONFIG
	case "SOFTWARE":
		registryHive = registry.LOCAL_MACHINE
		keyPath = `SOFTWARE\` + keyPath
	default:
		fail("Unknown registry hive: " + registryHiveText)
		return false, errors.New("Unknown registry hive " + registryHiveText)
	}

	debug("Getting key", keyPath, "from hive", registryHiveText)
	// Actually get the key
	k, err := registry.OpenKey(registryHive, keyPath, registry.QUERY_VALUE)
	if err != nil {
		if verboseEnabled {
			warn("Registry opening key failed:", err)
		}
		return false, err
	}
	defer k.Close()

	// Fetch registry value
	registrySlice := make([]byte, 256)
	regLength, valType, err := k.GetValue(keyLoc, registrySlice)
	if err != nil {
		// Error is probably about the key not existing. This is fine, some keys
		// are not defined until the setting is explicitly set. However, the
		// check should not pass for RegistryKey or RegistryKeyNot, so we return
		// an error.
		warn("Failed to open registry key:", err)
		return false, err
	}

	registrySlice = registrySlice[:regLength]
	debug("Retrieved registry value was", registrySlice, "length", regLength, "value", valType)

	// Determine value type to convert to string
	var registryValue string
	switch valType {
	case 1: // SZ
		registryValue, _, err = k.GetStringValue(keyLoc)
	case 2: // EXPAND_SZ
		registryValue, _, err = k.GetStringValue(keyLoc)
	case 3: // BINARY
		fail("Binary registry format not yet supported.")
	case 4: // DWORD
		registryValue = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(registrySlice)), 10)
	default:
		fail("Unknown registry type: " + fmt.Sprint(valType))
	}

	// fmt.Printf("Registry value: %s, keyvalue %s\n", registryValue, keyValue)
	if registryValue == c.Value {
		return true, err
	}
	return false, err
}

func (c cond) RegistryKeyExists() (bool, error) {
	c.requireArgs("Key")
	_, err := c.RegistryKey()
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (c cond) ScheduledTaskExists() (bool, error) {
	c.requireArgs("Name")
	return cond{
		Cmd:   "(Get-ScheduledTask -TaskName '" + c.Name + "').TaskName",
		Value: c.Name,
	}.CommandOutput()
}

func (c cond) SecurityPolicy() (bool, error) {
	c.requireArgs("Key", "Value")
	var desiredString string

	// If the passed key is one we know is in the registry, just wrap
	// RegistryKey.
	if regKey, ok := secpolToKey[c.Key]; ok {
		return cond{
			Key:   regKey,
			Value: c.Value,
		}.RegistryKey()
	}

	// Otherwise, we're going to grab and parse secedit output :/
	seceditOutput, err := getSecedit()
	if err != nil {
		return false, err
	}

	re := regexp.MustCompile("(?m)[\r\n]+^.*" + c.Key + ".*$")
	output := strings.TrimSpace(string(re.Find([]byte(seceditOutput))))
	if output == "" {
		return false, errors.New("SecurityPolicy item not found")
	}

	if c.Key == "NewAdministratorName" || c.Key == "NewGuestName" {
		// These two are strings, not numbers, so they have ""
		desiredString = c.Key + " = " + c.Value
	} else if c.Key == "MinimumPasswordAge" ||
		c.Key == "MinimumPasswordLength" ||
		c.Key == "LockoutDuration" ||
		c.Key == "ResetLockoutCount" ||
		c.Key == "MaximumPasswordAge" ||
		c.Key == "LockoutBadCount" ||
		c.Key == "PasswordHistorySize" {

		// These keys are integers, and support ranges.
		var outputResult, err = strconv.Atoi(strings.Split(output, " = ")[1])
		if err != nil {
			return false, err
		}

		if strings.Contains(c.Value, "-") {
			splitVal := strings.Split(c.Value, "-")
			if len(splitVal) != 2 {
				fail("Malformed range value:", c.Value)
				return false, errors.New("invalid c.Value range")
			}
			intLow, err := strconv.Atoi(splitVal[0])
			if err != nil {
				fail(splitVal[0] + " is not a valid integer for SecurityPolicy check")
				return false, err
			}
			intHigh, err := strconv.Atoi(splitVal[1])
			if err != nil {
				fail(splitVal[1] + " is not a valid integer for SecurityPolicy check")
				return false, err
			}
			if intLow <= outputResult && outputResult <= intHigh {
				return true, nil
			}
		} else {
			desiredValue, err := strconv.Atoi(c.Value)
			if err != nil {
				fail(c.Value + " is not a valid integer for SecurityPolicy check")
				return false, errors.New("invalid c.Value")
			}
			if outputResult == desiredValue {
				return true, nil
			}
		}
	} else {
		desiredString = c.Key + " = " + c.Value
	}

	return output == desiredString, nil
}

func (c cond) ServiceStartup() (bool, error) {
	c.requireArgs("Name", "Value")
	var startupNumber string
	switch c.Value = strings.ToLower(c.Value); c.Value {
	case "automatic":
		startupNumber = "2"
	case "manual":
		startupNumber = "3"
	case "disabled":
		startupNumber = "4"
	default:
		fail("Unknown startup type '"+c.Value+"' for", c.Name)
		return false, errors.New("Unknown status type found for " + c.Name)
	}
	serviceKey := `HKLM\SYSTEM\CurrentControlSet\Services\` + c.Name + `\Start`
	return cond{
		Key:   serviceKey,
		Value: startupNumber,
	}.RegistryKey()
}

func (c cond) ServiceUp() (bool, error) {
	c.requireArgs("Name")
	serviceStatus, err := getLocalServiceStatus(c.Name)
	return serviceStatus.IsRunning, err
}

func (c cond) ShareExists() (bool, error) {
	c.requireArgs("Name")
	return cond{
		Cmd:   "(Get-SmbShare -Name '" + c.Name + "').Name",
		Value: c.Name,
	}.CommandOutput()
}

func (c cond) UserExists() (bool, error) {
	c.requireArgs("User")
	user, err := getLocalUser(c.User)
	if err != nil {
		return false, err
	}
	if user.Username == "" {
		return false, nil
	}
	return true, nil
}

func (c cond) UserInGroup() (bool, error) {
	c.requireArgs("User", "Group")
	users, err := wapi.LocalGroupGetMembers(c.Group)
	if err != nil {
		// Error is returned if group is empty.
		return false, nil
	}
	for _, user := range users {
		justName := strings.Split(user.Name, `\`)[1]
		if c.User == user.Name || c.User == justName {
			return true, nil
		}
	}
	return false, nil
}

func (c cond) UserDetail() (bool, error) {
	c.requireArgs("User", "Key", "Value")
	c.Value = strings.TrimSpace(c.Value)
	c.Key = strings.TrimSpace(c.Key)
	splitVal := c.Value
	lookingFor := strings.ToLower(c.Value) == "yes"
	user, err := getLocalUser(c.User)
	if err != nil {
		return false, err
	}
	if c.Key == "PasswordAge" || c.Key == "BadPasswordCount" || c.Key == "NumberOfLogons" {
		var num int
		switch c.Key {
		case "PasswordAge":
			num = int(user.PasswordAge.Hours() / 24)
		case "BadPasswordCount":
			num = int(user.BadPasswordCount)
		case "NumberOfLogons":
			num = int(user.NumberOfLogons)
		}
		if len(c.Value) < 1 {
			fail("Invalid value input:", c.Value)
			return false, errors.New("invalid c.Value range")
		}
		var val int
		switch c.Value[0] {
		case '<':
			splitVal = strings.Split(c.Value, "<")[1]
			val, err = strconv.Atoi(splitVal)
			if err == nil {
				return num < val, nil
			}
		case '>':
			splitVal = strings.Split(c.Value, ">")[1]
			val, err = strconv.Atoi(splitVal)
			if err == nil {
				return num > val, nil
			}
		default:
			val, err = strconv.Atoi(splitVal)
			if err == nil {
				return num == val, nil
			}

		}
		fail("c.Value not an integer:", val)
		return false, err
	}

	//Monday, January 02, 2006 3:04:05 PM
	if c.Key == "LastLogon" {
		lastLogon := user.LastLogon.UTC()
		var timeComparison func(time.Time) bool
		var timeString string
		if len(c.Value) < 2 {
			fail("Could not parse date: \"" + c.Value + "\". Correct format is \"Monday, January 02, 2006 3:04:05 PM\" and in UTC time.")
			return false, errors.New("invalid c.Value date")
		}
		switch c.Value[0] {
		case '<':
			timeString = strings.Split(c.Value, "<")[1]
			timeComparison = lastLogon.Before
		case '>':
			timeString = strings.Split(c.Value, ">")[1]
			timeComparison = lastLogon.After
		default:
			timeComparison = lastLogon.Equal
		}
		parse, err := time.Parse("Monday, January 02, 2006 3:04:05 PM", timeString)
		if err != nil {
			fail("Could not parse date: \"" + c.Value + "\". Correct format is \"Monday, January 02, 2006 3:04:05 PM\" and in UTC time.")
			return false, err
		}
		return timeComparison(parse), nil
	}

	switch c.Key {
	case "FullName":
		if user.FullName == c.Value {
			return true, nil
		}
	case "IsEnabled":
		return user.IsEnabled == lookingFor, nil
	case "IsLocked":
		return user.IsLocked == lookingFor, nil
	case "IsAdmin":
		return user.IsAdmin == lookingFor, nil
	case "PasswordNeverExpires":
		return user.PasswordNeverExpires == lookingFor, nil
	case "NoChangePassword":
		return user.NoChangePassword == lookingFor, nil
	default:
		fail("c.Key (" + c.Key + ") passed to userDetail is invalid.")
		return false, errors.New("invalid detail")
	}
	return false, nil
}

func (c cond) UserRights() (bool, error) {
	// TODO consider /mergedpolicy when windows domain is active?
	// domain support is untested, it should be easy to add a domain
	// flag in the config though. then just make sure you're not getting
	// invalid local policies instead of gpo
	c.requireArgs("Name", "Value")

	// TODO: only get section of users -- this can also falsely score correct for other secedit fields (like LegalNoticeText)
	seceditOutput, err := getSecedit()
	if err != nil {
		return false, err
	}

	re := regexp.MustCompile("(?m)[\r\n]+^.*" + c.Value + ".*$")
	privilegeString := strings.TrimSpace(string(re.Find([]byte(seceditOutput))))
	debug("Privilege string for UserRights is:", privilegeString)
	if privilegeString == "" {
		return false, nil
	}

	if strings.Contains(privilegeString, c.Name) {
		// Sometimes, Windows just puts their user or group name instead of the
		// SID. Really cool
		return true, nil
	}

	privStringSplit := strings.Split(privilegeString, " ")
	if len(privStringSplit) != 3 {
		return false, errors.New("error splitting privilege")
	}

	privStringSplit = strings.Split(privStringSplit[2], ",")
	for _, sidValue := range privStringSplit {
		sidValue = strings.TrimSpace(sidValue)
		userForSid := strings.Split(sidToLocalUser(sidValue[1:]), "\\")
		userSid := strings.TrimSpace(userForSid[0])
		if len(userForSid) == 2 {
			userSid = strings.TrimSpace(userForSid[1])
		}
		if userSid == c.Name {
			return true, nil
		}
	}

	return false, err
}

func (c cond) WindowsFeature() (bool, error) {
	c.requireArgs("Name")
	return cond{
		Cmd:   "(Get-WindowsOptionalFeature -Online -FeatureName " + c.Name + ").State",
		Value: "Enabled",
	}.CommandOutput()
}
