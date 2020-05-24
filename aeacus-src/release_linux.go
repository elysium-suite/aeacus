package main

func writeDesktopFiles(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Creating or emptying TeamID.txt...")
	}
	shellCommand("echo 'YOUR-TEAMID-HERE' > /opt/aeacus/misc/TeamID.txt")
	if mc.Cli.Bool("v") {
		infoPrint("Writing shortcuts to Desktop...")
	}
	shellCommand("cp " + mc.DirPath + "misc/*.desktop /home/" + mc.Config.User + "/Desktop/")
	shellCommand("chmod +x /home/" + mc.Config.User + "/Desktop/*.desktop")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " /home/" + mc.Config.User + "/Desktop/*")

	// todo configure autologin user
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
		infoPrint("Changing perms to 755 in /opt/aeacus...")
	}
	shellCommand("chmod 755 -R /opt/aeacus")

	if mc.Cli.Bool("v") {
		infoPrint("Removing .viminfo files...")
	}
	shellCommand("find / -name '.viminfo' -delete")

	if mc.Cli.Bool("v") {
		infoPrint("Symlinking .bash_history and .zsh_history to /dev/null...")
	}
	shellCommand("find / -name '.bash_history' -exec ln -sf /dev/null {} \\;")
	shellCommand("find / -name '.zsh_history' -exec ln -sf /dev/null {} \\;")

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
	shellCommand("find / -type f -iname '*.swp' -delete")
	shellCommand("rm -rf /root/*~ /home/*/Desktop/*~")

	if mc.Cli.Bool("v") {
		infoPrint("Removing crash and VMWare data...")
	}
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	if mc.Cli.Bool("v") {
		infoPrint("Removing apt and dpkg logs...")
	}
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	if mc.Cli.Bool("v") {
		infoPrint("Removing logs (auth and syslog)")
	}
	shellCommand("rm -f /var/log/auth.log* /var/log/syslog*")

	if mc.Cli.Bool("v") {
		infoPrint("Removing initial package list")
	}
	shellCommand("rm -f /var/log/installer/initial-status.gz")

	if mc.Cli.Bool("v") {
		infoPrint("Removing scoring.conf...")
	}
	shellCommand("rm /opt/aeacus/scoring.conf*")

	if mc.Cli.Bool("v") {
		infoPrint("Removing other setup files...")
	}
	shellCommand("rm -rf /opt/aeacus/misc/LICENSE /opt/aeacus/misc/ReadMe* /opt/aeacus/misc/ScoringReport* /opt/aeacus/misc/TeamID.desktop /opt/aeacus/misc aeacus-client /opt/aeacus/misc/*.sh /opt/aeacus/misc/previous.txt /opt/aeacus/ReadMe.conf /opt/aeacus/README.md /opt/aeacus/TODO.md")

	if mc.Cli.Bool("v") {
		infoPrint("Removing aeacus binary...")
	}
	shellCommand("rm /opt/aeacus/aeacus")

	if mc.Cli.Bool("v") {
		infoPrint("Overwriting timestamps to obfuscate changes...")
	}
	shellCommand("find /etc -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
	shellCommand("find /home -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
	shellCommand("find /var -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
	shellCommand("find /opt -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
}
