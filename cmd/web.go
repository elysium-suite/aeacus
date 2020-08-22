package cmd

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func genReport(img imageData) {
	teamID := mc.TeamID
	if len(teamID) < 7 {
		teamID = "1010 1101"
	}
	header := `<!DOCTYPE html> <html> <head> <meta http-equiv="refresh" content="60"> <title>Aeacus Scoring Report</title> <style type="text/css"> h1 { text-align: center; } h2 { text-align: center; } body { font-family: Arial, Verdana, sans-serif; font-size: 14px; margin: 0; padding: 0; width: 100%; height: 100%; background: url('background.png'); background-size: cover; background-attachment: fixed; background-position: top center; background-color: #336699; } .red {color: red;} .green {color: green;} .blue {color: blue;} .main { margin-top: 10px; margin-bottom: 10px; margin-left: auto; margin-right: auto; padding: 0px; border-radius: 12px; background-color: white; width: 900px; max-width: 100%; min-width: 600px; box-shadow: 0px 0px 12px #003366; } .text { padding: 12px; -webkit-touch-callout: none; -webkit-user-select: none; -khtml-user-select: none; -moz-user-select: none; -ms-user-select: none; user-select: none; } .center { text-align: center; } .binary { position: relative; overflow: hidden; } .binary::before { position: absolute; top: -1000px; left: -1000px; display: block; width: 500%; height: 300%; -webkit-transform: rotate(-45deg); -moz-transform: rotate(-45deg); -ms-transform: rotate(-45deg); transform: rotate(-45deg); content: attr(data-binary); opacity: 0.15; line-height: 2em; letter-spacing: 2px; color: #369; font-size: 10px; pointer-events: none; } </style> <meta http-equiv="refresh"> </head> <body><div class="main"><div class="text"><div class="binary" data-binary="` + teamID + `"><p align=center style="width:100%;text-align:center"><img align=middle style="width:180px; float:middle" src="logo.png"></p>`

	footer := `</p> <br> <p align=center style="text-align:center"> The Aeacus project is free and open source software. This project is in no way endorsed or affiliated with the Air Force Association or the University of Texas at San Antonio. </p> </div> </div> </div> </body> </html>`

	var htmlFile strings.Builder
	htmlFile.WriteString(header)
	genTime := time.Now()
	htmlFile.WriteString("<h1>" + mc.Config.Title + "</h1>")
	htmlFile.WriteString("<h2>Report Generated At: " + genTime.Format("2006/01/02 15:04:05 MST") + " </h2>")
	htmlFile.WriteString(`<script>var bin = document.querySelectorAll('.binary'); [].forEach.call(bin, function(el) { el.dataset.binary = Array(10000).join(el.dataset.binary + ' ') }); var currentdate = new Date().getTime(); gendate = Date.parse('0000/00/00 00:00:00 UTC'); diff = Math.abs(currentdate - gendate); if ( gendate > 0 && diff > 1000 * 60 * 5 ) { document.write('<span style="color:red"><h2>WARNING: CCS Scoring service may not be running</h2></span>'); } </script>`)

	if mc.Config.Remote != "" {
		htmlFile.WriteString(`<h3 class="center">Current Team ID: ` + mc.TeamID + `</h3>`)
	}

	htmlFile.WriteString(fmt.Sprintf(`<h2> %d out of %d points received</h2>`, img.Score, img.TotalPoints))

	if mc.Config.Remote != "" {
		htmlFile.WriteString(`<a href="` + mc.Config.Remote + `">Click here to view the public scoreboard</a><br>`)
		htmlFile.WriteString(`<a href="` + mc.Config.Remote + `/announcements` + `">Click here to view the announcements</a><br>`)

		htmlFile.WriteString(`<p><h3>Connection Status: <span style="color:` + mc.Conn.OverallColor + `">` + mc.Conn.OverallStatus + `<span></h3>`)

		htmlFile.WriteString(`Internet Connectivity Check: <span style="color:` + mc.Conn.NetColor + `">` + mc.Conn.NetStatus + `</span><br>`)
		htmlFile.WriteString(`Aeacus Server Connection Status: <span style="color:` + mc.Conn.ServerColor + `">` + mc.Conn.ServerStatus + `</span></p>`)
	} else {
		htmlFile.WriteString(`<p><h3>Connection Status: <span style="color:` + mc.Conn.OverallColor + `">` + mc.Conn.OverallStatus + `<span></h3>`)
		htmlFile.WriteString(`Internet Connectivity Check: <span style="color:grey">N/A</span><br>`)
		htmlFile.WriteString(`Aeacus Server Connection Status: <span style="color:grey">N/A</span><br>`)
	}

	htmlFile.WriteString(fmt.Sprintf(`<h3> %d penalties assessed, for a loss of %.0f points: </h3> <p> <span style="color:red">`, len(img.Penalties), math.Abs(float64(img.Detracts))))

	// for each penalty
	for _, penalty := range img.Penalties {
		htmlFile.WriteString(fmt.Sprintf("%s - %.0f pts<br>", penalty.Message, math.Abs(float64(penalty.Points))))
	}

	htmlFile.WriteString(fmt.Sprintf(`</span> </p> <h3> %d out of %d scored security issues fixed, for a gain of %d points:</h3><p>`, len(img.Points), img.ScoredVulns, img.Contribs))

	// for each point:
	for _, point := range img.Points {
		deobfuscateData(&point.Message)
		htmlFile.WriteString(fmt.Sprintf("%s - %d pts<br>", point.Message, point.Points))
		obfuscateData(&point.Message)
	}

	htmlFile.WriteString(footer)

	infoPrint("Writing HTML to ScoringReport.html...")
	writeFile(mc.DirPath+"assets/ScoringReport.html", htmlFile.String())
}

