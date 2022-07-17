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

For these checks, the `modifier` field can be used to specify comparison.
Valid modifiers:
>`less`: The property is less than the `value` field.
>
>`greater`: The property is greater than the `value` field.
>
>`equal`: The property is equal to the `value` field.
- `PasswordAge`: Number of Days (e.g. "7")
  - How old is the password?
- `BadPasswordCount`: Number of incorrect passwords (e.g. "3")
  - How many incorrect passwords have been entered?
- `NumberOfLogons`: Number of total logons (e.g. "4")
  - How many times has the user logged on?

Valid modifiers for the `LastLogonTime` property:
> `before`: The property is before the `value` field's date.
>
> `after`: The property is after the `value` field's date.
- `LastLogonTime`: Date in format: Monday, January 02, 2006 3:04:05 PM
    - If the date is not in that exact format, the check will fail.
