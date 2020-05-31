package main

func writeDesktopFiles(mc *metaConfig) {
	firefoxBinary := `C:\Program Files (x86)\Mozilla Firefox\firefox.exe`
	if mc.Cli.Bool("v") {
		infoPrint("Writing ScoringReport.html shortcut to Desktop...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\web\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Writing ReadMe.html shortcut to Desktop...")
	}
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ReadMe.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\web\ReadMe.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Creating or emptying TeamID.txt file...")
	}
	cmdString = "echo 'YOUR-TEAMID-HERE' > C:\\aeacus\\misc\\TeamID.txt"
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Writing TeamID shortcut to Desktop...")
	}
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\TeamID.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Arguments = "-i yes"; $Shortcut.Save()`
	shellCommand(cmdString)

	// todo configure autologin user (netplwiz?)
	// domain compatibility? doubt
}

func installService(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Installing service with sc.exe...")
	}
	cmdString := `sc.exe create CSSClient binPath= "C:\aeacus\phocus.exe" start= "auto" DisplayName= "CSSClient"`
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Setting service description...")
	}
	cmdString = `sc.exe description CSSClient "This is Aeacus's Competition Scoring System client. Don't stop or mess with this unless you want to not get points, and maybe have your registry deleted."`
	shellCommand(cmdString)
}

func cleanUp(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Removing scoring.conf and ReadMe.conf...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\scoring.conf")
	shellCommand("Remove-Item -Force C:\\aeacus\\ReadMe.conf")
	if mc.Cli.Bool("v") {
		infoPrint("Removing previous.txt...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\misc\\previous.txt")
	if mc.Cli.Bool("v") {
		infoPrint("Emptying recycle bin...")
	}
	shellCommand("Clear-RecycleBin -Force")
	if mc.Cli.Bool("v") {
		infoPrint("Clearing recently used...")
	}
	shellCommand("Remove-Item -Force '${env:USERPROFILE}\\AppData\\Roaming\\Microsoft\\Windows\\Recent‌​*.lnk'")
	if mc.Cli.Bool("v") {
		warnPrint("Done cleaning up! You need to remove aeacus.exe manually.")
	}
}
