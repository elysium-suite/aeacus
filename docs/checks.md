# aeacus

## Checks

This is a list of vulnerability checks that can be used in the configuration for aeacus.


__Command__: pass if command succeeds (exit code `0`)
```
type="Command"
arg1="grep 'pam_history.so' /etc/pam.d/common-password"
```

> __Note!__ Each of these check types can be used for either `Pass` or `Fail` conditions, and there can be multiple conditions per check.

__FileExists__: pass if specified file exists
```
type="FileExists"
arg1="C:\importantprogram.exe"
```

> __Note!__ Each of the commands here can check for the opposite by appending "Not" to the check type. For example, `FileExistsNot` to pass if a file does not exist.

__FileContains__: pass if file contains string
```
type="FileContains"
arg1="/home/coolUser/Desktop/Forensic Question 1.txt"
arg2="ANSWER: SomeCoolAnswer"
```

__FileContainsRegex__: pass if file contains regex string
```
type="FileContains"
arg1="C:\Users\coolUser\Desktop\Forensic Question 1.txt"
arg2="ANSWER:\sCool[a-zA-Z]+VariedAnswer"
```

> __Note!__ A check passes by default. This means that you can use two failing conditions to simulate two conditions that should pass only if they're _BOTH_ true. That's confusing, so here's an example: you want a check to pass only if this file AND that file contain a string. So, instead of two pass conditions (pass if FileContains, pass if FileContains) (where the check will pass if either pass), you can code two fail conditions where the check will pass only if _BOTH_ fail conditions do not pass (fail if FileContainsNot, fail if FileContainsNot).

__FileEquals__: pass if file equals sha1 hash
```
type="FileEquals"
arg1="/etc/sysctl.conf"
arg2="403926033d001b5279df37cbbe5287b7c7c267fa"
```

__PackageInstalled__: pass if package is installed
```
type="PackageInstalled"
arg1="Mozilla Firefox 75 (x64 en-US)"
```

> For packages, Linux side uses `dpkg`, Windows side uses the Windows API

__ServiceUp__: pass if service is running
```
type="ServiceUp"
arg1="sshd"
```

> For services, Linux side uses `systemctl`, Windows side uses `Get-Service`

__UserExists__: pass if user exists on system
```
type="UserExists"
arg1="ballen"
```

> __Note!__ If a check has negative points assigned to it, it automatically becomes a penalty.

__FirewallUp__: pass if firewall is active
```
type="FirewallUp"
```

### Linux-Specific Checks

__UserIsInGroup__: pass if specified user is in specified group
```
type=UserIsInGroup
arg1="ballen"
arg2="sudo"
```

### Windows-Specific Checks

(WORK IN PROGRESS, dont use)
__UserDetail__: pass if user detail key is equal to value
```
type="UserDetail"
arg1="Administrator"
arg2="Enforce_password_history"
arg3=""
```
> Docs will be here so you can see what the detail names are

__UserRights__: pass if specified user or group has specified privilege
```
type="UserRights"
arg1="Administrators"
arg2="SeTimeZonePrivilege"
```

> A list of URA and Constant Names (which are used in the config) [can be found here](https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/user-rights-assignment).

__SecurityPolicy__: pass if key is equal to value
```
type="SecurityPolicy"
arg1="DisableCAD"
arg2="0"
```
> TODO: add specific settings for interactive logon/etc that take relative operators (ex. For Password age, should check if value is x or higher)

> Values are checking Registry Keys and `secedit.exe` behind the scenes. This means `0` is `Disabled` and `1` is `Enabled`. [See here for reference](securitypolicy.md).

__RegistryKey__: pass if key is equal to value
```
type="RegistryKey"
arg1="HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System\\DisableCAD"
arg2="0"
```
> Note: Make sure to escape your slashes (`\` --> `\\`)


(WORK IN PROGRESS dont use) __AdminTemplate__: pass if specified template item is equal to value
```
type="AdminTemplate"
arg1="Turn off background refresh of Group Policy"
arg2="1"
```
> `AdminTemplate` is still very much a work in progress. However, you can check every administrative template item with registry keys (unlike that nasty secpol.) See https://docs.google.com/spreadsheets/d/1N7uuke4Jg1R9FBhj8o5dxJQtEntQlea0McYz5upaiTk/edit?usp=sharing for the registry keys list
