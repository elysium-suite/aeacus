package main

// writeDesktopFiles creates TeamID.txt and its shortcut, as well as links
// to the ScoringReport, ReadMe, and other needed files.
func writeDesktopFiles() {
	info("Creating or emptying TeamID.txt...")
	shellCommand("echo 'YOUR-TEAMID-HERE' > " + dirPath + "TeamID.txt")
	shellCommand("chmod 666 " + dirPath + "TeamID.txt")
	shellCommand("chown " + conf.User + ":" + conf.User + " " + dirPath + "TeamID.txt")
	info("Writing shortcuts to Desktop...")
	shellCommand("mkdir -p /home/" + conf.User + "/Desktop/")
	shellCommand("cp " + dirPath + "misc/desktop/*.desktop /home/" + conf.User + "/Desktop/")
	shellCommand("chmod +x /home/" + conf.User + "/Desktop/*.desktop")
	shellCommand("chown " + conf.User + ":" + conf.User + " /home/" + conf.User + "/Desktop/*")
}

// configureAutologin configures the auto-login capability for LightDM and
// GDM3, so that the image automatically logs in to the main user's account
// on boot.
func configureAutologin() {
	lightdm, _ := cond{Path: "/usr/share/lightdm"}.PathExists()
	gdm, _ := cond{Path: "/etc/gdm3/"}.PathExists()
	if lightdm {
		info("LightDM detected for autologin.")
		shellCommand(`echo "autologin-user=` + conf.User + `" >> /usr/share/lightdm/lightdm.conf.d/50-ubuntu.conf`)
	} else if gdm {
		info("GDM3 detected for autologin.")
		shellCommand(`echo -e "AutomaticLoginEnable=True\nAutomaticLogin=` + conf.User + `" >> /etc/gdm3/daemon.conf`)
	} else {
		fail("Unable to configure autologin! Please do so manually.")
	}
}

// installFont is skipped for Linux.
func installFont() {
	info("Skipping font install for Linux...")
}

// installService for Linux installs and starts the CSSClient init.d service.
func installService() {
	info("Installing service...")
	shellCommand("cp " + dirPath + "misc/dev/CSSClient /etc/init.d/")
	shellCommand("chmod +x /etc/init.d/CSSClient")
	shellCommand("systemctl enable CSSClient")
	shellCommand("systemctl start CSSClient")
}

// cleanUp for Linux is primarily focused on removing cached files, history,
// and other pieces of forensic evidence. It also removes the non-required
// files in the aeacus directory.
func cleanUp() {
	findPaths := "/bin /etc /home /opt /root /sbin /srv /usr /mnt /var"

	info("Changing perms to 755 in " + dirPath + "...")
	shellCommand("chmod 755 -R " + dirPath)

	info("Removing aeacus binary...")
	shellCommand("rm " + dirPath + "aeacus")

	info("Removing scoring.conf...")
	shellCommand("rm " + dirPath + "scoring.conf*")

	info("Removing other setup files...")
	shellCommand("rm -rf " + dirPath + "misc/")
	shellCommand("rm -rf " + dirPath + "ReadMe.conf")
	shellCommand("rm -rf " + dirPath + "README.md")
	shellCommand("rm -rf " + dirPath + ".git")
	shellCommand("rm -rf " + dirPath + ".github")
	shellCommand("rm -rf " + dirPath + "*.go")
	shellCommand("rm -rf " + dirPath + "Makefile")
	shellCommand("rm -rf " + dirPath + "go.*")
	shellCommand("rm -rf " + dirPath + "*.exe")
	shellCommand("rm -rf " + dirPath + "docs")

	if !ask("Do you want to remove cache and log files, overwrite timestamps, and remove other forensic data from this machine? This may impact data used for your forensic questions!") {
		return
	}

	info("Removing .viminfo and .swp files...")
	shellCommand("find " + findPaths + " -iname '*.viminfo*' -delete -iname '*.swp' -delete")

	info("Symlinking .bash_history and .zsh_history to /dev/null...")
	shellCommand(`find ` + findPaths + ` -iname '*.bash_history' -exec ln -sf /dev/null {} \;`)
	shellCommand(`find ` + findPaths + ` -name '.zsh_history' -exec ln -sf /dev/null {} \;`)

	info("Removing .mysql_history...")
	shellCommand(`find ` + findPaths + ` -name '.mysql_history' -exec rm {} \;`)

	info("Removing .local files...")
	shellCommand("rm -rf /root/.local /home/*/.local/")

	info("Removing cache...")
	shellCommand("rm -rf /root/.cache /home/*/.cache/")

	info("Removing temp root and Desktop files...")
	shellCommand("rm -rf /root/*~ /home/*/Desktop/*~")

	info("Removing crash and VMWare data...")
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	info("Removing apt and dpkg logs...")
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	info("Removing logs (auth and syslog)...")
	shellCommand("rm -f /var/log/auth.log* /var/log/syslog*")

	info("Removing initial package list...")
	shellCommand("rm -f /var/log/installer/initial-status.gz")

	info("Installing BleachBit...")
	shellCommand("apt-get install -y bleachbit")

	info("Clearing Firefox cache and browsing history...")
	shellCommand("bleachbit --clean firefox.url_history; bleachbit --clean firefox.cache")

	info("Overwriting timestamps to obfuscate changes...")
	shellCommand(`find /etc /home /var -exec touch --date='2012-12-12 12:12' {} \; 2>/dev/null`)
}
