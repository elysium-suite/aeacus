//go:build !phocus

package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

//////////////////////////////////////////////////////////////////
//  .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o //
// `P  )88b  d88' `88b `P  )88b  d88' `"Y8 `888  `888  d88  "8  //
//  .oP"888  888ooo888  .oP"888  888        888   888  `"Y88b.  //
// d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b //
// `Y888""8o `Y8bod8P' `Y888""8o `Y8bod8P'  `V88V"V8P' 8""888P' //
//////////////////////////////////////////////////////////////////

const (
	DEBUG_BUILD = true
)

func main() {
	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "score image vulnerabilities",
		Before: func(c *cli.Context) error {
			if debugEnabled {
				verboseEnabled = true
			}
			err := determineDirectory()
			if err != nil {
				return err
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			permsCheck()
			readConfig()
			scoreImage()
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Print extra information",
				Destination: &verboseEnabled,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "Print a lot of information",
				Destination: &debugEnabled,
			},
			&cli.BoolFlag{
				Name:        "yes",
				Aliases:     []string{"y"},
				Usage:       "Automatically answer 'yes' to any prompts",
				Destination: &yesEnabled,
			},
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"r"},
				Usage:       "Directory for aeacus and its files",
				Destination: &dirPath,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "Score image with current scoring config",
				Action: func(c *cli.Context) error {
					permsCheck()
					readConfig()
					scoreImage()
					return nil
				},
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					readConfig()
					return nil
				},
			},
			{
				Name:    "readme",
				Aliases: []string{"rd"},
				Usage:   "Compile the README",
				Action: func(c *cli.Context) error {
					permsCheck()
					readConfig()
					genReadMe()
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "Encrypt scoring configuration",
				Action: func(c *cli.Context) error {
					permsCheck()
					readConfig()
					writeConfig()
					return nil
				},
			},
			{
				Name:    "prompt",
				Aliases: []string{"p"},
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
					permsCheck()
					verboseEnabled = true
					getInfo(c.Args().Get(0))
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the current version of aeacus",
				Action: func(c *cli.Context) error {
					println("aeacus version " + version)
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					permsCheck()
					confirm("Are you sure you want to begin the image release process?")
					releaseImage()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fail(err.Error())
	}
}

// releaseImage goes through the process of checking the config,
// writing the ReadMe/Desktop Files, installing the system service,
// and cleaning the image for release.
func releaseImage() {
	readConfig()
	writeConfig()
	genReadMe()
	writeDesktopFiles()
	configureAutologin()
	installFont()
	installService()
	confirm("Everything is done except cleanup. Are you sure you want to continue, and remove your scoring configuration and other aeacus files?")
	cleanUp()
}
