# aeacus

## Checks

This is a list of vulnerability checks that can be used in the configuration for `aeacus`. The notes on this page contain a lot of important information, please be sure to read them.

> **Note!** Each of the commands here can check for the opposite by appending 'Not' to the check type. For example, `PathExistsNot` to pass if a file does not exist.

> **Note!** If a check has negative points assigned to it, it automatically becomes a penalty.

> **Note!** Each of these check types can be used for `Pass`, `PassOverride` or `Fail` conditions, and there can be multiple conditions per check. See [configuration](config.md) for more details.

> Note: `Command*` checks are prone to interception, modification, and tomfoolery. Your scoring configuration will be much more robust if you rely on checks using native mechanisms rather than shell commands (for example, `PathExists` instead of ls).

**CommandContains**: pass if command output contains string. If it returns an error, check never passes. Use of this check is discouraged.

```
type = 'CommandContains'
cmd = 'ufw status'
value = 'Status: active'
```

> **Note!** If any check returns an error (e.g., something that it was not expecting), it will _never_ pass, even if it's a `Not` condition. This varies by check, but for example, if you try to check the content of a file that doesn't exist, it will return an error and not succeed-- even if you were doing `FileContainsNot`.

**CommandOutput**: pass if command output matches string exactly. If it returns an error, check never passes. Use of this check is discouraged.

```
type = 'CommandOutput'
cmd = '(Get-NetFirewallProfile -Name Domain).Enabled'
value = 'True'
```

**DirContains**: pass if directory contains regular expression (regex) string

> **Note!** Read more about regex [here](regex.md).

```
type = 'DirContains'
path = '/etc/sudoers.d/'
value = 'NOPASSWD'
```

> `DirContains` is recursive! This means it checks every folder and subfolder. It currently is capped at 10,000 files, so you should begin your search at the deepest folder possible.


