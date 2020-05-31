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

// grab idgui flag
var idgui *string = flag.String("i", "", "Spawn TeamID gui")

// program implements svc.Service
type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

func main() {
	flag.Parse()
	prg := &program{}

	// Call svc.Run to start your program/service.
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	//if ! env.IsWindowsService() && *idgui != "yes" {
	//    failPrint("Sorry! You need to run this as a Windows service.")
	//    os.Exit(1)
	//}
	return nil
}

func (p *program) Start() error {
	p.quit = make(chan struct{})
	p.wg.Add(1)
	if *idgui == "yes" {
		go launchIDPromptWrapper(p.quit)
	} else {
		go phocusStart(p.quit)
	}
	return nil
}

func (p *program) Stop() error {
	log.Println("Stopping...")
	close(p.quit)
	os.Exit(1) // Causes windows stopping error, but it stops
	// Quit struct doesn't work... todo
	p.wg.Wait()
	log.Println("Stopped.")
	return nil
}

type metaConfig struct {
	Cli     *cli.Context
	TeamID  string
	DirPath string
	Config  scoringChecks
}

var teamID string
var dirPath string

func launchIDPromptWrapper(quit chan struct{}) {
	launchIDPrompt()
	os.Exit(0) // Kind of ghetto-- would prefer actually
	// using the quit struct
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

	cli.AppHelpTemplate = "" // No help! >:(
	app := &cli.App{
		Name:  "phocus",
		Usage: "score vulnerabilities",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, teamID, dirPath, scoringChecks{}}
			decryptedData, err := tryDecodeString(readData(&mc))
			if err != nil {
				return errors.New("Error in reading scoring.dat!")
			}
			parseConfig(&mc, decryptedData)
			rand.Seed(time.Now().UnixNano())
			for {
				id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0, []string{"green", "OK", "green", "OK", "green", "OK"}, false}
				infoPrint("Scoring image...")
				scoreImage(&mc, &id)
				jitter := rand.Intn(6) + 6
				infoPrint("Scored image, sleeping for a bit...")
				for s := 0; s < jitter; s++ {
					time.Sleep(1 * time.Second)
					// Todo: Check every second if Windows wants us to die
					//select {
					//case <-quit:
					//	break
					//}
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
