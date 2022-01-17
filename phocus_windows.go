//go:build phocus
// +build phocus

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/elysium-suite/aeacus/cmd"
	"github.com/judwhite/go-svc"
)

func phocusStart(quit chan struct{}) {
	go cmd.StartSocketWin()
	app := cmd.GenPhocusApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// idgui is set using the flag package in order to grab its value before
// the Windows service is initialized.
var idgui *bool = flag.Bool("i", false, "Spawn TeamID gui")

// Program implements svc.Service, for Windows Services.
type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

// main for phocus_windows.go will
func main() {
	flag.Parse()
	prg := &program{}
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

// Init for our windows service *program will prevent people from running
// the production binary outside of the Windows service.
func (p *program) Init(env svc.Environment) error {
	if !env.IsWindowsService() && !*idgui {
		fmt.Println("Sorry! Don't run this binary yourself. It's probably already running as a Windows service (CSSClient).")
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
	cmd.LaunchIDPrompt()
	os.Exit(0) // This is temporary solution
}
