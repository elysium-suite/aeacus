// +build phocus

package main

import (
	"log"
	"os"

	"github.com/elysium-suite/aeacus/cmd"
)

func main() {
	go cmd.StartSocketLin()
	app := cmd.GenPhocusApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
