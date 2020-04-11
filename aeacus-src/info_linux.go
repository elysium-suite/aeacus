package main

// gather info about system (packages, other stuff that is easy to mess up)

import (
    "os"
    "fmt"
)

func getInfo(infoType string) {
    switch infoType {
    case "packages":
        ListPackages()
    default:
        if infoType == "" {
            failPrint("No info type provided.")
        } else {
            failPrint("No info for \"" + infoType + "\" found.")
        }
        os.Exit(1)
    }
}

func ListPackages() {
    fmt.Println("just run dpkg -l")
}
