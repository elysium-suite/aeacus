package main

import (
	"os"
	"log"
    "time"
	"runtime"
    "math/rand"
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

type metaConfig struct {
	Cli        *cli.Context
    TeamID      string
	DirPath    string
	Config     scoringChecks
}

func main() {

    var dirPath string
    if runtime.GOOS == "linux" {
        if ! adminCheckL() {
            failPrint("You need to run this binary as root!")
            os.Exit(1)
        }
    	dirPath = "/opt/aeacus/"
    } else if runtime.GOOS == "windows" {
        if ! adminCheckW() {
            failPrint("You need to run this binary as Administrator!")
        }
        dirPath = "C:\\aeacus\\"
    } else {
        failPrint("What are you doing?")
        os.Exit(1)
    }

    cli.AppHelpTemplate = "" // No help! >:(
	app := &cli.App{
		Name:                   "phocus",
		Usage:                  "score vulnerabilities",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, "", dirPath, scoringChecks{}}
            parseConfig(&mc, readData(&mc))
            rand.Seed(time.Now().UnixNano())
            for {
                id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0}
    			scoreImage(&mc, &id)
                jitter := time.Duration(rand.Intn(20) + 6)
                time.Sleep(jitter * time.Second)
            }
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
