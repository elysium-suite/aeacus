package main

// gather info about system (packages, other stuff that is easy to mess up)

import (
	"fmt"
	"os"
)

func getInfo(infoType string) {
	switch infoType {
	case "packages":
		ListPackages()
	case "kernel":
		fmt.Println(shellCommandOutput("uname -r"))
	default:
		if infoType == "" {
			failPrint("No info type provided.")
		} else {
			failPrint("No info for \"" + infoType + "\" found.")
		}
		os.Exit(1)
	}
}

func ListPackages() string {
	// fmt.Println("just run dpkg -l") bro what if i'm too lazy ~ safin
	pkgs = shellCommand("dpkg -l")
	return pkgs
}
