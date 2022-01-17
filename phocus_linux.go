//go:build phocus

package main

import (
	"log"
	"os"
)

func main() {
	app := genPhocusApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
