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

For these checks, you can specify comparison through less, greater, and equal to through adding characters to the beginning of the value field.
This character **must** be added as the first character of the value field, as shown below.
>`<[value]`: The property is less than the `value` field.
>
>`>[value]`: The property is greater than the `value` field.
>
>`[value]`: The property is equal to the `value` field. This is the default.
- `PasswordAge`: Number of Days (e.g. "7")
  - How old is the password?
- `BadPasswordCount`: Number of incorrect passwords (e.g. "3")
  - How many incorrect passwords have been entered?
- `NumberOfLogons`: Number of total logons (e.g. "4")
  - How many times has the user logged on?

For the `LastLogonTime` property:
> `<[value]`: The property is before the `value` field's date.
>
> `>[value]`: The property is after the `value` field's date.
>
> `[value]`: The property is equal to the `value` field's date. This is the default.
- `LastLogonTime`: Date in format: Monday, January 02, 2006 3:04:05 PM
    - If the date is not in that exact format, the check will fail.
    - Time **must** be in UTC to pass.
