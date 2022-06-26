package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func phocusLoop() {
	info("Initializing engine context...")
	phocusEnvironment()
	if conf.Shell {
		go shellSocket()
	}
	for {
		scoreImage()
		jitter := time.Duration(rand.Intn(8) + 10)
		info("Scored image, sleeping for a bit...")
		time.Sleep(jitter * time.Second)
	}
}

// phocusEnvironment runs functions needed in order for phocus to successfully
// run on first start.
func phocusEnvironment() {
	// Make sure we're running as admin.
	permsCheck()
	// Make sure phocus is not being traced or debugged.
	checkTrace()
	// Read in scoring data from the scoring data file.
	if err := readScoringData(); err != nil {
		fail(err)
		os.Exit(1)
	}
	// Seed the random function for scoring at "random" intervals.
	rand.Seed(time.Now().UnixNano())
}

// genPhocusApp generates a basic CLI interface that is OS-independent.
func genPhocusApp() *cli.App {
	return &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			phocusLoop()
			return nil
		},
		Before: func(c *cli.Context) error {
			err := determineDirectory()
			if err != nil {
				return err
			}
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
					info("phocus version " + version)
					return nil
				},
			},
		},
	}
}