> **Note!** You don't have to escape any characters because we're using single quotes, which are literal strings in TOML. If you need use single quotes, use a TOML multi-line string literal `''' like this! that's neat! C:\path\here '''`), or just normal quotes (but you'll have to escape characters with those).

**FileContains**: pass if file contains regex

> Note: `FileContains` will never pass if file does not exist! Add an additional PassOverride check for PathExistsNot, if you want to score that a file does not contain a line, OR it doesn't exist.

```
type = 'FileContains'
path = 'C:\Users\coolUser\Desktop\Forensic Question 1.txt'
value = 'ANSWER:\sCool[a-zA-Z]+VariedAnswer'
```

**FileEquals**: pass if file equals sha256 hash

```
type = 'FileEquals'
path = '/etc/sysctl.conf'
name = 'e61ff3fb83b51fe9f2cd03cc0408afa15d4e8e69b8488b4ed1ecb854ae25da9b'
```

**FirewallUp**: pass if firewall is active

```
type = 'FirewallUp'
```

> **Note**: On Linux, only `ufw` is supported (checks `/etc/ufw/ufw.conf`). On Window, this passes if all three Windows Firewall profiles are active.


**PathExists**: pass if specified path exists. This works for both files AND folders (directories).

```
type = 'PathExists'
path = '/var/www/backup.zip'
```

```
type = 'PathExists'
path = 'C:\importantfolder\'
```

> **Note!**: Please use absolute paths (rather than relative) for safety and specificity.

**ProgramInstalled**: pass if program is installed. On Linux, will use `dpkg`, and on Windows, checks if any installed programs contain your program string.

```
type = 'ProgramInstalled'
name = 'Mozilla Firefox 75 (x64 en-US)'

```

**ProgramVersion**: pass if a program meets the version requirements

```
# Linux: get version from dpkg -s programhere
type = 'ProgramVersion'
name = 'Firefox'
value = '88.0.1+build1-0ubuntu0.20.04.2'
```

```
# Windows: get versions from .\aeacus.exe info programs
# Checks version on first matching substring. E.g., for program name 'Ace',
# it may match on 'Ace Of Spades' rather than 'Ace Ventura'. Make your program
# name as detailed as possible.
type = 'ProgramVersion'
name = 'Firefox'
value = '95.0.1'
```

> **Note**: We recommend you use the `Not` version of this check to score a program's version being different from its version at the beginning of the image. You can't guarantee that the latest version of the program you're scoring will be the same once your round is released, and it's unlikely that a competitor will intentionally downgrade a package.

> For packages, Linux uses `dpkg`, Windows uses the Windows API

**ServiceUp**: pass if service is running

```
type = 'ServiceUp'
name = 'sshd'
```

```
# Windows: check the service 'Properties' to find the real service name
type = 'ServiceUp'
name = 'tapisrv' # this is telephony

```

> For services, Linux uses `systemctl`, Windows uses `Get-Service`

**UserExists**: pass if user exists on system

```
type = 'UserExists'
user = 'ballen'
```

**UserInGroup**: pass if specified user is in specified group

```
type = 'UserInGroupNot'
user = 'HackerUser'
group = 'Administrators'
```

> Linux reads `/etc/group` and Windows uses the Windows API.

<hr>

### Linux-Specific Checks

**AutoCheckUpdatesEnabled**: pass if the system is configured to automatically check for updates

```
type = 'AutoCheckUpdatesEnabled'
```

**Command**: pass if command succeeds. Use of this check is discouraged. This check will NOT return an error if the command is not found

```
type = 'Command'
cmd = 'cat coolfile.txt'
```

**GuestDisabledLDM**: pass if guest is disabled (for LightDM)

```
type = 'GuestDisabledLDM'
```

**KernelVersion**: pass if kernel version is equal to specified

```
type = 'KernelVersion'
value = '5.4.0-42-generic'
```

> Tip: Check your `KernelVersion` with `uname -r`. This check performs the `uname` syscall.

> Only works for standard `apt` installs.

**PasswordChanged**: pass if user's password hash is not next to their username in `/etc/shadow`. If you don't use the whole hash, make sure you start it from the beginning (typically `$X$...` where X is a number).

```
type = 'PasswordChanged'
user = 'bob'
value = '$6$BgBsRlajjwVOoQCY$rw5WBSha4nkpynzfCzc3yYkV1OyDhr.ELoJOPpidwZoygUzRFBFSrtE3fyP0ITubCwN9Bb9DUqVV3mzTHL8sw/'
```

> This check will never pass if the user does not exist, so don't use this with users that should be removed.

**PermissionIs**: pass if the specified file has octal permissions specified. Use question marks to omit bits you don't care about.

```
type = 'PermissionIs'
path = '/etc/shadow
value = 'rw-rw----'
```

For example, this one checks that /bin/bash is not SUID and not world writable:
```
type = 'PermissionIsNot'
path = '/bin/bash'
value = 's???????w?'
```

<hr>

### Windows-Specific Checks

**BitlockerEnabled**: pass if a drive has been fully encrypted with bitlocker drive encription or is in the process of being encrypted

```
type = "BitlockerEnabled"
```
> This check will succeed if the drive is either encrypted or encryption is in progress.

**FileOwner**: pass if specified user/group owns a given file

```
type = 'FileOwner'
path = 'C:\test.txt'
name = 'BUILTIN\Administrators'
```

> Get owner of the file using PowerShell: `(Get-Acl [FILENAME]).Owner`

**FirewallDefaultBehavior**: pass if the firewall profile's default behavior is set to the specified value

```
type = 'FirewallDefaultBehavior'
name = 'Domain'
value = 'Allow'
key = 'Inbound'
```
> Valid "name" (profile) values are: Domain, Public, Private, All
>
> Valid "value" (behavior) values are: Allow, Block
>
> Valid "key" (direction) values are: Inbound, Outbound



**PasswordChanged**: pass if user password has changed after the specified date

```
type = 'PasswordChanged'
user = 'username'
after = 'Monday, January 02, 2006 3:04:05 PM'
```
> You should take the value from `(Get-LocalUser <USERNAME>).PasswordLastSet` and use it as `after`. This check will never pass if the user does not exist, so don't use this with users that should be removed.

**RegistryKey**: pass if key is equal to value

```
type = 'RegistryKey'
key = 'HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\DisableCAD'
value = '0'
```

> Note: This check will never pass if retrieving the key fails (wrong hive, key doesn't exist, etc). If you want to check that a key was deleted, use `RegistryKeyExists`.

> **Administrative Templates**: There are 4000+ admin template fields. See [this list of registry keys and descriptions](https://docs.google.com/spreadsheets/d/1N7uuke4Jg1R9FBhj8o5dxJQtEntQlea0McYz5upaiTk/edit?usp=sharing), then use the `RegistryKey` or `RegistryKeyExists` check.

**RegistryKeyExists**: pass if key exists

```
type = 'RegistryKeyExists'
key = 'SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\DisableCAD'
```

> **Note!**: Notice the single quotes `'` on the above argument! This means it's a _string literal_ in TOML. If you don't do this, you have to make sure to escape your slashes (`\` --> `\\`)

> Note: You can use `SOFTWARE` as a shortcut for `HKEY_LOCAL_MACHINE\SOFTWARE`.

**ScheduledTaskExists**: pass if scheduled task exists

```
type = 'ScheduledTaskExists'
name = 'Disk Cleanup'
```

**SecurityPolicy**: pass if key is within the bounds for value

```
type = 'SecurityPolicy'
key = 'DisableCAD'
value = '0'
```

> Values are checking Registry Keys and `secedit.exe` behind the scenes. This means `0` is `Disabled` and `1` is `Enabled`. [See here for reference](securitypolicy.md).

> **Note**: For all integer-based values (such as `MinimumPasswordAge`), you can provide a range of values, as seen below. The lower value must be specified first.

```
type = 'SecurityPolicy'
key = 'MaximumPasswordAge'
value = '80-100'
```

**ServiceStartup**: pass if service is set to a given startup type (manual, automatic, or disabled)

```
type = "ServiceStartup"
name = "TermService"
value = "manual"
```

> This check is a wrapper around RegistryKey to fetch the proper key for you. Also, Automatic (Delayed) and Automatic are the same value for the key we're checking.

**ShareExists**: pass if SMB share exists

```
type = 'ShareExists'
name = 'ADMIN$'
```

> **Note!** Don't use any single quotes (`'`) in your parameters for Windows options like this. If you need to, use a double-quoted string instead (ex. `"Admin's files"`)


**UserDetail**: pass if user detail key is equal to value

> **Note!** The valid boolean values for this command (when the field is only True or False) are 'yes', if you want the value to be true, or literally anything else for false (like 'no').

```
type = 'UserDetailNot'
user = 'Administrator'
key = 'PasswordNeverExpires'
value = 'No'
```

> See [here](userproperties.md) for all `UserDetail` properties.

> **Note!** For non-boolean details, you can use modifiers in the value field to specify the comparison.
> This is specified in the above property document.
```
type = 'UserDetail'
user = 'Administrator'
key = 'PasswordAge'
value = '>90'
```


**UserRights**: pass if specified user or group has specified privilege

```
type = 'UserRights'
name = 'Administrators'
value = 'SeTimeZonePrivilege'
```

> A list of URA and Constant Names (which are used in the config) [can be found here](https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/user-rights-assignment). On your local machine, check Local Security Policy > User Rights Assignments to see the current assignments.


**WindowsFeature**: pass if Windows Feature is enabled

```
type = 'WindowsFeature'
name = 'SMB1Protocol'
```

> **Note:** Use the PowerShell tool `Get-WindowsOptionalFeature -Online` to find the feature you want!
