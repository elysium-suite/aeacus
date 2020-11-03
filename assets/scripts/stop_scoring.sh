#!/usr/bin/env sh

if zenity --question \
	--text="Would you like to stop scoring for this image?" \
	--title="Aeacus SE"; then
	notify-send -i /opt/aeacus/assets/img/logo.png "Stopping scoring, and shutting down."
	service CSSClient stop
	pkill -9 phocus
	rm -f /opt/aeacus/phocus /opt/aeacus/scoring.dat
	shutdown now
else
	notify-send -i /opt/aeacus/assets/img/logo.png "Confirmation failed!"
fi
