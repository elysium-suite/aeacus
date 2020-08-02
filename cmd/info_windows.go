package cmd

import (
	wapi "github.com/iamacarpet/go-win64api"
	"os"
)

func getInfo(infoType string) {
	switch infoType {
	case "packages":
		packageList := getPackages()
		for _, p := range packageList {
			infoPrint(p)
		}
	default:
		if infoType == "" {
			failPrint("No info type provided.")
		} else {
			failPrint("No info for \"" + infoType + "\" found.")
		}
		os.Exit(1)
	}
}

func getPackages() []string {
	sw, err := wapi.InstalledSoftwareList()
	if err != nil {
		failPrint(err.Error())
	}
	softwareList := []string{}
	for _, s := range sw {
		softwareList = append(softwareList, s.Name())
	}
	return softwareList
}
