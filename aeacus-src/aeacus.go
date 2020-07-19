package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/urfave/cli"
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

var teamID string
var dirPath string

func main() {
	id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0, []string{"green", "OK", "green", "OK", "green", "OK"}, false}

	if runtime.GOOS == "linux" {
		dirPath = "/opt/aeacus/"
	} else if runtime.GOOS == "windows" {
		dirPath = "C:\\aeacus\\"
	} else {
		failPrint("This operating system (" + runtime.GOOS + ") is not supported!")
		os.Exit(1)
	}

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
			timeCheck(&mc)
			runningPermsCheck()
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
					runningPermsCheck()
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					timeCheck(&mc)
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
					timeCheck(&mc)
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
					timeCheck(&mc)
					writeConfig(&mc)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "Check that scoring.dat is valid",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					timeCheck(&mc)
					decryptedData, err := tryDecodeString(readData(&mc))
					if err != nil {
						return errors.New("error in reading scoring.dat")
					}
					parseConfig(&mc, decryptedData)
					infoPrint("Config looks good! Decryption successful.")
					return nil
				},
			},
			{
				Name:    "forensics",
				Aliases: []string{"f"},
				Usage:   "Create forensic question files",
				Action: func(c *cli.Context) error {
					runningPermsCheck()
					numFqs, err := strconv.Atoi(c.Args().First())
					if err != nil {
						return errors.New("Invalid or missing number passed to forensics")
					}
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					checkConfig(&mc)
					createFQs(&mc, numFqs)
					return nil
				},
			},
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
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Get info about the system",
				Action: func(c *cli.Context) error {
					runningPermsCheck()
					getInfo(c.Args().Get(0))
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the current version of aeacus",
				Action: func(c *cli.Context) error {
					fmt.Println("=== aeacus ===")
					fmt.Println("version", aeacusVersion)
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					runningPermsCheck()
					mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
					timeCheck(&mc)
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
