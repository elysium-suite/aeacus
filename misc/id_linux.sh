#!/bin/bash

teamid=$(
	zenity --entry= \
		--text="Enter in your TeamID here"
)
if [[ ${#teamid} > 0 ]]; then
	echo $teamid >/opt/aeacus/TeamID.txt
else
	notify-send -i /opt/aeacus/assets/logo.png "Aeacus SE" "No ID specified!"
fi
