#!/bin/bash

zenity --question \
        --text="Would you like to stop scoring for this image?" \
        --title="Aeacus SE"


if [[ $? -eq 0 ]]; then
    notify-send -i /opt/aeacus/web/assets/logo.png "Stopping scoring and shutting down"
    service aeacus-client stop
    rm -rf /opt/aeacus/phocus
    rm -rf /opt/aeacus/scoring.dat
    shutdown now
else
    notify-send -i /opt/aeacus/web/assets/logo.png "Confirmation failed!"
fi