// +build !phocus

package main

import (
	"errors"
	"log"
	"os"
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

func main() {

	fillConstants()
	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Action: func(c *cli.Context) error {
			parseFlags(c)
			runningPermsCheck()
			checkConfig(scoringConf)
			scoreImage()
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Print extra information",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Print a lot of information",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "Score image with current scoring config",
				Action: func(c *cli.Context) error {
					parseFlags(c)
					runningPermsCheck()
					checkConfig(scoringConf)
					scoreImage()
					return nil
				},
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					parseFlags(c)
					checkConfig(scoringConf)
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "Encrypt scoring configuration",
				Action: func(c *cli.Context) error {
					parseFlags(c)
					writeConfig(scoringConf, scoringData)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "Check that scoring data file is valid",
				Action: func(c *cli.Context) error {
					parseFlags(c)
					decryptedData, err := decodeString(readData(scoringData))
					if err != nil {
						return errors.New("error in reading scoring.dat")
					}
					parseConfig(decryptedData)
					if verboseEnabled {
						infoPrint("Config looks good! Decryption successful.")
					}
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
					parseFlags(c)
					checkConfig(scoringConf)
					createFQs(numFqs)
					return nil
				},
			},
			{
				Name:    "configure",
				Aliases: []string{"g"},
				Usage:   "Launch configuration GUI",
				Action: func(c *cli.Context) error {
					launchConfigGui()
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
					infoPrint("=== aeacus ===")
					infoPrint("version " + aeacusVersion)
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					runningPermsCheck()
					parseFlags(c)
					releaseImage()
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

// parseFlags sets the global variable values, for example,
// verboseEnabled.
func parseFlags(c *cli.Context) {
	if c.Bool("v") {
		verboseEnabled = true
	}
	if c.Bool("d") {
		debugEnabled = true
	}
}

// checkConfig parses and checks the validity of the current
// `scoring.conf` file.
func checkConfig(fileName string) {
	fileContent, err := readFile(mc.DirPath + fileName)
	if err != nil {
		failPrint("Configuration file (" + fileName + "not found!")
		os.Exit(1)
	}
	parseConfig(fileContent)
	if verboseEnabled {
		printConfig()
	}
}

// releaseImage goes through the process of checking the config,
// writing the ReadMe/Desktop Files, installing the system service,
// and cleaning the image for release.
func releaseImage() {
	checkConfig(scoringConf)
	writeConfig(scoringConf, scoringData)
	genReadMe()
	writeDesktopFiles()
	configureAutologin()
	installService()
	cleanUp()
}
