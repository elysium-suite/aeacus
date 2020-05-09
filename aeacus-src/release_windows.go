package main

func writeDesktopFiles(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Writing shortcuts to Desktop...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "C:\Program Files\Mozilla Firefox\firefox.exe C:\aeacus\web\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)

	if mc.Cli.Bool("v") {
		infoPrint("Creating TeamID.txt file...")
	}
	shellCommand(`echo 'YOUR-TEAMID-HERE' > C:\Users\` + mc.Config.User + `\Desktop\TeamID.txt`)
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
    shellCommand("Remove-Item C:\\aeacus\\scoring.conf")
	if mc.Cli.Bool("v") {
		infoPrint("Removing aeacus.exe...")
	}
    shellCommand("Remove-Item C:\\aeacus\\aeacus.exe")
	warnPrint("oops cleanup doesnt do too much yet")
	warnPrint("just empty trash bin? and recently used? idk look around")
}
