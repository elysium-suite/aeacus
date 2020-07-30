# aeacus

## Checks

This is a list of vulnerability checks that can be used in the configuration for `aeacus`. The notes on this page contain a lot of important information, please be sure to read them.

> **Note!** Each of the commands here can check for the opposite by appending 'Not' to the check type. For example, `FileExistsNot` to pass if a file does not exist.

**Command**: pass if command succeeds (exit code `0`, checks `$?`)

```
type='Command'
arg1='grep "pam_history.so" /etc/pam.d/common-password'
```

> **Note!** Each of these check types can be used for either `Pass` or `Fail` conditions, and there can be multiple conditions per check.

**CommandOutput**: pass if command output matches

```
type='CommandOutput'
arg1='(Get-NetFirewallProfile -Name Domain).Enabled'
arg2='True'
```

**FileExists**: pass if specified file exists

```
type='FileExists'
arg1='C:\importantprogram.exe'
```

> **Note!** You don't have to escape any characters because we're using single quotes, which are literal strings in TOML. If you need use single quotes, use a TOML multi-line string literal `''' like this! that's neat! '''`).

**FileContains**: pass if file contains string

> Note: `FileContains` will never pass if file does not exist! Add an additional pass check for FileExistsNot, for example, if you want to score that a file does not contain a line OR it doesn't exist.

```
type='FileContains'
arg1='/home/coolUser/Desktop/Forensic Question 1.txt'
arg2='ANSWER: SomeCoolAnswer'
```

> **Note!**: Please use absolute paths (rather than relative) for safety and specificity.

**FileContainsRegex**: pass if file contains regex string

```
type='FileContains'
arg1='C:\Users\coolUser\Desktop\Forensic Question 1.txt'
arg2='ANSWER:\sCool[a-zA-Z]+VariedAnswer'
```

> **Note!** A check passes by default. This means that you can use two failing conditions to simulate two conditions that should pass only if they're _BOTH_ true. That's confusing, so here's an example: you want a check to pass only if this file AND that file contain a string. So, instead of two pass conditions (pass if FileContains, pass if FileContains) (where the check will pass if either pass), you can code two fail conditions where the check will pass only if _BOTH_ fail conditions do not pass (fail if FileContainsNot, fail if FileContainsNot).

_Example of a check with two mutually inclusive pass conditions (User1 must be Group1 and Group2):_

```
[[check]]
message = "User1 is in Group1 and Group2"
[[check.fail]]
type="UserInGroupNot"
arg1="User1"
arg2="Group1"
[[check.fail]]
type="UserInGroupNot"
arg1="User1"
arg2="Group2"
```

**DirContainsRegex**: pass if directory contains regex string

```
type='DirContainsRegex'
arg1='/etc/sudoers.d/'
arg2='NOPASSWD'
```

> `DirContainsRegex` is recursive! This means it checks every folder and subfolder. It currently is capped at 10,000 files so it doesn't segfault if you try to search `/`...

**FileEquals**: pass if file equals sha1 hash

```
type='FileEquals'
arg1='/etc/sysctl.conf'
arg2='403926033d001b5279df37cbbe5287b7c7c267fa'
```

> **Note!** If a check has negative points assigned to it, it automatically becomes a penalty.

**PackageInstalled**: pass if package is installed

```
type='PackageInstalled'
arg1='Mozilla Firefox 75 (x64 en-US)'
```

> For packages, Linux uses `dpkg`, Windows uses the Windows API

**ServiceUp**: pass if service is running

```
type='ServiceUp'
arg1='sshd'
```

> For services, Linux uses `systemctl`, Windows uses `Get-Service`

**UserExists**: pass if user exists on system

```
type='UserExists'
arg1='ballen'
```

**UserInGroup**: pass if specified user is in specified group

```
type='UserInGroupNot'
arg1='HackerUser'
arg2='Administrators'
```

> Linux reads `/etc/group` and Windows checks `net user` behind the scenes.

**FirewallUp**: pass if firewall is active

```
type='FirewallUp'
```

> **Note**: On Linux, unfortunately uses `ufw` at the moment. On Window, this passes if all three Windows Firewall profiles are active.

<hr>

### Linux-Specific Checks

