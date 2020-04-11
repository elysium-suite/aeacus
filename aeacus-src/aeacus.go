package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"runtime"
)

//////////////////////////////////////////////////////////////////
//  .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o //
// `P  )88b  d88' `88b `P  )88b  d88' `"Y8 `888  `888  d88  "8  //
//  .oP"888  888ooo888  .oP"888  888        888   888  `"Y88b.  //
// d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b //
// `Y888""8o `Y8bod8P' `Y888""8o `Y8bod8P'  `V88V"V8P' 8""888P' //
//////////////////////////////////////////////////////////////////

type metaConfig struct {
	Cli     *cli.Context
	TeamID  string
	DirPath string
	Config  scoringChecks
}

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
		failPrint("This operating system (" + runtime.GOOS + ") is not supported!")
		os.Exit(1)
	}

	id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0, []string{"green", "OK", "green", "OK", "green", "OK"}, false}

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
			checkConfig(&mc)
			scoreImage(&mc, &id)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Print extra information",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "(default) Score image with current scoring config",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					checkConfig(&mc)
					scoreImage(&mc, &id)
					return nil
				},
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					checkConfig(&mc)
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "Encrypt scoring.conf to scoring.dat",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					writeConfig(&mc)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "decrypt lol",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
                    parseConfig(&mc, readData(&mc))
					scoreImage(&mc, &id)
					return nil
				},
			},
			{
				Name:    "createfqs",
				Aliases: []string{"f"},
				Usage:   "Create forensic question files (3 by default)",
				Action: func(c *cli.Context) error {
					fmt.Println("todo")
					return nil
				},
			},
			{
				Name:    "gooey",
				Aliases: []string{"g"},
				Usage:   "Launch gui tests",
				Action: func(c *cli.Context) error {
					launchGui()
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					releaseImage(&mc)
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

///////////////////////
// CONTROL FUNCTIONS //
///////////////////////

func checkConfig(mc *metaConfig) {
	fileContent, err := readFile(mc.DirPath + "scoring.conf")
	if err != nil {
		failPrint("Configuration file not found!")
		os.Exit(1)
	}
	parseConfig(mc, fileContent)
	if mc.Cli.Bool("v") {
		printConfig(mc)
	}
}

func releaseImage(mc *metaConfig) {
	checkConfig(mc)
	writeConfig(mc)
	genReadMe(mc)
	writeDesktopFiles(mc)
	installService(mc)
	cleanUp(mc)
}
