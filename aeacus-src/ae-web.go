package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func genReport(mc *metaConfig, id *imageData, connStatus []string) {

	header := `<!DOCTYPE html> <html> <head> <meta http-equiv="refresh" content="10"> <title>Aeacus Scoring Report</title> <style type="text/css"> h1 { text-align: center; } h2 { text-align: center; } body { font-family: Arial, Verdana, sans-serif; font-size: 14px; margin: 0; padding: 0; width: 100%; height: 100%; background: url('assets/background.png'); background-size: cover; background-attachment: fixed; background-position: top center; background-color: #336699; } .red {color: red;} .green {color: green;} .blue {color: blue;} .main { margin-top: 10px; margin-bottom: 10px; margin-left: auto; margin-right: auto; padding: 0px; border-radius: 12px; background-color: white; width: 900px; max-width: 100%; min-width: 600px; box-shadow: 0px 0px 12px #003366; } .text { padding: 12px; -webkit-touch-callout: none; -webkit-user-select: none; -khtml-user-select: none; -moz-user-select: none; -ms-user-select: none; user-select: none; } .center { text-align: center; } .binary { position: relative; overflow: hidden; } .binary::before { position: absolute; top: -75%; left: -125%; display: block; width: 200%; height: 150%; -webkit-transform: rotate(-45deg); -moz-transform: rotate(-45deg); -ms-transform: rotate(-45deg); transform: rotate(-45deg); content: attr(data-binary); opacity: 0.16; line-height: 2em; letter-spacing: 2px; color: #369; font-size: 10px; pointer-events: none; } </style> <meta http-equiv="refresh"> </head> <body><div class="main"><div class="text"><div class="binary" data-binary="0000 0000 11010000 01100100"><p align=center style="width:100%;text-align:center"><img align=middle style="width:180px; float:middle" src="assets/logo.png"></p>`

	footer := `</p> <br> <p align=center style="text-align:center"> The Aeacus project is free and open source software. This project is in no way endorsed or affiliated with the Air Force Association or the University of Texas at San Antonio. </p> </div> </div> </div> </body> </html>`

	var htmlFile strings.Builder
	htmlFile.WriteString(header)
	genTime := time.Now()
	htmlFile.WriteString(fmt.Sprintf("<h1>%s</h1>", mc.Config.Title))
	htmlFile.WriteString(fmt.Sprintf("<h2>Report Generated At: %s </h2>", genTime.Format("2006/01/02 15:04:05 MST")))
	htmlFile.WriteString(`<script language="Javascript"> var bin = document.querySelectorAll('.binary'); [].forEach.call(bin, function(el) { el.dataset.binary = Array(4096).join(el.dataset.binary + ' ') }); var currentdate = new Date().getTime(); gendate = Date.parse('0000/00/00 00:00:00 UTC'); diff = Math.abs(currentdate - gendate); if ( gendate > 0 && diff > 1000 * 60 * 5 ) { document.write('<span style="color:red"><h2>WARNING: CCS Scoring service may not be running</h2></span>'); } </script>`)

	// Who needs timers, am I right
	//htmlFile.WriteString(`<h3 class="center">Approximate Image Running Time: 00:00:00</h3>`)
	//htmlFile.WriteString(`<h3 class="center">Approximate Team Running Time: 00:00:00</h3>`)

	if mc.Config.Remote != "" {
		htmlFile.WriteString(fmt.Sprintf(`<h3 class="center">Current Team ID: %s</h3>`, mc.TeamID))
	}

	htmlFile.WriteString(fmt.Sprintf(`<h2> %d out of %d points received</h2>`, id.Score, id.TotalPoints))

	if mc.Config.Remote != "" {
		htmlFile.WriteString(fmt.Sprintf(`<a href="http://%s/scores/css">Click here to view the public scoreboard</a><br>`, mc.Config.Remote))

		htmlFile.WriteString(fmt.Sprintf(`<p><h3>Connection Status: <span style="color:%s">%s<span></h3>`, connStatus[0], connStatus[1]))
		htmlFile.WriteString(fmt.Sprintf(`Internet Connectivity Check: <span style="color:%s">%s</span><br>`, connStatus[2], connStatus[3]))
		htmlFile.WriteString(fmt.Sprintf(`Aeacus Server Connection Status: <span style="color:%s">%s</span></p>`, connStatus[4], connStatus[5]))
	}

	htmlFile.WriteString(fmt.Sprintf(`<h3> %d penalties assessed, for a loss of %.0f points: </h3> <p> <span style="color:red">`, len(id.Penalties), math.Abs(float64(id.Detracts))))

	// for each penalty
	for _, penalty := range id.Penalties {
		htmlFile.WriteString(fmt.Sprintf("%s - %.0f pts<br>", penalty.Message, math.Abs(float64(penalty.Points))))
	}

	htmlFile.WriteString(fmt.Sprintf(`</span> </p> <h3> %d out of %d scored security issues fixed, for a gain of %d points:</h3><p>`, len(id.Points), id.ScoredVulns, id.Contribs))

	// for each point:
	for _, point := range id.Points {
		htmlFile.WriteString(fmt.Sprintf("%s - %d pts<br>", point.Message, point.Points))
	}

	htmlFile.WriteString(footer)

	if mc.Cli.Bool("v") {
		infoPrint("Writing HTML to ScoringReport.html...")
	}
	writeFile(mc.DirPath+"web/ScoringReport.html", htmlFile.String())
}

