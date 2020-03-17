package main

import (
	"fmt"
    "os/exec"
)

func cleanUp(mc *metaConfig) {

    if mc.Cli.Bool("v") {
    	infoPrint("Removing .viminfo files...")
    }
    cmd := exec.Command("sh", "-c", "find / -name \".viminfo\" -exec rm {} \\;")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing .bash_history...")
    }
    cmd = exec.Command("sh", "-c", "find / -name \".bash_history\" -exec rm {} \\;")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing recently-used...")
    }
    cmd = exec.Command("sh", "-c", "rm -rf /home/*/.local/share/recently-used.xbel")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing swap and temp Desktop files")
    }
    cmd = exec.Command("sh", "-c", "rm -rf /home/*/Desktop/*~")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing crash and VMWare data...")
    }
    cmd = exec.Command("sh", "-c", "rm -f /var/VMwareDnD/* /var/crash/*.crash")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing apt and dpkg logs...")
    }
    cmd = exec.Command("sh", "-c", "rm -rf /var/log/apt/* /var/log/dpkg.log")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing scoring.conf...")
    }
    cmd = exec.Command("sh", "-c", "rm /opt/aeacus/scoring.conf")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Removing aeacus binary...")
    }
    cmd = exec.Command("sh", "-c", "rm /opt/aeacus/aeacus")
    cmd.Run()

}

func writeDesktopFiles(mc *metaConfig) {
    if mc.Cli.Bool("v") {
    	infoPrint("Writing shortcut to ReadMe.html...")
    	infoPrint("Writing shortcut to ScoringReport.html...")
    	infoPrint("Creating TeamID.txt file...")
    }
}

func installService(mc *metaConfig) {
    if mc.Cli.Bool("v") {
    	infoPrint("Installing service...")
    }
}

func destroyImage() {
	// destroy the image if outside time range
	fmt.Println("destroying the system lol. todo")
}

func sendNotification(notifyTitle string, notifyBody string) {
	cmd := exec.Command("notify-send", "-i", "/opt/aeacus/web/assets/logo.png", notifyTitle, notifyBody)
    cmd.Run()
}
