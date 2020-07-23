#!/bin/bash

FILE="/opt/aeacus/TeamID.txt"

teamid=$(zenity --entry= \
		--text="Enter in your TeamID here"
		)
if [[ ${#teamid} > 0 ]]; then
	echo $teamid > $FILE
else
	notify-send "No ID Specified"
fi