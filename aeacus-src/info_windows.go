package main

import (
    "os"
    "fmt"
    wapi "github.com/iamacarpet/go-win64api"
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
    sw, err := wapi.InstalledSoftwareList()
    if err != nil {
        fmt.Printf("%s\r\n", err.Error())
    }

    for _, s := range sw {
        infoPrint(s.Name())
    }
}
