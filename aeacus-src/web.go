package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
    "time"
)

func writeToHtml(mc *metaConfig, htmlFile string) {
	if mc.Cli.Bool("v") {
		infoPrint("Writing web output to " + mc.WebName)
	}
	err := ioutil.WriteFile(mc.WebName, []byte(htmlFile), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func genTemplate(mc *metaConfig, id *imageData, connStatus []string) {

	header := `<!DOCTYPE html> <html> <head> <title>Aeacus Scoring Report</title> <style type="text/css"> h1 { text-align: center; } h2 { text-align: center; } body { font-family: Arial, Verdana, sans-serif; font-size: 14px; margin: 0; padding: 0; width: 100%; height: 100%; background: url('background.png'); background-size: cover; background-attachment: fixed; background-position: top center; background-color: #336699; } .red {color: red;} .green {color: green;} .blue {color: blue;} .main { margin-top: 10px; margin-bottom: 10px; margin-left: auto; margin-right: auto; padding: 0px; border-radius: 12px; background-color: white; width: 900px; max-width: 100%; min-width: 600px; box-shadow: 0px 0px 12px #003366; } .text { padding: 12px; -webkit-touch-callout: none; -webkit-user-select: none; -khtml-user-select: none; -moz-user-select: none; -ms-user-select: none; user-select: none; } .center { text-align: center; } .binary { position: relative; overflow: hidden; } .binary::before { position: absolute; top: -75%; left: -125%; display: block; width: 200%; height: 150%; -webkit-transform: rotate(-45deg); -moz-transform: rotate(-45deg); -ms-transform: rotate(-45deg); transform: rotate(-45deg); content: attr(data-binary); opacity: 0.16; line-height: 2em; letter-spacing: 2px; color: #369; font-size: 10px; pointer-events: none; } </style> <meta http-equiv="refresh"> </head> <body><div class="main"><div class="text"><div class="binary" data-binary="0000 0000 11010000 01100100"><p align=center style="width:100%;text-align:center"><img align=middle style="width:180px; float:middle" src="logo.png"></p>`

	footer := `</p> <br> <p align=center style="text-align:center"> The Aeacus project is free and open source software. This project is in no way endorsed or affiliated with the Air Force Association or the University of Texas at San Antonio. </p> </div> </div> </div> </body> </html>`

	var htmlFile strings.Builder
    genTime := time.Now()
	htmlFile.WriteString(header)
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
	writeToHtml(mc, htmlFile.String())
}
