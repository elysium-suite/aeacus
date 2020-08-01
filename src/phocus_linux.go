// +build phocus

package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
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

func main() {

	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}

	checkTrace()
	fillConstants()

	daemonTest := os.Getenv("INVOCATION_ID")
	if daemonTest == "" {
		failPrint("Sorry! You're not supposed to run this binary manually. It's probably already running as a daemon (CSSClient).")
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			parseConfig(readData(scoringData))
			rand.Seed(time.Now().UnixNano())
			for {
				checkTrace()
				timeCheck()
				infoPrint("Scoring image...")
				scoreImage()
				jitter := time.Duration(rand.Intn(8) + 10)
				infoPrint("Scored image, sleeping for a bit...")
				time.Sleep(jitter * time.Second)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "idprompt",
				Aliases: []string{"d"},
				Usage:   "Launch TeamID GUI prompt",
				Action: func(c *cli.Context) error {
					launchIDPrompt()
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the current version of phocus",
				Action: func(c *cli.Context) error {
					infoPrint("=== phocus (linux) ===")
					infoPrint("version " + aeacusVersion)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
