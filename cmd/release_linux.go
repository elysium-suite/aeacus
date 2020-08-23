package cmd

// WriteDesktopFiles creates TeamID.txt and its shortcut, as well as links
// to the ScoringReport, ReadMe, and other needed files.
func WriteDesktopFiles() {
	infoPrint("Creating or emptying TeamID.txt...")
	shellCommand("echo 'YOUR-TEAMID-HERE' > " + mc.DirPath + "TeamID.txt")
	shellCommand("chmod 666 " + mc.DirPath + "TeamID.txt")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " " + mc.DirPath + "TeamID.txt")
	infoPrint("Writing shortcuts to Desktop...")
	shellCommand("cp " + mc.DirPath + "misc/desktop/*.desktop /home/" + mc.Config.User + "/Desktop/")
	shellCommand("chmod +x /home/" + mc.Config.User + "/Desktop/*.desktop")
	shellCommand("chown " + mc.Config.User + ":" + mc.Config.User + " /home/" + mc.Config.User + "/Desktop/*")
}

// ConfigureAutologin configures the auto-login capability for LightDM and
// GDM3, so that the image automatically logs in to the main user's account
// on boot.
func ConfigureAutologin() {
	lightdm, _ := pathExists("/usr/share/lightdm")
	gdm, _ := pathExists("/etc/gdm3/")
	if lightdm {
		infoPrint("LightDM detected for autologin.")
		shellCommand(`echo "autologin-user=` + mc.Config.User + `" >> /usr/share/lightdm/lightdm.conf.d/50-ubuntu.conf`)
	} else if gdm {
		infoPrint("GDM3 detected for autologin.")
		shellCommand(`echo -e "AutomaticLogin=True\nAutomaticLogin=` + mc.Config.User + `" >> /etc/gdm3/custom.conf`)
	} else {
		failPrint("Unable to configure autologin! Please do so manually.")
	}
}

// InstallFont is skipped for Linux Builds
func InstallFont() {
	infoPrint("Skipping font install for Linux...")
}

// InstallService for Linux installs and starts the CSSClient init.d service.
func InstallService() {
	infoPrint("Installing service...")
	shellCommand("cp " + mc.DirPath + "misc/dev/CSSClient /etc/init.d/")
	shellCommand("chmod +x /etc/init.d/CSSClient")
	shellCommand("systemctl enable CSSClient")
	shellCommand("systemctl start CSSClient")
}

// CleanUp for Linux is primarily focused on removing cached files, history,
// and other pieces of forensic evidence. It also removes the non-required
// files in the aeacus directory.
func CleanUp() {
	infoPrint("Installing BleachBit...")
	shellCommand("apt install -y bleachbit")

	findPaths := "/bin /etc /home /opt /root /sbin /srv /usr /mnt /var"

	infoPrint("Changing perms to 755 in " + mc.DirPath + "...")
	shellCommand("chmod 755 -R " + mc.DirPath)

	infoPrint("Removing .viminfo and .swp files...")
	shellCommand("find " + findPaths + " -iname '*.viminfo*' -delete -iname '*.swp' -delete")

	infoPrint("Symlinking .bash_history and .zsh_history to /dev/null...")
	shellCommand(`find " + findPaths + " -iname '*.bash_history' -exec ln -sf /dev/null {} \;`)
	shellCommand(`"find " + findPaths + " -name '.zsh_history' -exec ln -sf /dev/null {} \;`)

	infoPrint("Removing .local files...")
	shellCommand("rm -rf /root/.local /home/*/.local/")

	infoPrint("Removing cache...")
	shellCommand("rm -rf /root/.cache /home/*/.cache/")

	infoPrint("Removing temp root and Desktop files...")
	shellCommand("rm -rf /root/*~ /home/*/Desktop/*~")

	infoPrint("Removing crash and VMWare data...")
	shellCommand("rm -f /var/VMwareDnD/* /var/crash/*.crash")

	infoPrint("Removing apt and dpkg logs...")
	shellCommand("rm -rf /var/log/apt/* /var/log/dpkg.log")

	infoPrint("Removing logs (auth and syslog)...")
	shellCommand("rm -f /var/log/auth.log* /var/log/syslog*")

	infoPrint("Removing initial package list...")
	shellCommand("rm -f /var/log/installer/initial-status.gz")

	infoPrint("Removing scoring.conf...")
	shellCommand("rm " + mc.DirPath + "scoring.conf*")

	infoPrint("Removing other setup files...")
	shellCommand("rm -rf " + mc.DirPath + "misc/")
	shellCommand("rm -rf " + mc.DirPath + "ReadMe.conf")
	shellCommand("rm -rf " + mc.DirPath + "README.md")
	shellCommand("rm -rf " + mc.DirPath + "TODO.md")
	shellCommand("rm -rf " + mc.DirPath + ".git")
	shellCommand("rm -rf " + mc.DirPath + ".github")

	infoPrint("Removing aeacus binary...")
	shellCommand("rm " + mc.DirPath + "aeacus")

	infoPrint("Overwriting timestamps to obfuscate changes...")
	shellCommand(`find /etc /home /var -exec  touch --date='2012-12-12 12:12' {} \; 2>/dev/null`)

	infoPrint("Clearing firefox cache and browsing history...")
	shellCommand("bleachbit --clean firefox.url_history; bleachbit --clean firefox.cache")
}
