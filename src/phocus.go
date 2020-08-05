package main

import (
	"math/rand"
	"time"

	"github.com/urfave/cli"
)

func phocusLoop() {
	infoPrint("Initializing engine context...")
	phocusEnvironment()
	for {
		checkTrace()
		timeCheck()
		infoPrint("Scoring image...")
		scoreImage()
		jitter := time.Duration(rand.Intn(8) + 10)
		infoPrint("Scored image, sleeping for a bit...")
		time.Sleep(jitter * time.Second)
	}
}

// phocusEnvironment runs functions needed in order for phocus to successfully
// run on first start.
func phocusEnvironment() {
	// Make sure we're running as admin.
	runningPermsCheck()
	// Fill constants (ex. mc.DirPath) based on OS.
	fillConstants()
	// Make sure phocus is not being traced or debugged.
	checkTrace()
	// Read in scoring data from the scoring data file.
	readScoringData()
	// Seed the random function for scoring at random intervals.
	rand.Seed(time.Now().UnixNano())
}

func genPhocusApp() *cli.App {
	return &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			phocusLoop()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "idprompt",
				Aliases: []string{"p"},
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
					infoPrint("=== phocus ===")
					infoPrint("version " + aeacusVersion)
					return nil
				},
			},
		},
	}
}
