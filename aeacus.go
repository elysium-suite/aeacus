// +build !phocus

package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/elysium-suite/aeacus/cmd"
	"github.com/urfave/cli/v2"
)

//////////////////////////////////////////////////////////////////
//  .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o //
// `P  )88b  d88' `88b `P  )88b  d88' `"Y8 `888  `888  d88  "8  //
//  .oP"888  888ooo888  .oP"888  888        888   888  `"Y88b.  //
// d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b //
// `Y888""8o `Y8bod8P' `Y888""8o `Y8bod8P'  `V88V"V8P' 8""888P' //
//////////////////////////////////////////////////////////////////

func main() {
	cmd.FillConstants()
	cmd.RunningPermsCheck()
	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Before: func(c *cli.Context) error {
			cmd.ParseFlags(c)
			return nil
		},
		Action: func(c *cli.Context) error {
			cmd.CheckConfig(cmd.ScoringConf)
			cmd.ScoreImage()
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
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Automatically answer 'yes' to any prompts",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "Score image with current scoring config",
				Action: func(c *cli.Context) error {
					cmd.CheckConfig(cmd.ScoringConf)
					cmd.ScoreImage()
					return nil
				},
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					cmd.CheckConfig(cmd.ScoringConf)
					return nil
				},
			},
			{
				Name:    "readme",
				Aliases: []string{"c"},
				Usage:   "Compile the readme",
				Action: func(c *cli.Context) error {
					cmd.GenReadMe()
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"c"},
				Usage:   "Score the image and render a readme",
				Action: func(c *cli.Context) error {
					cmd.GenReadMe()
					cmd.CheckConfig(cmd.ScoringConf)
					cmd.ScoreImage()
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "Encrypt scoring configuration",
				Action: func(c *cli.Context) error {
					cmd.WriteConfig(cmd.ScoringConf, cmd.ScoringData)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "Check that scoring data file is valid",
				Action: func(c *cli.Context) error {
					err := cmd.ReadScoringData()
					return err
				},
			},
			{
				Name:    "forensics",
				Aliases: []string{"f"},
				Usage:   "Create forensic question files",
				Action: func(c *cli.Context) error {
					numFqs, err := strconv.Atoi(c.Args().First())
					if err != nil {
						return errors.New("Invalid or missing number passed to forensics")
					}
					cmd.CheckConfig(cmd.ScoringConf)
					cmd.CreateFQs(numFqs)
					return nil
				},
			},
			{
				Name:    "configure",
				Aliases: []string{"g"},
				Usage:   "Launch configuration GUI",
				Action: func(c *cli.Context) error {
					cmd.LaunchConfigGui()
					return nil
				},
			},
			{
				Name:    "idprompt",
				Aliases: []string{"p"},
				Usage:   "Launch TeamID GUI prompt",
				Action: func(c *cli.Context) error {
					cmd.LaunchIDPrompt()
					return nil
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Get info about the system",
				Action: func(c *cli.Context) error {
					cmd.SetVerbose(true)
					cmd.GetInfo(c.Args().Get(0))
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the current version of aeacus",
				Action: func(c *cli.Context) error {
					println("=== aeacus ===")
					println("version " + cmd.AeacusVersion)
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					if !cmd.YesEnabled {
						cmd.ConfirmPrint("Are you sure you want to begin the image release process?")
					}
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

// releaseImage goes through the process of checking the config,
// writing the ReadMe/Desktop Files, installing the system service,
// and cleaning the image for release.
func releaseImage() {
	cmd.CheckConfig(cmd.ScoringConf)
	cmd.WriteConfig(cmd.ScoringConf, cmd.ScoringData)
	cmd.GenReadMe()
	cmd.WriteDesktopFiles()
	cmd.ConfigureAutologin()
	cmd.InstallFont()
	cmd.InstallService()
	cmd.CleanUp()
}
