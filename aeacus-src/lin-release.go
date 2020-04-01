package main

import (
	"os/exec"
)

func writeDesktopFilesL(mc *metaConfig) {

	if mc.Cli.Bool("v") {
		infoPrint("Writing shortcuts to Desktop...")
	}

	cmd := exec.Command("sh", "-c", "cp /opt/aeacus/desktop/*.desktop /home/"+mc.Config.User+"/Desktop/")
	cmd.Run()
	cmd = exec.Command("sh", "-c", "chmod +x /home/"+mc.Config.User+"/Desktop/*.desktop")
	cmd.Run()
}

func installServiceL(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Installing service...")
	}
	cmd := exec.Command("sh", "-c", "cp /opt/aeacus/misc/aeacus-client /etc/init.d/")
	cmd.Run()
	cmd = exec.Command("sh", "-c", "chmod +x /etc/init.d/aeacus-client")
	cmd.Run()
	cmd = exec.Command("sh", "-c", "systemctl enable aeacus-client")
	cmd.Run()
}

func cleanUpL(mc *metaConfig) {

	if mc.Cli.Bool("v") {
		infoPrint("Changing perms to 644 in /opt/aeacus...")
	}
	cmd := exec.Command("sh", "-c", "chmod 644 -R /opt/aeacus")
	cmd.Run()

	if mc.Cli.Bool("v") {
		infoPrint("Removing .viminfo files...")
	}
	cmd = exec.Command("sh", "-c", "find / -name \".viminfo\" -exec rm {} \\;")
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
		infoPrint("Removing .swp files")
	}
	cmd = exec.Command("sh", "-c", "rm -rf find / -type f -iname '*.swp' -delete")
    cmd.Run()

    if mc.Cli.Bool("v") {
		infoPrint("Removing .swp files")
	}
	cmd = exec.Command("sh", "-c", "rm -rf /home/*/.local/")
    cmd.Run()

    if mc.Cli.Bool("v") {
		infoPrint("Removing cache...")
	}
	cmd = exec.Command("sh", "-c", "rm -rf /home/*/.cache/")
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
		infoPrint("Removing other setup files...")
	}
	cmd = exec.Command("sh", "-c", "rm -rf /opt/aeacus/misc ReadMe.conf README.md TODO.md")
	cmd.Run()

	if mc.Cli.Bool("v") {
		infoPrint("Removing aeacus binary...")
	}
	cmd = exec.Command("sh", "-c", "rm /opt/aeacus/aeacus")
	cmd.Run()
}
