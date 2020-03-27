package main

import (
	"fmt"
    "os/exec"
)


func writeDesktopFilesL(mc *metaConfig) {

    if mc.Cli.Bool("v") {
    	infoPrint("Writing shortcuts to Desktop...")
    }

    cmd := exec.Command("sh", "-c", "ln -sf " + mc.DirPath + "web/ReadMe.html /home/"+ mc.Config.User + "/Desktop/ReadMe")
    cmd.Run()
    cmd = exec.Command("sh", "-c", "ln -sf " + mc.DirPath + "web/ScoringReport.html /home/" + mc.Config.User + "/Desktop/ScoringReport")
    cmd.Run()

    if mc.Cli.Bool("v") {
    	infoPrint("Creating TeamID.txt file...")
    }

    cmd = exec.Command("sh", "-c", "ln -sf " + mc.DirPath + "web/ReadMe.html /home/"+ mc.Config.User + "/Desktop/ReadMe")
    cmd.Run()
    cmd = exec.Command("sh", "-c", "ln -sf " + mc.DirPath + "web/ScoringReport.html /home/" + mc.Config.User + "/Desktop/ScoringReport")
    cmd.Run()
}

func installServiceL(mc *metaConfig) {
    if mc.Cli.Bool("v") {
    	infoPrint("Installing service...")
    }
    cmd := exec.Command("sh", "-c", "echo '* * * * * root /opt/aeacus/phocus' >> /etc/crontab")
    cmd.Run()
    fmt.Println("Not really sure how to do that... atm doing cronjob. it works tm?")
}

func cleanUpL(mc *metaConfig) {

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
