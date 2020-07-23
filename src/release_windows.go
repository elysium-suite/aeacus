package main

func writeDesktopFiles(mc *metaConfig) {
	firefoxBinary := `C:\Program Files (x86)\Mozilla Firefox\firefox.exe`
	if verboseEnabled {
		infoPrint("Writing ScoringReport.html shortcut to Desktop...")
	}
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if verboseEnabled {
		infoPrint("Writing ReadMe.html shortcut to Desktop...")
	}
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ReadMe.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ReadMe.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	if verboseEnabled {
		infoPrint("Creating or emptying TeamID.txt file...")
	}
	cmdString = "echo 'YOUR-TEAMID-HERE' > C:\\aeacus\\TeamID.txt"
	shellCommand(cmdString)
	if verboseEnabled {
		infoPrint("Writing TeamID shortcut to Desktop...")
	}
	powershellPermission := `
	$ACL = Get-ACL C:\aeacus\TeamID.txt
	$ACL.SetOwner([System.Security.Principal.NTAccount] $env:USERNAME)
	Set-Acl -Path C:\aeacus\TeamID.txt -AclObject $ACL
	`
	shellCommand(powershellPermission)
	if verboseEnabled {
		infoPrint("Changing Permissions of TeamID")
	}

	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\TeamID.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Arguments = "-i yes"; $Shortcut.Save()`
	shellCommand(cmdString)

	// domain compatibility? doubt
}

func configureAutologin(mc *metaConfig) {
	if verboseEnabled {
		infoPrint("Setting Up autologin for " + mc.Config.User + "...")
	}
	powershellAutoLogin := `
	$RegPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\DefaultUsername"
	If(!(Test-Path $RegPath)) {
		New-ItemProperty HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\ -name "DefaultUsername" -Value "$env:USERNAME" -type String
	} else {
		Set-ItemProperty HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\ -name "DefaultUsername" -Value "$env:USERNAME" -type String
	}
	`
	shellCommand(powershellAutoLogin)
}

func installService() {
	if verboseEnabled {
		infoPrint("Installing service with sc.exe...")
	}
	cmdString := `sc.exe create CSSClient binPath= "C:\aeacus\phocus.exe" start= "auto" DisplayName= "CSSClient"`
	shellCommand(cmdString)
	if verboseEnabled {
		infoPrint("Setting service description...")
	}
	cmdString = `sc.exe description CSSClient "This is Aeacus's Competition Scoring System client. Don't stop or mess with this unless you want to not get points, and maybe have your registry deleted."`
	shellCommand(cmdString)
}

func cleanUp() {
	if verboseEnabled {
		infoPrint("Removing scoring.conf and ReadMe.conf...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\scoring.conf")
	shellCommand("Remove-Item -Force C:\\aeacus\\ReadMe.conf")
	if verboseEnabled {
		infoPrint("Removing previous.txt...")
	}
	shellCommand("Remove-Item -Force C:\\aeacus\\previous.txt")
	if verboseEnabled {
		infoPrint("Emptying recycle bin...")
	}
	shellCommand("Clear-RecycleBin -Force")
	if verboseEnabled {
		infoPrint("Clearing recently used...")
	}
	shellCommand("Remove-Item -Force '${env:USERPROFILE}\\AppData\\Roaming\\Microsoft\\Windows\\Recent‌​*.lnk'")
	if verboseEnabled {
		warnPrint("Done cleaning up! You need to remove aeacus.exe manually. The only things you need in the C:\\aeacus directory is phocus, scoring.dat, TeamID.txt, and the assets directory.")
	}
}
