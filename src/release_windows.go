package main

func writeDesktopFiles() {
	firefoxBinary := `C:\Program Files (x86)\Mozilla Firefox\firefox.exe`
	infoPrint("Writing ScoringReport.html shortcut to Desktop...")
	cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ScoringReport.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ScoringReport.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	infoPrint("Writing ReadMe.html shortcut to Desktop...")
	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\ReadMe.lnk"); $Shortcut.TargetPath = "` + firefoxBinary + `"; $Shortcut.Arguments = "C:\aeacus\assets\ReadMe.html"; $Shortcut.Save()`
	shellCommand(cmdString)
	infoPrint("Creating or emptying TeamID.txt file...")
	cmdString = "echo 'YOUR-TEAMID-HERE' > C:\\aeacus\\TeamID.txt"
	shellCommand(cmdString)
	infoPrint("Writing TeamID shortcut to Desktop...")
	powershellPermission := `
	$ACL = Get-ACL C:\aeacus\TeamID.txt
	$ACL.SetOwner([System.Security.Principal.NTAccount] $env:USERNAME)
	Set-Acl -Path C:\aeacus\TeamID.txt -AclObject $ACL
	`
	shellCommand(powershellPermission)
	infoPrint("Changing Permissions of TeamID")

	cmdString = `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\` + mc.Config.User + `\Desktop\TeamID.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Arguments = "-i yes"; $Shortcut.Save()`
	shellCommand(cmdString)

	// domain compatibility? doubt
}

func configureAutologin() {
	infoPrint("Setting Up autologin for " + mc.Config.User + "...")
	powershellAutoLogin := `
	function Test-RegistryValue {

		param (

		 [parameter(Mandatory=$true)]
		 [ValidateNotNullOrEmpty()]$Path,

		[parameter(Mandatory=$true)]
		 [ValidateNotNullOrEmpty()]$Value
		)

		try {

		Get-ItemProperty -Path $Path | Select-Object -ExpandProperty $Value -ErrorAction Stop | Out-Null
		 return $true
		 }

		catch {

		return $false

		}

	}
	$RegPath1Exists = Test-RegistryValue -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Value "DefaultUsername"
	if ($RegPath1Exists -eq $false) {
		New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "DefaultUsername" -Value $env:USERNAME -type String
	}
	elseif ($RegPath1Exists -eq $true) {
		Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "DefaultUsername" -Value $env:USERNAME -type String
	}

	$RegPath2Exists = Test-RegistryValue -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Value "AutoAdminLogon"
	if ($RegPath2Exists -eq $false) {
		New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "AutoAdminLogon" -Value 1 -type String
	}
	elseif ($RegPath2Exists -eq $true) {
		Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -name "AutoAdminLogon" -Value 1 -type String
	}
	`
	shellCommand(powershellAutoLogin)
}

func installFont() {
	infoPrint("Installing Raleway font for Winform...")
	powershellFontInstall := `
	$SourceDir   = "C:\aeacus\assets\Raleway"
	$Source      = "C:\aeacus\assets\Raleway\*"
	$Destination = (New-Object -ComObject Shell.Application).Namespace(0x14)
	$TempFolder  = "C:\Windows\Temp\Fonts"

	# Create the source directory if it doesn't already exist
	New-Item -ItemType Directory -Force -Path $SourceDir | Out-Null

	New-Item $TempFolder -Type Directory -Force | Out-Null

	Get-ChildItem -Path $Source -Include '*.ttf','*.ttc','*.otf' -Recurse | ForEach {
		If (-not(Test-Path "C:\Windows\Fonts\$($_.Name)")) {

			$Font = "$TempFolder\$($_.Name)"

			# Copy font to local temporary folder
			Copy-Item $($_.FullName) -Destination $TempFolder

			# Install font
			$Destination.CopyHere($Font,0x10)

			# Delete temporary copy of font
			Remove-Item $Font -Force
		}
	}
	`
	shellCommand(powershellFontInstall)
}

func installService() {
	infoPrint("Installing service with sc.exe...")
	cmdString := `sc.exe create CSSClient binPath= "C:\aeacus\phocus.exe" start= "auto" DisplayName= "CSSClient"`
	shellCommand(cmdString)
	infoPrint("Setting service description...")
	cmdString = `sc.exe description CSSClient "This is Aeacus's Competition Scoring System client. Don't stop or mess with this unless you want to not get points, and maybe have your registry deleted."`
	shellCommand(cmdString)
}

func cleanUp() {
	infoPrint("Removing scoring.conf and ReadMe.conf...")
	shellCommand("Remove-Item -Force C:\\aeacus\\scoring.conf")
	shellCommand("Remove-Item -Force C:\\aeacus\\ReadMe.conf")
	infoPrint("Removing previous.txt...")
	shellCommand("Remove-Item -Force C:\\aeacus\\previous.txt")
	infoPrint("Emptying recycle bin...")
	shellCommand("Clear-RecycleBin -Force")
	infoPrint("Clearing recently used...")
	shellCommand("Remove-Item -Force '${env:USERPROFILE}\\AppData\\Roaming\\Microsoft\\Windows\\Recent‌​*.lnk'")
	infoPrint("Clearing run.exe command history...")
	clearRunScript := `$path = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\RunMRU" 
	$arr = (Get-Item -Path $path).Property 
	foreach($item in $arr) 
	{ 
	   if($item -ne "MRUList") 
	   { 
		 Remove-ItemProperty -Path $path -Name $item -ErrorAction SilentlyContinue 
	   } 
	}`
	shellCommand(clearRunScript)
	infoPrint("Removing Command History for Powershell")
	shellCommand("Remove-Item (Get-PSReadlineOption).HistorySavePath")
	warnPrint("Done with automatic cleanup! You need to remove aeacus.exe manually. The only things you need in the C:\\aeacus directory is phocus, scoring.dat, TeamID.txt, and the assets directory.")
}
