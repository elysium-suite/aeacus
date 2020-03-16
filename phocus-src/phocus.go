package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"runtime"
)

/////////////////////////////////////////////////////////////////////
//            oooo                                                 //
//            `888                                                 //
// oo.ooooo.   888 .oo.    .ooooo.   .ooooo.  oooo  oooo   .oooo.o //
//  888' `88b  888P"Y88b  d88' `88b d88' `"Y8 `888  `888  d88(  "8 //
//  888   888  888   888  888   888 888        888   888  `"Y88b.  //
//  888   888  888   888  888   888 888   .o8  888   888  o.  )88b //
//  888bod8P' o888o o888o `Y8bod8P' `Y8bod8P'  `V88V"V8P' 8""888P' //
//  888                                                            //
// o888o                                                           //
/////////////////////////////////////////////////////////////////////

type metaConfig struct {
	Cli        *cli.Context
    TeamID      string
	ConfigName string
	DataName   string
	WebName    string
	Config     scoringChecks
}

func main() {

	var configName string
	var dataName string
	var webName string
	if runtime.GOOS == "linux" {
		configName = "/opt/aeacus/scoring.conf"
		dataName = "/opt/aeacus/scoring.dat"
		webName = "/opt/aeacus/web/ScoringReport.html"
	} else if runtime.GOOS == "windows" {
		configName = "C:\\aeacus\\scoring.conf"
		dataName = "C:\\aeacus\\scoring.dat"
		webName = "C:\\aeacus\\web\\ScoringReport.html"
	} else {
		failPrint("What are you doing?")
		os.Exit(1)
	}

	id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0}

    // read TeamID
    teamID := "booger"

    cli.AppHelpTemplate = "" // No help! >:(

	app := &cli.App{
		Name:                   "phocus",
		Usage:                  "score vulnerabilities",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
			scoreImage(&mc, &id)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

///////////////////////
// CONTROL FUNCTIONS //
///////////////////////

func scoreImage(mc *metaConfig, id *imageData) {
	parseConfig(mc, readData(mc))
    connStatus := []string{"green", "OK", "green", "OK", "green", "OK"}
    if mc.Config.Remote != "" {
        connStatus, connection := checkServer(mc)
        if ! connection {
            failPrint("No connection to server found!")
            genTemplate(mc, id, connStatus)
            os.Exit(1)
        }
    }
    if runtime.GOOS == "linux" {
        scoreLinux(mc, id)
    } else {
        //scoreWindows(mc, id)
        fmt.Println("score wondows")
    }
    genTemplate(mc, id, connStatus)
}
