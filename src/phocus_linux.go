// +build phocus

package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
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

	var teamID string
	var dirPath string

	if !adminCheck() {
		failPrint("You need to run this binary as root or Administrator!")
		os.Exit(1)
	}
	if runtime.GOOS == "linux" {
		dirPath = "/opt/aeacus/"
	} else if runtime.GOOS == "windows" {
		dirPath = "C:\\aeacus\\"
	} else {
		failPrint("What are you up to?")
		os.Exit(1)
	}

	id := newImageData()
	mc := metaConfig{teamID, dirPath, scoringChecks{}}

	app := &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			parseConfig(&mc, readData(&mc))
			rand.Seed(time.Now().UnixNano())
			for {
				timeCheck(&mc)
				infoPrint("Scoring image...")
				scoreImage(&mc, &id)
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
