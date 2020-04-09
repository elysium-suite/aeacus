package main

import (
	"fmt"
)

func writeDesktopFiles(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Writing shortcuts to Desktop...")
	}
    // scoringreport: file:///C:/aeacus/web/ScoringReport.html

	fmt.Println("xxd")

	if mc.Cli.Bool("v") {
		infoPrint("Creating TeamID.txt file...")
	}

	fmt.Println("xxd")
}

func installService(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		infoPrint("Throwing shortcut into the startup folder...")
	}
    cmdString := `$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut("C:\Users\"` + mc.Config.User + `"\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\aeacus-client.lnk"); $Shortcut.TargetPath = "C:\aeacus\phocus.exe"; $Shortcut.Save()`
	cmd = exec.Command("powershell.exe", "-c", cmdString)
	cmd.Run()
}

func cleanUp(mc *metaConfig) {
    warnPrint("oops cleanup doesnt do anything yet")

}