func genReadMe(mc *metaConfig) {

	header := `<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Aeacus README</title><link rel="stylesheet" type="text/css" href="./assets/bootstrap3.min.css"><link rel="stylesheet" type="text/css" href="./assets/bootstrap3-custom.min.css"><link rel="stylesheet" type="text/css" href="./assets/style.css"><style>body{background-image:url("./assets/background.png");background-size:cover}</style><div style="height: 100%; width: 100%"><div id="centerarea" class="container"><div id="MainRow" class="row"><div class="col-sm-12" style="min-height:725px;">s<div id="DeltaPlaceHolderMain"><div class="row"><div class="col-md-12"><div class="mod-wrap"><div class="article article-body"><div class="article-content"><div><div id="MSOZoneCell_WebPartctl00_ctl48_g_8c02beb0_d06b_46f1_9ba2_95bef95154c8" class="s4-wpcell-plain ms-webpartzone-cell ms-webpart-cell-vertical ms-fullWidth "><div class="ms-webpart-chrome ms-webpart-chrome-vertical ms-webpart-chrome-fullWidth "><style type="text/css">h1{text-align:center;font-family:Helvetica,Arial,sans-serif;font-size:36px;margin:10px;padding:30px 14px 10px 0px;width:100%;height:100%;color:#0D2E5B !important}h2{font-family:Helvetica,Arial,sans-serif;font-size:18px;margin:30px 0 10px 0;padding:0;width:100%;height:100%;color:#0D2E5B !important}body{font-family:Helvetica,Arial,sans-serif;font-size:16px;margin:0;padding:0;width:100%;height:100%;background-color:#0D2E5B}pre{font-family:Helvetica,Arial,sans-serif;font-size:16px}.main{margin-top:25px;margin-bottom:10px;margin-left:auto;margin-right:auto;padding:0px;background-color:white;max-width:100%}.text{padding-top:12px;padding-bottom:12px;padding-left:40px;padding-right:40px}.center{text-align:center}</style><div class="main"><div class="text"><p align="center"> <img src="./assets/logo.png" height="210" width="230"></p>`

	footer := `<h2>Competition Guidelines</h2><ul><li> In order to provide a better competition experience, you are <b>NOT</b> required to change the password of the primary, auto-login, user account. Changing the password of a user that is set to automatically log in may lock you out of your computer.</li><li> Authorized administrator passwords were correct the last time you did a password audit, but are not guaranteed to be currently accurate.</li><li> Do not stop or disable the Aeacus-Client service or process.</li><li> Do not remove any authorized users or their home directories.</li><li> The time zone of this image is set to UTC. Please do not change the time zone, date, or time on this image.</li><li> You can view your current scoring report by double-clicking the "Aeacus Scoring Report" desktop icon.</li><li> JavaScript is required for some error messages that appear on the "Aeacus Scoring Report." To ensure that you only receive correct error messages, please do not disable JavaScript.</li><li> Some security settings may prevent the Stop Scoring application from running. If this happens, the safest way to stop scoring is to suspend the virtual machine. You should <b>NOT</b> power on the VM again before deleting.</li></ul><p align="center" style="text-align:center"> The Aeacus Project is in no way affiliated or endorsed by the Air Force Association or the University of Texas at San Antonio.</p></div></div></div></div></div></div></div></div></div></div></div></div></div></div></div></div><div class="footer"><div class="footer-copyright-wrap"><div class="container"><div class="footer-copyright-content"><ul><li>Copyright Never &copy;</li><li>No rights reserved</li></ul></div></div></div></div></body></html>  `

	header_the_sequel := `<p> Please read the entire README thoroughly before modifying anything on this computer.</p><h2>Unique Identifier</h2><p> If you have not yet entered a valid Team ID, please do so immediately by double clicking on the "Aeacus Set Team ID" icon on the desktop. If you do not enter a valid Team ID this VM may stop functioning after a short period of time.</p><h2>Forensics Questions</h2><p> If there are "Forensics Questions" on your Desktop, you will receive points for answering these questions correctly. Valid (scored) "Forensics Questions" will only be located directly on your Desktop. Please read all "Forensics Questions" thoroughly before modifying this computer, as you may change something that prevents you from answering the question correctly.</p><h2>Competition Scenario</h2><p> This company's security policies require that all user accounts be password protected. Employees are required to choose secure passwords, however this policy may not be currently enforced on this computer. The presence of any non-work related media files and "hacking tools" on any computers is strictly prohibited. This company currently does not use any centralized maintenance or polling tools to manage their IT equipment. This computer is for official business use only by authorized users. This is a critical computer in a production environment. Please do <b>NOT</b> attempt to upgrade the operating system on this machine.</p>`

	var htmlFile strings.Builder
	htmlFile.WriteString(header)
	htmlFile.WriteString(fmt.Sprintf("<h1><b>%s %s README</b></h1>", mc.Config.OS, mc.Config.Title))
	htmlFile.WriteString(header_the_sequel)

	htmlFile.WriteString(fmt.Sprintf("<h2><b>%s</b></h2>", mc.Config.OS))

	htmlFile.WriteString(fmt.Sprintf(`<p>
    It is company policy to use only %s on this
    computer. It is also company policy to use only the
    latest, official, stable %s packages available
    for required software and services on this computer.
    Management has decided that the default web browser for
    all users on this computer should be the latest stable
    version of Firefox.`, mc.Config.OS, mc.Config.OS))

	if runtime.GOOS == "linux" {
		htmlFile.WriteString(` Company policy is to never let users log in as root. If administrators need to run commands as root, they are required to use the "sudo" command.`)
	}

	htmlFile.WriteString("</p>")
	userReadMe, err := readFile("ReadMe.conf")
	if err != nil {
		failPrint("No ReadMe.conf file found!")
		os.Exit(1)
	}
	htmlFile.WriteString(userReadMe)
	htmlFile.WriteString(footer)
	if mc.Cli.Bool("v") {
		infoPrint("Writing HTML to ReadMe.html...")
	}
	writeFile(mc.DirPath+"web/ReadMe.html", htmlFile.String())
}
