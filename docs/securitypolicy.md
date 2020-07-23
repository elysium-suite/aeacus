# aeacus

## Windows Security Settings 🤯🔫

> A note on using `secedit.exe` and just parsing it... even [more reputable projects](https://github.com/dsccommunity/SecurityPolicyDsc/blob/8c318e43171cd32b14fe914b9c18c307093ba964/Modules/SecurityPolicyResourceHelper/SecurityPolicyResourceHelper.psm1) found it to be usable solution.

> List is sourced from `secedit.exe` and [this god-awful spreadsheet from Microsoft](https://www.microsoft.com/en-us/download/details.aspx?id=25250).

### Account Policies

### Password Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/0b40db09-d95d-40a6-8467-32aedec8140c

- `MinimumPasswordAge`
- `MaximumPasswordAge`
- `MinimumPasswordLength`
- `PasswordComplexity`
- `ClearTextPassword`
  > Also known as "Store passwords using reversible encryption".
- `PasswordHistorySize`

### Account Lockout Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/2cd39c97-97cd-4859-a7b4-1229dad5f53d

- `ForceLogoffWhenHourExpire`
- `LockoutDuration`
- `LockoutBadCount`
- `ResetLockoutCount`

## Local Policies

### Event Audit

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/01f8e057-f6a8-4d6e-8a00-99bcd241b403

- `AuditSystemEvents`
- `AuditLogonEvents`
- `AuditObjectAccess`
- `AuditPrivilegeUse`
- `AuditPolicyChange`
- `AuditAccountManage`
- `AuditProcessTracking`
- `AuditDSAccess`
- `AuditAccountLogon`

> These should be set to `3` for "Success and Failure Audits".

<hr>

Everything from this point and below should be in the `Security Options` pane of `gpedit.msc` or `secpol.msc`. https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/security-options

### Local Account Policies

> https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpsb/d6eaa54a-f609-48e9-8461-b32738d77a47

- `LSAAnonymousNameLookup`
- `EnableAdminAccount`
- `EnableGuestAccount`
- `NewAdministratorName`
- `NewGuestName`

## Security Options

I'm seriously going to suffer brain damage if I have to format all of these again... [See the spreadsheet](https://docs.google.com/spreadsheets/d/1N7uuke4Jg1R9FBhj8o5dxJQtEntQlea0McYz5upaiTk/edit#gid=1772229936).