**GuestDisabledLDM**: pass if guest is disabled (for LightDM)

```
type='GuestDisabledLDM'
```

**PasswordChanged**: pass if user's hashed password is not in /etc/shadow

```
type='PasswordChanged'
arg1='user'
arg2='password-hash-here'
```

**PackageVersion**: pass if package version is equal to specified

```
type='PackageVersion'
arg1='git'
arg2='1:2.17.1-1ubuntu0.4'
```

> `PackageVersion` checks `uname -r`.

**KernelVersion**: pass if kernel version is equal to specified

```
type='KernelVersion'
arg1='5.4.0-42-generic'
```
> `KernelVersion` checks `dpkg -l | awk '$2=="<PACKAGENAME>" { print $3 }'`.

**AutoCheckUpdatesEnabled**: pass if the system is configured to automatically check for updates

```
type='AutoCheckUpdatesEnabled'
```

<hr>

### Windows-Specific Checks

> todo: allow SID input or auto-translation for system account names that can change (Guest, Administrator)
**PasswordChanged**: pass if user password has changed
```
type='PasswordChanged'
arg1='user'
arg2='01/17/2019 20:57:41'
```
> You should take the value from `Get-LocalUser user | select PasswordLastSet` and use it as `arg2`.
**windowsFeature**: pass if Feature Enabled

```
type='WindowsFeature'
arg1='SMB1Protocol'
```

**UserDetail**: pass if user detail key is equal to value

```
type='UserDetailNot'
arg1='Administrator'
arg2='Password expires'
arg3='Never'
```

> `UserDetail` checks `net user` behind the scenes. [See here for reference](userproperties.md).

**UserRights**: pass if specified user or group has specified privilege

```
type='UserRights'
arg1='Administrators'
arg2='SeTimeZonePrivilege'
```

> A list of URA and Constant Names (which are used in the config) [can be found here](https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/user-rights-assignment).

**ShareExists**: pass if SMB share exists

```
type='ShareExists'
arg1='ADMIN$'
```

> **Note!** Don't use any single quotes (`'`) in your parameters for Windows options like this. If you need to, use a double-quoted string instead (ex. `"Admin's files"`)

**ScheduledTaskExists**: pass if scheduled task exists

```
type='ScheduledTaskExists'
arg1='Disk Cleanup'
```

(WORK IN PROGRESS, dont use)
**StartupProgramExists**: pass if startup program exists

```
type='StartupProgramExists'
arg1='backdoor.exe'
```

> (WIP) **StartupProgramExists** checks the startup folder, Run and RunOnce registry keys, and (other startup methods on windows)

**SecurityPolicy**: pass if key is equal to value

```
type='SecurityPolicy'
arg1='DisableCAD'
arg2='0'
```

> **Note**: If your value should be X or higher (for example, MinimumPasswordAge should be 1 or higher), or if your value should be X or lower (but not 0) (ex. MaximumPasswordAge should be between 1 and 999), the `SecurityPolicy` check will intelligently score it for you.

> Values are checking Registry Keys and `secedit.exe` behind the scenes. This means `0` is `Disabled` and `1` is `Enabled`. [See here for reference](securitypolicy.md).

**RegistryKey**: pass if key is equal to value

```
type='RegistryKey'
arg1='HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\DisableCAD'
arg2='0'
```

> Note: This check will never pass if retrieving the key fails (wrong hive, key doesn't exist, etc). If you want to check that a key was deleted, use `RegistryKeyExistsNot`.

> **Administrative Templates**: There are 4000+ admin template fields. See [this list of registry keys and descriptions](https://docs.google.com/spreadsheets/d/1N7uuke4Jg1R9FBhj8o5dxJQtEntQlea0McYz5upaiTk/edit?usp=sharing), then use the `RegistryKey` or `RegistryKeyExists` check.

**RegistryKeyExists**: pass if key exists

```
type='RegistryKeyExists'
arg1='SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\DisableCAD'
```

> **Note!**: Notice the single quotes `'` on the above argument! This means it's a _string literal_ in TOML. If you don't do this, you have to make sure to escape your slashes (`\` --> `\\`)

> Note: You can use `SOFTWARE` as a shortcut to mean `HKEY_LOCAL_MACHINE\SOFTWARE`.
