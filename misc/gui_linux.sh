#!/bin/bash
global=$(zenity --forms --add-entry="Round Name" --add-entry="Round Title" --add-entry="Main User" --add-entry="OS" --add-entry="Remote" --add-entry="Local" --add-entry="EndDate" --add-entry="NoDestroy" --separator=$'\034' --title=AeacusSE)

IFS=$'\034' read -r name title user os remote local end destroy <<< "$global"

next=""
vulns=""

while [[ "$next" != "None" ]]
do

next=$(zenity --list --title="Add vulnerabilities" --text "Which category do you want to add?" --radiolist --column "Pick" --column "Vuln Category" FALSE "Command" FALSE "CommandNot" FALSE "FileExists" FALSE "FileExistsNot" FALSE "FileContains" FALSE "FileContainsNot" FALSE "FileContainsRegex" FALSE "FileContainsRegexNot" FALSE "DirContainsRegex" FALSE "DirContainsRegexNot" FALSE "FileEquals" FALSE "FileEqualsNot" FALSE "PackageInstalled" FALSE "PackageInstalledNot" FALSE "ServiceUp" FALSE "ServiceUpNot" FALSE "UserExists" FALSE "UserExistsNot" FALSE "FirewallUp" FALSE "UserInGroup" FALSE "UserInGroupNot" FALSE "GuestDisabledLDM" FALSE "None" --height=700 --width=600 --title=AeacusSE)

case $next in

"Command")
	data=$(zenity --forms --text="Pass if command succeeds" --add-entry "Message" --add-entry="Points" --add-entry="Command" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"Command\"\narg1 = \"${arg1}\"\n"
	;;
"CommandNot")
	data=$(zenity --forms --text="Pass if command doesn't succeed (exit code 0, checks $?)" --add-entry "Message" --add-entry="Points" --add-entry="Command" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vvulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"CommandNot\"\narg1 = \"${arg1}\"\n"
	;;
"FileExists")
	data=$(zenity --forms --text="Pass if Specified File Exists" --add-entry "Message" --add-entry="Points" --add-entry="File" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileExists\"\narg1 = \"${arg1}\"\n"
	;;
"FileExistsNot")
	data=$(zenity --forms --text="Pass if Specified File Does not Exist" --add-entry "Message" --add-entry="Points" --add-entry="File" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileExistsNot\"\narg1 = \"${arg1}\"\n"
	;;
"FileContains")
	data=$(zenity --forms --text="Pass if File Contains String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileContains\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"FileContainsNot")
	data=$(zenity --forms --text="Pass if File Doesn't Contain String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileContainsNot\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"FileContainsRegex")
	data=$(zenity --forms --text="Pass if file Contains Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileContainsRegex\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"FileContainsRegexNot")
	data=$(zenity --forms --text="Pass if file doesn't contain Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileContainsRegexNot\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"DirContainsRegex")
	data=$(zenity --forms --text="Pass if Directory Contains Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"DirContainsRegex\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"DirContainsRegexNot")
	data=$(zenity --forms --text="Pass if Directory doesn't contain Regex String" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="String" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"DirContainsRegexNot\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"FileEquals")
	data=$(zenity --forms --text="Pass if file equals SHA1 Hash" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="hash" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileEquals\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;

"FileEqualsNot")
	data=$(zenity --forms --text="Pass if file equals SHA1 Hash" --add-entry "Message" --add-entry="Points" --add-entry="File" --add-entry="hash" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FileEqualsNot\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;

"PackageInstalled")
	data=$(zenity --forms --text="Pass if Package is Installed" --add-entry "Message" --add-entry="Points" --add-entry="Package" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"PackageInstalled\"\narg1 = \"${arg1}\"\n"
	;;

"PackageInstalledNot")
	data=$(zenity --forms --text="Pass if Package is Not Installed" --add-entry "Message" --add-entry="Points" --add-entry="Package" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"PackageInstalledNot\"\narg1 = \"${arg1}\"\n"
	;;

"ServiceUp")
	data=$(zenity --forms --text="Pass if Service is Running" --add-entry "Message" --add-entry="Points" --add-entry="Service" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"ServiceUp\"\narg1 = \"${arg1}\"\n"
	;;

"ServiceUpNot")
	data=$(zenity --forms --text="Pass if Service is not running" --add-entry "Message" --add-entry="Points" --add-entry="Service" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"ServiceUpNot\"\narg1 = \"${arg1}\"\n"
	;;


"UserExists")
	data=$(zenity --forms --text="Pass if User Exists On System" --add-entry "Message" --add-entry="Points" --add-entry="User" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"UserExists\"\narg1 = \"${arg1}\"\n"
	;;

"UserExistsNot")
	data=$(zenity --forms --text="Pass if User does not exist On System" --add-entry "Message" --add-entry="Points" --add-entry="User" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"UserExistsNot\"\narg1 = \"${arg1}\"\n"
	;;


"FireWallUp")
	data=$(zenity --forms --text="Pass if Firewall is Up, (Just Click Ok)" --separator=$'\034')
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"FireWallUp\"\n"
	;;
	
"UserInGroup")
	data=$(zenity --forms --text="Pass if User is in Group" --add-entry "Message" --add-entry="Points" --add-entry="User" --add-entry="Group" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"UserInGroup\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"UserInGroupNot")
	data=$(zenity --forms --text="Pass if User is not in Group" --add-entry "Message" --add-entry="Points" --add-entry="User" --add-entry="Group" --separator=$'\034')
	IFS=$'\034' read -r mess pts arg1 arg2 <<< "$data"
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"UserInGroupNot\"\narg1 = \"${arg1}\"\narg2 = \"${arg2}\"\n"
	;;
"GuestDisabledLDM")
	data=$(zenity --forms --text="Pass if Gues is Disabled" --separator=$'\034')
	vulns="${vulns}\n[[check]]\nmessage = \"${mess}\"\npoints = \"${pts}\"\n[[check.pass]]\ntype= \"GuestDisblaedLDM\"\n"
	;;
esac

done

FILE=/opt/aeacus/scoring.conf
if [ ! -f "$FILE" ]; then
    touch /opt/aeacus/scoring.conf
fi

echo "name = \"$name\"" >> /opt/aeacus/scoring.conf
echo "title = \"$title\"" >> /opt/aeacus/scoring.conf
echo "user = \"$user\"" >> /opt/aeacus/scoring.conf
echo "os = \"$os\"" >> /opt/aeacus/scoring.conf
echo "remote = \"$remote\"" >> /opt/aeacus/scoring.conf
echo "local = \"$local\"" >> /opt/aeacus/scoring.conf
echo "enddate = \"$end\"" >> /opt/aeacus/scoring.conf
echo "nodestory = \"$destroy\"" >> /opt/aeacus/scoring.conf

echo -e "$vulns" >> /opt/aeacus/scoring.conf
