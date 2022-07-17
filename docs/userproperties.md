# aeacus

## Windows User Properties

These are the user properties that you would normally find the user edit dialog. (Plus some extra goodies!)

**Usernames are case-sensitive.** `Yes`/`No` are not case-sensitive. If you input something other than "yes" or "no" for a boolean check, it defaults to "no".

- `FullName`: Full name of user (case-sensitive)
- `IsEnabled`: Yes or No
  - This returns for whether the account is **disabled**.
- `IsLocked`: Yes or No
  - This returns for whether the account is **locked out**.
  - (incorrect password attempts, temporary)
- `IsAdmin`: Yes or No
  - Does the account have admin access?
- `PasswordNeverExpires`: Yes or No
  - Password does not expire.
- `NoPasswordChange`: Yes or No
  - Password cannot be changed.
- `PasswordAge`: Number of Days (e.g. "7")
  - How old is the password?
- `LastLogonTime`: Date in format: Monday, January 02, 2006 3:04:05 PM
  - If the date is not in that exact format, the check will fail.
- `BadPasswordCount`: Number of incorrect passwords (e.g. "3")
  - How many incorrect passwords have been entered?
- `NumberOfLogons`: Number of total logons (e.g. "4")
  - How many times has the user logged on?
