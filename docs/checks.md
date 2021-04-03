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

**CommandOutput**: pass if command output matches exactly. if error, never passes

```
type='CommandOutput'
arg1='(Get-NetFirewallProfile -Name Domain).Enabled'
arg2='True'
```

**CommandContains**: pass if command output contains string. if error, never passes

```
type='CommandContains'
arg1='firewall status'
arg2='Active'
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
type='FileContainsRegex'
arg1='C:\Users\coolUser\Desktop\Forensic Question 1.txt'
arg2='ANSWER:\sCool[a-zA-Z]+VariedAnswer'
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

**PasswordChanged**: pass if user's hashed password is not in `/etc/shadow`

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

> `PackageVersion` checks `dpkg -l | awk '$2=="<PACKAGENAME>" { print $3 }'`.

**KernelVersion**: pass if kernel version is equal to specified

```
type='KernelVersion'
arg1='5.4.0-42-generic'
```

> `KernelVersion` checks `uname -r`.

**AutoCheckUpdatesEnabled**: pass if the system is configured to automatically check for updates

```
type='AutoCheckUpdatesEnabled'
```

> Only works for standard `apt` installs.

**PermissionIs**: pass if the specified file has the octal permissions specified

```
type='PermissionIs'
arg1='/etc/passwd'
arg2='644'
```

<hr>

### Windows-Specific Checks

**ServiceStatus**: pass if service status and service startup type is the same as specified

```
type="ServiceStatus"
arg1="TermService"
arg2="Running"
arg3="Automatic"
```

> This check uses the windows API to check the service current status and the windows registry for the startuptype

> todo: allow SID input or auto-translation for system account names that can change (Guest, Administrator)

**PasswordChanged**: pass if user password has changed after the given date

```
type='PasswordChanged'
arg1='user'
arg2='01/17/2019 20:57:41'
```

> You should take the value from `Get-LocalUser user | select PasswordLastSet` and use it as `arg2`.

**WindowsFeature**: pass if Feature Enabled

```
type='WindowsFeature'
arg1='SMB1Protocol'
```

**UserDetail**: pass if user detail key is equal to value

```
type='UserDetailNot'
arg1='Administrator'
arg2='PasswordNeverExpires'
arg3='No'
```

> See [here](userproperties.md) for all `UserDetail` properties.

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

> Note: You can use `SOFTWARE` as a shortcut for `HKEY_LOCAL_MACHINE\SOFTWARE`.

**FileOwner**: pass if specified owner is the owner of the specified file

```
type='FileOwner'
arg1='C:\test.txt'
arg2='BUILTIN\Administrators'
```

> Get owner of the file using `(Get-Acl [FILENAME]).Owner`.

**ProgramVersion**: pass if a program meets the version requirements

```
type='ProgramVersion'
arg1='Firefox'
arg2='87'
arg3='1'
```

> Variable `arg3` is the version comparing method. 0 is equal to, 1 is greater than, 2 is equal to or greater than.
> `arg2` is the version that you are comparing the program to.

**ProgramVersionNot**: pass if a program does not meet the version requirements

```
type='ProgramVersion'
arg1='Firefox'
arg2='87'
arg3='1'
```

> Variable `arg3` is the version comparing method. 0 is equal to, 1 is greater than, 2 is equal to or greater than.
> `arg2` is the version that you are comparing the program to.
