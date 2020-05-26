package main

func writeDesktopFiles(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Writing ScoringReport.html shortcut to Desktop...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "C:\Program Files\Mozilla Firefox\firefox.exe C:\aeacus\web\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Creating or emptying TeamID.txt file...")
	}
	cmdString = "echo 'YOUR-TEAMID-HERE' > C:\\aeacus\\misc\\TeamID.txt"
	shellCommand(cmdString)
	if mc.Cli.Bool("v") {
		infoPrint("Writing TeamID.txt shortcut to Desktop...")
	}
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\TeamID.lnk"); $Shortcut.TargetPath = "C:\aeacus\misc\TeamID.txt"; $Shortcut.Save()`
	shellCommand(cmdString)

	// todo configure autologin user (netplwiz?)
	// domain compatibility? doubt
}

func installService(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Throwing shortcut into the startup folder...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\aeacus-client.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Save()`
	shellCommand(cmdString)
}

func cleanUp(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Removing scoring.conf...")
	}
	shellCommand("Remove-Item C:\\aeacus\\misc\\previous.txt")
	if mc.Cli.Bool("v") {
		infoPrint("Removing previous.txt...")
	}
	shellCommand("Remove-Item C:\\aeacus\\web.conf")
	if mc.Cli.Bool("v") {
		infoPrint("Removing aeacus.exe...")
	}
	shellCommand("Remove-Item C:\\aeacus\\aeacus.exe")
	if mc.Cli.Bool("v") {
		infoPrint("Emptying recycle bin...")
	}
	shellCommand("Clear-RecycleBin -Force")
	if mc.Cli.Bool("v") {
		infoPrint("Clearing recently used...")
	}
	shellCommand("Remove-Item -Force '${env:USERPROFILE}\\AppData\\Roaming\\Microsoft\\Windows\\Recent‌​*.lnk'")
}
