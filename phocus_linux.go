// +build phocus

package main

import (
	"log"
	"os"

	"github.com/elysium-suite/aeacus/cmd"
)

func main() {
	/* Only works for systemd >= v232.
	daemonTest := os.Getenv("INVOCATION_ID")
	if daemonTest == "" {
		failPrint("Sorry! You're not supposed to run this binary manually. It's probably already running as a daemon (CSSClient).")
		os.Exit(1)
	}
	*/
	app := cmd.GenPhocusApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
