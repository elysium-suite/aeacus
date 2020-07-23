#!/bin/bash

global=$(
	zenity --forms \
		--add-entry="Round Name" \
		--text="AeacusSE Round Information" \
		--add-entry="Round Title" \
		--add-entry="Main User" \
		--add-entry="OS" \
		--add-entry="Remote" \
		--add-entry="Local" \
		--add-entry="EndDate" \
		--add-entry="NoDestroy" \
		--separator=$'\034' \
		--title=AeacusSE
)

IFS=$'\034'
FILE="/opt/aeacus/scoring.conf"

# write_check(message, points, check.pass fields)
write_check() {
	echo "[[check]]" >>$FILE
	if [[ -n $1 ]]; then
		echo "message = \"$1\"" >>$FILE
	fi
	if [[ -n $2 ]]; then
		echo "points = $2" >>$FILE
	fi
	echo "[[check.pass]]" >>$FILE
	echo -e "$3" >>$FILE
	echo "" >>$FILE
}

# write_if_exist(key, value)
write_if_exist() {
	if [[ -n $2 ]]; then
		echo "$1 = \"$2\"" >>$FILE
	fi
}

read -r name title user os remote local end destroy <<<"$global"
write_if_exist "name" $name
write_if_exist "title" $title
write_if_exist "user" $user
write_if_exist "os" $os
write_if_exist "remote" $remote
write_if_exist "local" $local
write_if_exist "enddate" $end
write_if_exist "destroy" $destroy
echo >>$FILE

next="tmp"

while [[ "$next" != "None" ]] && [[ -n "$next" ]]; do
	next=$(
		zenity \
			--list \
			--title="Add vulnerabilities" \
			--text "Which category do you want to add?" \
			--radiolist \
			--column "Pick" \
			--column "Vuln Category" \
			FALSE "Command" \
			FALSE "CommandNot" \
			FALSE "FileExists" \
			FALSE "FileExistsNot" \
			FALSE "FileContains" \
			FALSE "FileContainsNot" \
			FALSE "FileContainsRegex" \
			FALSE "FileContainsRegexNot" \
			FALSE "DirContainsRegex" \
			FALSE "DirContainsRegexNot" \
			FALSE "FileEquals" \
			FALSE "FileEqualsNot" \
			FALSE "PackageInstalled" \
			FALSE "PackageInstalledNot" \
			FALSE "ServiceUp" \
			FALSE "ServiceUpNot" \
			FALSE "UserExists" \
			FALSE "UserExistsNot" \
			FALSE "FirewallUp" \
			FALSE "UserInGroup" \
			FALSE "UserInGroupNot" \
			FALSE "GuestDisabledLDM" \
			FALSE "PasswordChanged" \
			FALSE "None" \
			--height=700 \
			--width=600 \
			--title=AeacusSE
	)

	case $next in

	"Command")
		data=$(
			zenity --forms \
				--text="Pass if command succeeds" \
				--add-entry "Message" \
				--add-entry="Points" \
				--add-entry="Command" \
				--separator=$'\034'
		)
		read -r mess pts arg1 <<<"$data"
		vulns="type='Command'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"CommandNot")
		data=$(
			zenity --forms \
				--text="Pass if command doesn't succeed (exit code 0, checks \$?)" \
				--add-entry "Message" \
				--add-entry="Points" \
				--add-entry="Command" \
				--separator=$'\034'
		)
		read -r mess pts arg1 <<<"$data"
		vulns="type='CommandNot'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileExists")
		data=$(
			zenity --forms \
				--text="Pass if Specified File Exists" \
				--add-entry "Message" \
				--add-entry="Points" \
				--add-entry="File" \
				--separator=$'\034'
		)
		read -r mess pts arg1 <<<"$data"
		vulns="type='FileExists'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileExistsNot")
		data=$(zenity --forms --text="Pass if Specified File Does not Exist" --add-entry "Message" --add-entry="Points" --add-entry="File" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='FileExistsNot'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileContains")
		data=$(zenity --forms --text="Pass if File Contains String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileContains'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileContainsNot")
		data=$(zenity --forms --text="Pass if File Doesn't Contain String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileContainsNot'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileContainsRegex")
		data=$(zenity --forms --text="Pass if file Contains Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileContainsRegex'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileContainsRegexNot")
		data=$(zenity --forms --text="Pass if file doesn't contain Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileContainsRegexNot'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"DirContainsRegex")
		data=$(zenity --forms --text="Pass if Directory Contains Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='DirContainsRegex'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"DirContainsRegexNot")
		data=$(zenity --forms --text="Pass if Directory doesn't contain Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='DirContainsRegexNot'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileEquals")
		data=$(zenity --forms --text="Pass if file equals SHA1 Hash" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="hash" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileEquals'\narg1='${arg1}'\narg2='${arg2}'\n"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FileEqualsNot")
		data=$(zenity --forms --text="Pass if file equals SHA1 Hash" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="hash" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='FileEqualsNot'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"PackageInstalled")
		data=$(zenity --forms --text="Pass if Package is Installed" --add-entry "Message" --add-entry="Points" --add-entry="Package" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='PackageInstalled'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"PackageInstalledNot")
		data=$(zenity --forms --text="Pass if Package is Not Installed" --add-entry "Message" --add-entry="Points" --add-entry="Package" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='PackageInstalledNot'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"ServiceUp")
		data=$(zenity --forms --text="Pass if Service is Running" --add-entry "Message" --add-entry="Points" --add-entry="Service" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='ServiceUp'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"ServiceUpNot")
		data=$(zenity --forms --text="Pass if Service is not running" --add-entry "Message" --add-entry="Points" --add-entry="Service" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='ServiceUpNot'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"UserExists")
		data=$(zenity --forms --text="Pass if User Exists On System" --add-entry "Message" --add-entry="Points" --add-entry="User" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='UserExists'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"UserExistsNot")
		data=$(zenity --forms --text="Pass if User does not exist On System" --add-entry "Message" --add-entry="Points" --add-entry="User" --separator=$'\034')
		read -r mess pts arg1 <<<"$data"
		vulns="type='UserExistsNot'\narg1='${arg1}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"FirewallUp")
		data=$(zenity --forms --text="Pass if Firewall is Up, (Just Click Ok)" --separator=$'\034')
		vulns="type='FirewallUp'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"UserInGroup")
		data=$(zenity --forms --text="Pass if User is in Group" --add-entry "Message" --add-entry="Points" --add-entry="User" --add-entry="Group" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='UserInGroup'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"UserInGroupNot")
		data=$(zenity --forms --text="Pass if User is not in Group" --add-entry "Message" --add-entry="Points" --add-entry="User" --add-entry="Group" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='UserInGroupNot'\narg1='${arg1}'\narg2 = '${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"GuestDisabledLDM")
		data=$(zenity --forms --text="Pass if Guest is Disabled" --separator=$'\034')
		vulns="type='GuestDisabledLDM'"
		write_check "$mess" "$pts" "$vulns"
		;;

	"PasswordChanged")
		data=$(zenity --forms --text="Pass if user's password has changed" --add-entry "Message" --add-entry="Points" --add-entry="User" --add-entry="Hash" --separator=$'\034')
		read -r mess pts arg1 arg2 <<<"$data"
		vulns="type='PasswordChanged'\narg1='${arg1}'\narg2='${arg2}'"
		write_check "$mess" "$pts" "$vulns"
		;;
	esac
done
