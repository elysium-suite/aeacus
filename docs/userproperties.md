# aeacus

## Windows User Properties

These are the user properties that you would normally find the the user edit dialog.

**Usernames are case sensitive.** `Yes`/`No` are not case sensitive. If you input something other than "yes" or "no" for a boolean check, it defaults to "no".

-   `FullName`: Full name of user (case sensitive)
-   `IsEnabled`: Yes or No 
	- This returns for whether or not the account is **disabled**.
-   `IsLocked`: Yes or No
	- This returns for whether or not the account is **locked out**.
-   `IsAdmin`: Yes or No
	- This returns for whether or not the account has administrative privileges.
-   `PasswordNeverExpires`: Yes or No
	- This returns for whether or not the password can never expire
-   `NoChangePassword` : Yes or No
	- This returns for whether or not the account can change their own password.
