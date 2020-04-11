package main

func writeDesktopFiles(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Writing shortcuts to Desktop...")
	}
	shellCommand("cp " + mc.DirPath + "misc/*.desktop /home/" + mc.Config.User + "/Desktop/")
	shellCommand("chmod +x /home/" + mc.Config.User + "/Desktop/*.desktop")
}

func installService(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Installing service...")
	}
	shellCommand("cp /opt/aeacus/misc/aeacus-client /etc/init.d/")
	shellCommand("chmod +x /etc/init.d/aeacus-client")
	shellCommand("systemctl enable aeacus-client")
}

func cleanUp(mc *metaConfig) {

	if mc.Cli.Bool("v") {
		infoPrint("Changing perms to 644 in /opt/aeacus...")
	}
	shellCommand("chmod 644 -R /opt/aeacus")

	if mc.Cli.Bool("v") {
		infoPrint("Removing .viminfo files...")
	}
	shellCommand("find / -name \".viminfo\" -delete")

	if mc.Cli.Bool("v") {
		infoPrint("Symlinking .bash_history to /dev/null...")
	}
	shellCommand("find / -name \".bash_history\" -exec ln -sf /dev/null {} \\;")

	if mc.Cli.Bool("v") {
		infoPrint("Removing .swp files")
	}
	shellCommand("find / -type f -iname '*.swp' -delete")

	if mc.Cli.Bool("v") {
		infoPrint("Removing .local files")
	}
	shellCommand("rm -rf /root/.local /home/*/.local/")

	if mc.Cli.Bool("v") {
		infoPrint("Removing cache...")
	}
	shellCommand("rm -rf /root/.cache /home/*/.cache/")

	if mc.Cli.Bool("v") {
		infoPrint("Removing swap and temp Desktop files")
	}
	shellCommand("rm -rf /home/*/Desktop/*~")

	if mc.Cli.Bool("v") {
		infoPrint("Removing crash and VMWare data...")
	}
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	if mc.Cli.Bool("v") {
		infoPrint("Removing apt and dpkg logs...")
	}
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	if mc.Cli.Bool("v") {
		infoPrint("Removing scoring.conf...")
	}
	shellCommand("rm /opt/aeacus/scoring.conf*")

	if mc.Cli.Bool("v") {
		infoPrint("Removing other setup files...")
	}
	shellCommand("rm -rf /opt/aeacus/misc /opt/aeacus/ReadMe.conf /opt/aeacus/README.md /opt/aeacus/TODO.md")

	if mc.Cli.Bool("v") {
		infoPrint("Removing aeacus binary...")
	}
	shellCommand("rm /opt/aeacus/aeacus")
}
