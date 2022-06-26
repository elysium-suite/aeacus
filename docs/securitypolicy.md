# aeacus

## Windows Security Settings

> A note on using `secedit.exe` and just parsing it... even [more reputable projects](https://github.com/dsccommunity/SecurityPolicyDsc/blob/8c318e43171cd32b14fe914b9c18c307093ba964/Modules/SecurityPolicyResourceHelper/SecurityPolicyResourceHelper.psm1) found it to be usable solution.

> List is sourced from `secedit.exe` and [this god-awful spreadsheet from Microsoft](https://www.microsoft.com/en-us/download/details.aspx?id=25250).

These are all aliases. You can (and probably should) use the RegistryKey check for more control and to score obscure policies.

## Account Policies

### Password Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b40db09-d95d-40a6-8467-32aedec8140c

-   `MinimumPasswordAge`
-   `MaximumPasswordAge`
-   `MinimumPasswordLength`
-   `PasswordComplexity`
-   `ClearTextPassword`
	> Also known as "Store passwords using reversible encryption".
-   `PasswordHistorySize`

### Account Lockout Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/2cd39c97-97cd-4859-a7b4-1229dad5f53d

-   `ForceLogoffWhenHourExpire`
-   `LockoutDuration`
-   `LockoutBadCount`
-   `ResetLockoutCount`

## Local Policies

### Event Audit

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/01f8e057-f6a8-4d6e-8a00-99bcd241b403

-   `AuditSystemEvents`
-   `AuditLogonEvents`
-   `AuditObjectAccess`
-   `AuditPrivilegeUse`
-   `AuditPolicyChange`
-   `AuditAccountManage`
-   `AuditProcessTracking`
-   `AuditDSAccess`
-   `AuditAccountLogon`

> These should be set to `3` for "Success and Failure Audits".

<hr>

Everything from this point and below should be in the `Security Options` pane of `gpedit.msc` or `secpol.msc`. https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/security-options

### Local Account Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/d6eaa54a-f609-48e9-8461-b32738d77a47

-   `LSAAnonymousNameLookup`
-   `EnableAdminAccount`
-   `EnableGuestAccount`
-   `NewAdministratorName`
-   `NewGuestName`


### Other Options

Welcome to hell.

| Aeacus Key Name              | Policy Name                                                                                                        | Registry Key                                                                                                     |
|------------------------------|--------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|
| LimitBlankPasswordUse        | Accounts: Limit local account use of blank passwords to console logon only                                         | `MACHINE\System\CurrentControlSet\Control\Lsa\LimitBlankPasswordUse`                                             |
| AuditBaseObjects             | Audit: Audit the accesss of global system objects                                                                  | `9:44:09 PM`                                                                                                     |
| FullPrivilegeAuditing        | Audit: Audit the use of Backup and Restore privilege                                                               | `MACHINE\System\CurrentControlSet\Control\Lsa\FullPrivilegeAuditing`                                             |
| SCENoApplyLegacyAuditPolicy  | Audit: Force audit policy subcategory settings (Windows Vista or later) to override audit policy category settings | `MACHINE\System\CurrentControlSet\Control\Lsa\SCENoApplyLegacyAuditPolicy`                                       |
| CrashOnAuditFail             | Audit: Shut down system immediately if unable to log security audits                                               | `MACHINE\System\CurrentControlSet\Control\Lsa\CrashOnAuditFail`                                                  |
| MachineAccessRestriction     | DCOM: Machine Access Restrictions in Security Descriptor Definition Language (SDDL) syntax                         | `MACHINE\SOFTWARE\policies\Microsoft\windows NT\DCOM\MachineAccessRestriction`                                   |
| MachineLaunchRestriction     | DCOM: Machine Launch Restrictions in Security Descriptor Definition Language (SDDL) syntax                         | `MACHINE\SOFTWARE\policies\Microsoft\windows NT\DCOM\MachineLaunchRestriction`                                   |
| UndockWithoutLogon           | Devices: Allow undock without having to log on                                                                     | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\UndockWithoutLogon`                           |
| AllocateDASD                 | Devices: Allowed to format and eject removable media                                                               | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\AllocateDASD`                                     |
| AddPrinterDrivers            | Devices: Prevent users from installing printer drivers                                                             | `MACHINE\System\CurrentControlSet\Control\Print\Providers\LanMan Print Services\Servers\AddPrinterDrivers`       |
| AllocateCDRoms               | Devices: Restrict CD-ROM access to locally logged-on user only                                                     | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\AllocateCDRoms`                                   |
| AllocateFloppies             | Devices: Restrict floppy access to locally logged-on user only                                                     | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\AllocateFloppies`                                 |
| SubmitControl                | Domain controller: Allow server operators to schedule tasks                                                        | `MACHINE\System\CurrentControlSet\Control\Lsa\SubmitControl`                                                     |
| LDAPServerIntegrity          | Domain controller: LDAP server signing requirements                                                                | `MACHINE\System\CurrentControlSet\Services\NTDS\Parameters\LDAPServerIntegrity`                                  |
| RefusePasswordChange         | Domain controller: Refuse machine account password changes                                                         | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\RefusePasswordChange`                             |
| RequireSignOrSeal            | Domain member: Digitally encrypt or sign secure channel data (always)                                              | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\RequireSignOrSeal`                                |
| SealSecureChannel            | Domain member: Digitally encrypt secure channel data (when possible)                                               | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\SealSecureChannel`                                |
| SignSecureChannel            | Domain member: Digitally sign secure channel data (when possible)                                                  | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\SignSecureChannel`                                |
| DisablePasswordChange        | Domain member: Disable machine account password changes                                                            | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\DisablePasswordChange`                            |
| MaximumPasswordAge           | Domain member: Maximum machine account password age                                                                | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\MaximumPasswordAge`                               |
| RequireStrongKey             | Domain member: Require strong (Windows 2000 or later) session key                                                  | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\RequireStrongKey`                                 |
| DontDisplayLockedUserId      | Interactive Logon: Display user information when session is locked                                                 | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System, value=DontDisplayLockedUserId`               |
| DisableCAD                   | Interactive logon: Do not require CTRL+ALT+DEL                                                                     | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\DisableCAD`                                   |
| DontDisplayLastUserName      | Interactive logon: Don't display last signed-in                                                                    | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\DontDisplayLastUserName`                      |
| LegalNoticeText              | Interactive logon: Message text for users attempting to logon                                                      | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\LegalNoticeText`                              |
| LegalNoticeCaption           | Interactive logon: Message title for users attempting to logon                                                     | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\LegalNoticeCaption`                           |
| CachedLogonsCount            | Interactive logon: Number of previous logons to cache (in case domain controller is not available)                 | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\CachedLogonsCount`                                |
| PasswordExpiryWarning        | Interactive logon: Prompt user to change password before expiration                                                | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\PasswordExpiryWarning`                            |
| ForceUnlockLogon             | Interactive logon: Require Domain Controller authentication to unlock workstation                                  | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\ForceUnlockLogon`                                 |
| ScForceOption                | Interactive logon: Require smart card                                                                              | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\ScForceOption`                                |
| ScRemoveOption               | Interactive logon: Smart card removal behavior                                                                     | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Winlogon\ScRemoveOption`                                   |
| EnableSecuritySignature      | Microsoft network client: Digitally sign communications (if server agrees)                                         | `MACHINE\System\CurrentControlSet\Services\LanmanWorkstation\Parameters\EnableSecuritySignature`                 |
| EnablePlainTextPassword      | Microsoft network client: Send unencrypted password to third-party SMB servers                                     | `MACHINE\System\CurrentControlSet\Services\LanmanWorkstation\Parameters\EnablePlainTextPassword`                 |
| AutoDisconnect               | Microsoft network server: Amount of idle time required before suspending session                                   | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\AutoDisconnect`                               |
| RequireSecuritySignature     | Microsoft network server: Digitally sign communications (always)                                                   | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\RequireSecuritySignature`                     |
| EnableForcedLogOff           | Microsoft network server: Disconnect clients when logon hours expire                                               | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\EnableForcedLogOff`                           |
| SmbServerNameHardeningLevel  | Microsoft network server: Server SPN target name validation level                                                  | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\SmbServerNameHardeningLevel`                  |
| RestrictAnonymousSAM         | Network access: Do not allow anonymous enumeration of SAM accounts                                                 | `MACHINE\System\CurrentControlSet\Control\Lsa\RestrictAnonymousSAM`                                              |
| RestrictAnonymous            | Network access: Do not allow anonymous enumeration of SAM accounts and shares                                      | `MACHINE\System\CurrentControlSet\Control\Lsa\RestrictAnonymous`                                                 |
| DisableDomainCreds           | Network access: Do not allow storage of passwords and credentials for network authentication                       | `MACHINE\System\CurrentControlSet\Control\Lsa\DisableDomainCreds`                                                |
| EveryoneIncludesAnonymous    | Network access: Let Everyone permissions apply to anonymous users                                                  | `MACHINE\System\CurrentControlSet\Control\Lsa\EveryoneIncludesAnonymous`                                         |
| NullSessionPipes             | Network access: Named Pipes that can be accessed anonymously                                                       | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\NullSessionPipes`                             |
| Machine                      | Network access: Remotely accessible registry paths                                                                 | `MACHINE\System\CurrentControlSet\Control\SecurePipeServers\Winreg\AllowedPaths\Machine`                         |
| N/A                          | Network access: Remotely accessible registry paths and sub-paths                                                   | `MACHINE\System\CurrentControlSet\Control\SecurePipeServers\Winreg\AllowedPaths\Machine`                         |
| NullSessionShares            | Network access: Restrict anonymous access to Named Pipes and Shares                                                | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\NullSessionShares`                            |
| NullSessionShares            | Network access: Shares that can be accessed anonymously                                                            | `MACHINE\System\CurrentControlSet\Services\LanManServer\Parameters\NullSessionShares`                            |
| ForceGuest                   | Network access: Sharing and security model for local accounts                                                      | `MACHINE\System\CurrentControlSet\Control\Lsa\ForceGuest`                                                        |
| UseMachineId                 | Network security: Allow Local System to use computer identity for NTLM                                             | `MACHINE\System\CurrentControlSet\Control\Lsa\UseMachineId`                                                      |
| allownullsessionfallback     | Network security: Allow LocalSystem NULL session fallback                                                          | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\allownullsessionfallback`                                   |
| AllowOnlineID                | Network security: Allow PKU2U authentication requests to this computer to use online identities.                   | `MACHINE\System\CurrentControlSet\Control\Lsa\pku2u\AllowOnlineID`                                               |
| SupportedEncryptionTypes     | Network security: Configure encryption types allowed for Kerberos                                                  | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\Kerberos\Parameters\SupportedEncryptionTypes` |
| NoLMHash                     | Network security: Do not store LAN Manager hash value on next password change                                      | `MACHINE\System\CurrentControlSet\Control\Lsa\NoLMHash`                                                          |
| LmCompatibilityLevel         | Network security: LAN Manager authentication level                                                                 | `MACHINE\System\CurrentControlSet\Control\Lsa\LmCompatibilityLevel`                                              |
| LDAPClientIntegrity          | Network security: LDAP client signing requirements                                                                 | `MACHINE\System\CurrentControlSet\Services\LDAP\LDAPClientIntegrity`                                             |
| NTLMMinClientSec             | Network security: Minimum session security for NTLM SSP based (including secure RPC) clients                       | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\NTLMMinClientSec`                                           |
| NTLMMinServerSec             | Network security: Minimum session security for NTLM SSP based (including secure RPC) servers                       | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\NTLMMinServerSec`                                           |
| RestrictNTLMInDomain         | Network security: Restrict NTLM:  NTLM authentication in this domain                                               | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\RestrictNTLMInDomain`                             |
| ClientAllowedNTLMServers     | Network security: Restrict NTLM: Add remote server exceptions for NTLM authentication                              | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\ClientAllowedNTLMServers`                                   |
| DCAllowedNTLMServers         | Network security: Restrict NTLM: Add server exceptions in this domain                                              | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\DCAllowedNTLMServers`                             |
| AuditReceivingNTLMTraffic    | Network security: Restrict NTLM: Audit Incoming NTLM Traffic                                                       | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\AuditReceivingNTLMTraffic`                                  |
| AuditNTLMInDomain            | Network security: Restrict NTLM: Audit NTLM authentication in this domain                                          | `MACHINE\System\CurrentControlSet\Services\Netlogon\Parameters\AuditNTLMInDomain`                                |
| RestrictReceivingNTLMTraffic | Network security: Restrict NTLM: Incoming NTLM traffic                                                             | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\RestrictReceivingNTLMTraffic`                               |
| RestrictSendingNTLMTraffic   | Network security: Restrict NTLM: Outgoing NTLM traffic to remote servers                                           | `MACHINE\System\CurrentControlSet\Control\Lsa\MSV1_0\RestrictSendingNTLMTraffic`                                 |
| SecurityLevel                | Recovery console: Allow automatic administrative logon                                                             | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Setup\RecoveryConsole\SecurityLevel`                       |
| SetCommand                   | Recovery console: Allow floppy copy and access to all drives and all folders                                       | `MACHINE\Software\Microsoft\Windows NT\CurrentVersion\Setup\RecoveryConsole\SetCommand`                          |
| ShutdownWithoutLogon         | Shutdown: Allow system to be shut down without having to log on                                                    | `MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System\ShutdownWithoutLogon`                         |
| ClearPageFileAtShutdown      | Shutdown: Clear virtual memory pagefile                                                                            | `MACHINE\System\CurrentControlSet\Control\Session Manager\Memory Management\ClearPageFileAtShutdown`             |
| ForceKeyProtection           | System cryptography: Force strong key protection for user keys stored on the computer                              | `MACHINE\Software\Policies\Microsoft\Cryptography\ForceKeyProtection`                                            |
| FIPSAlgorithmPolicy          | System cryptography: Use FIPS compliant algorithms for encryption, hashing, and signing                            | `MACHINE\System\CurrentControlSet\Control\Lsa\FIPSAlgorithmPolicy`                                               |
| NoDefaultAdminOwner          | System objects: Default owner for objects created by members of the Administrators group                           | `MACHINE\System\CurrentControlSet\Control\Lsa\NoDefaultAdminOwner`                                               |
| ObCaseInsensitive            | System objects: Require case insensitivity for non-Windows subsystems                                              | `MACHINE\System\CurrentControlSet\Control\Session Manager\Kernel\ObCaseInsensitive`                              |
| ProtectionMode               | System objects: Strengthen default permissions of internal system objects (e.g., Symbolic Links)                   | `MACHINE\System\CurrentControlSet\Control\Session Manager\ProtectionMode`                                        |
| optional                     | System settings: Optional subsystems                                                                               | `MACHINE\System\CurrentControlSet\Control\Session Manager\SubSystems\optional`                                   |
| AuthenticodeEnabled          | System settings: Use Certificate Rules on Windows Executables for Software Restriction Policies                    | `MACHINE\Software\Policies\Microsoft\Windows\Safer\CodeIdentifiers\AuthenticodeEnabled`                          |
| FilterAdministratorToken     | User Account Control: Admin Approval Mode for the Built-in Administrator account                                   | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\FilterAdministratorToken`                             |
| EnableUIADesktopToggle       | User Account Control: Allow UIAccess applications to prompt for elevation without using the secure desktop.        | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\EnableUIADesktopToggle`                               |
| ConsentPromptBehaviorAdmin   | User Account Control: Behavior of the elevation prompt for administrators in Admin Approval Mode                   | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\ConsentPromptBehaviorAdmin`                           |
| ConsentPromptBehaviorUser    | User Account Control: Behavior of the elevation prompt for standard users                                          | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\ConsentPromptBehaviorUser`                            |
| EnableInstallerDetection     | User Account Control: Detect application installations and prompt for elevation                                    | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\EnableInstallerDetection`                             |
| ValidateAdminCodeSignatures  | User Account Control: Only elevate executables that are signed and validated                                       | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\ValidateAdminCodeSignatures`                          |
| EnableSecureUIAPaths         | User Account Control: Only elevate UIAccess applications that are installed in secure locations                    | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\EnableSecureUIAPaths`                                 |
| EnableLUA                    | User Account Control: Run all administrators in Admin Approval Mode                                                | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\EnableLUA`                                            |
| PromptOnSecureDesktop        | User Account Control: Switch to the secure desktop when prompting for elevation                                    | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\PromptOnSecureDesktop`                                |
| EnableVirtualization         | User Account Control: Virtualize file and registry write failures to per-user locations                            | `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\EnableVirtualization`                                 |
