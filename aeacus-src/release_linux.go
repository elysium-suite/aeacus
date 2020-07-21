package main

func writeDesktopFiles(mc *metaConfig) {
	if verboseEnabled {
		infoPrint("Creating or emptying TeamID.txt...")
	}
	shellCommand("touch /opt/aeacus/TeamID.txt")
	shellCommand("echo 'YOUR-TEAMID-HERE' > /opt/aeacus/TeamID.txt")
	if verboseEnabled {
		infoPrint("Writing shortcuts to Desktop...")
	}
	shellCommand("cp " + mc.DirPath + "misc/*.desktop /home/" + mc.Config.User + "/Desktop/")
	shellCommand("chmod +x /home/" + mc.Config.User + "/Desktop/*.desktop")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " /home/" + mc.Config.User + "/Desktop/*")
}

/*
func configureAutologin(mc *metaConfig) {
	// set up user autologin
}
*/

func installService() {
	if verboseEnabled {
		infoPrint("Installing service...")
	}
	shellCommand("cp /opt/aeacus/misc/aeacus-client /etc/init.d/")
	shellCommand("chmod +x /etc/init.d/aeacus-client")
	shellCommand("systemctl enable aeacus-client")
	shellCommand("systemctl start aeacus-client")
}

func cleanUp() {

	findPaths := "/bin /etc /home /opt /root /sbin /srv /usr /mnt /var"
	if verboseEnabled {
		infoPrint("Changing perms to 755 in /opt/aeacus...")
	}
	shellCommand("chmod 755 -R /opt/aeacus")

	if verboseEnabled {
		infoPrint("Removing .viminfo and .swp files...")
	}
	shellCommand("find " + findPaths + " -iname '*.viminfo*' -delete -iname '*.swp' -delete")

	if verboseEnabled {
		infoPrint("Symlinking .bash_history and .zsh_history to /dev/null...")
	}
	shellCommand("find " + findPaths + " -iname '*.bash_history' -exec ln -sf /dev/null {} \\;")
	shellCommand("find " + findPaths + " -name '.zsh_history' -exec ln -sf /dev/null {} \\;")

	if verboseEnabled {
		infoPrint("Removing .local files...")
	}
	shellCommand("rm -rf /root/.local /home/*/.local/")

	if verboseEnabled {
		infoPrint("Removing cache...")
	}
	shellCommand("rm -rf /root/.cache /home/*/.cache/")

	if verboseEnabled {
		infoPrint("Removing temp root and Desktop files...")
	}
	shellCommand("rm -rf /root/*~ /home/*/Desktop/*~")

	if verboseEnabled {
		infoPrint("Removing crash and VMWare data...")
	}
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	if verboseEnabled {
		infoPrint("Removing apt and dpkg logs...")
	}
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	if verboseEnabled {
		infoPrint("Removing logs (auth and syslog)...")
	}
	shellCommand("rm -f /var/log/auth.log* /var/log/syslog*")

	if verboseEnabled {
		infoPrint("Removing initial package list...")
	}
	shellCommand("rm -f /var/log/installer/initial-status.gz")

	if verboseEnabled {
		infoPrint("Removing scoring.conf...")
	}
	shellCommand("rm /opt/aeacus/scoring.conf*")

	if verboseEnabled {
		infoPrint("Removing other setup files...")
	}
	shellCommand("rm -rf /opt/aeacus/misc/ /opt/aeacus/ReadMe.conf /opt/aeacus/README.md /opt/aeacus/TODO.md")

	if verboseEnabled {
		infoPrint("Removing aeacus binary...")
	}
	shellCommand("rm /opt/aeacus/aeacus")

	if verboseEnabled {
		infoPrint("Overwriting timestamps to obfuscate changes...")
	}
	shellCommand("find /etc /home /var -exec  touch --date='2012-12-12 12:12' {} \\; 2>/dev/null")
}
