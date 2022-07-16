# aeacus

## Windows User Properties

These are the user properties that you would normally find the the user edit dialog.

**Usernames are case sensitive.** `Yes`/`No` are not case sensitive. If you input something other than "yes" or "no" for a boolean check, it defaults to "no".

- `FullName`: Full name of user (case sensitive)
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
