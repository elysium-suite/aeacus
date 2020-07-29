// +build phocus

package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/judwhite/go-svc/svc"
	"github.com/urfave/cli"
)

/////////////////////////////////////////////////////////////////////
//            oooo                                                 //
//            `888                                                 //
// oo.ooooo.   888 .oo.    .ooooo.   .ooooo.  oooo  oooo   .oooo.o //
//  888' `88b  888P"Y88b  d88' `88b d88' `"Y8 `888  `888  d88(  "8 //
//  888   888  888   888  888   888 888        888   888  `"Y88b.  //
//  888   888  888   888  888   888 888   .o8  888   888  o.  )88b //
//  888bod8P' o888o o888o `Y8bod8P' `Y8bod8P'  `V88V"V8P' 8""888P' //
//  888       production binary for aeacus                         //
// o888o                                                           //
/////////////////////////////////////////////////////////////////////

// idgui is set using the flag package in order to grab its value before
// the Windows service is initialized.
var idgui *bool = flag.Bool("i", false, "Spawn TeamID gui")

// Program implements svc.Service, for Windows Services.
type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

func main() {
	flag.Parse()
	prg := &program{}
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if !env.IsWindowsService() && !*idgui {
		failPrint("Sorry! Don't run this binary. It's probably already running as a Windows service (CSSClient).")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	return nil
}

func (p *program) Start() error {
	p.quit = make(chan struct{})
	p.wg.Add(1)
	if *idgui {
		go launchIDPromptWrapper(p.quit)
	} else {
		go phocusStart(p.quit)
	}
	return nil
}

func (p *program) Stop() error {
	log.Println("Stopping...")
	close(p.quit)
	os.Exit(1)
	/*
		Causes windows service stopping error (but it works)
		Quit struct doesn't work... todo
		p.wg.Wait()
		log.Println("Stopped.")
	*/
	return nil
}

func launchIDPromptWrapper(quit chan struct{}) {
	launchIDPrompt()
	os.Exit(0) // This is temporary solution
}

func phocusStart(quit chan struct{}) {

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
		failPrint("What are you up to?")
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			decryptedData, err := decodeString(readData(&mc))
			if err != nil {
				return errors.New("Error in reading scoring.dat!")
			}
			parseConfig(decryptedData)
			rand.Seed(time.Now().UnixNano())
			for {
				timeCheck()
				mc.Image = imageData{}
				infoPrint("Scoring image...")
				scoreImage(&id)
				jitter := time.Duration(rand.Intn(8) + 8)
				infoPrint("Scored image, sleeping for a bit...")
				time.Sleep(jitter * time.Second)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the current version of phocus",
				Action: func(c *cli.Context) error {
					infoPrint("=== phocus (windows) ===")
					infoPrint("version " + aeacusVersion)
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
