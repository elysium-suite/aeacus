package main

import (
    "fmt"
)


func writeDesktopFilesW(mc *metaConfig) {
    if mc.Cli.Bool("v") {
    	infoPrint("Writing shortcuts to Desktop...")
    }

    fmt.Println("xxd")

    if mc.Cli.Bool("v") {
    	infoPrint("Creating TeamID.txt file...")
    }

    fmt.Println("xxd")
}

func installServiceW(mc *metaConfig) {
    if mc.Cli.Bool("v") {
    	infoPrint("Installing service...")
    }
}

func cleanUpW(mc *metaConfig) {

}
