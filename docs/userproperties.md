# aeacus

## Windows User Properties

These are the user properties that you would normally find the the user edit dialog (`User cannot change password`, `Password never expires`, etc).

> Date format is "M/D/YYYY HH:MM:SS AM", but we only use the first part, for example, `5/23/2020` or `12/4/2021`. You really shouldn't be using dates though. You should probably be using `UserDetailNot` if you want to make sure a password expires (not Never).

> Values are cAsE sEnSiTiVe (`Yes` != `yes`).

- `Full name`: Full name of user (string)
- `Comment`: User account comment (string)
- `Account active`: Yes or No
- `Account expires`: Never or date
- `Password expires`: Never or date
- `Password required`: Yes or No
- `User may change password`: Yes or No