func GenReadMe() {
	header := `
<!DOCTYPE html>
<html>

<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Aeacus README</title>
	<style>
        body {
            background-image: url("./background.png");
            background-size: cover;
            font-family: Helvetica, Arial, sans-serif;
        }

        h1 {
	        text-align: center;
	        font-size: 36px;
	        margin: 10px;
	        padding: 0 14px 10px 0px;
	        width: 100%;
	        height: 100%;
	        color: #0D2E5B !important
        }

        h2 {
	        font-size: 18px;
	        margin: 30px 0 10px 0;
	        padding: 0;
	        width: 100%;
	        height: 100%;
	        color: #0D2E5B !important
        }

        pre {
	        font-size: 16px
        }

        .main {
	        margin-top: 25px;
	        margin-bottom: 10px;
	        margin-left: auto;
	        margin-right: auto;
	        padding: 0px;
	        background-color: white;
	        max-width: 100%
        }

        .text {
	        padding-top: 12px;
	        padding-bottom: 12px;
	        padding-left: 40px;
	        padding-right: 40px
        }

        .center {
	        text-align: center
        }
	</style>
        <div style="width: 80%; margin-left: auto; margin-right: auto; display: block" >
			<div class="main">
				<div class="text">
					<p align="center">
						<img src="./logo.png "width="180">
					</p>
`

	footer := `<h2>Competition Guidelines</h2><ul><li> In order to provide a better competition experience, you are <b>NOT</b> required to change the password of the primary, auto-login, user account. Changing the password of a user that is set to automatically log in may lock you out of your computer.</li><li> Authorized administrator passwords were correct the last time you did a password audit, but are not guaranteed to be currently accurate.</li><li> Do not disable or stop the CSSClient service or process.</li><li> Do not remove any authorized users or their home directories.</li><li> The time zone of this image is set to UTC. Please do not change the time zone, date, or time on this image.</li><li> You can view your current scoring report by double-clicking the "Scoring Report" desktop icon.</li><li> JavaScript is required for some error messages that appear on the "Scoring Report." To ensure that you only receive correct error messages, please do not disable JavaScript.</li><li> Some security settings may prevent the Stop Scoring application from running. If this happens, the safest way to stop scoring is to suspend the virtual machine. You should <b>NOT</b> power on the VM again before deleting.</li></ul><p align="center" style="text-align:center"> The Aeacus Project is in no way affiliated or endorsed by the Air Force Association or the University of Texas at San Antonio.</p></div></div></div></div></div></div></div></div></div></div></div></div></div></div></div></div><div class="footer"><div class="footer-copyright-wrap"><div class="container"><div class="footer-copyright-content"><ul><li>Copyright Never &copy;</li><li>No rights reserved</li></ul></div></div></div></div></body></html>  `

	headerTheSequel := `<p> Please read the entire README thoroughly before modifying anything on this computer.</p><h2>Unique Identifier</h2><p> If you have not yet entered a valid Team ID, please do so immediately by double clicking on the "Team ID" icon on the desktop. If you do not enter a valid Team ID this VM may stop functioning after a short period of time.</p><h2>Forensics Questions</h2><p> If there are "Forensics Questions" on your Desktop, you will receive points for answering these questions correctly. Valid (scored) "Forensics Questions" will only be located directly on your Desktop. Please read all "Forensics Questions" thoroughly before modifying this computer, as you may change something that prevents you from answering the question correctly.</p><h2>Competition Scenario</h2><p> This company's security policies require that all user accounts be password protected. Employees are required to choose secure passwords, however this policy may not be currently enforced on this computer. The presence of any non-work related media files and "hacking tools" on any computers is strictly prohibited. This company currently does not use any centralized maintenance or polling tools to manage their IT equipment. This computer is for official business use only by authorized users. This is a critical computer in a production environment. Please do <b>NOT</b> attempt to upgrade the operating system on this machine.</p>`

	var htmlFile strings.Builder
	htmlFile.WriteString(header)
	htmlFile.WriteString("<h1><b>" + mc.Config.OS + " " + mc.Config.Title + " README</b></h1>")
	htmlFile.WriteString(headerTheSequel)

	htmlFile.WriteString("<h2><b>" + mc.Config.OS + "</b></h2>")

	htmlFile.WriteString(`<p>
    It is company policy to use only ` + mc.Config.OS + ` on this
    computer. It is also company policy to use only the
    latest, official, stable ` + mc.Config.OS + ` packages available
    for required software and services on this computer.
    Management has decided that the default web browser for
    all users on this computer should be the latest stable
    version of Firefox.`)

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
	infoPrint("Writing HTML to ReadMe.html...")
	writeFile(mc.DirPath+"assets/ReadMe.html", htmlFile.String())
}
