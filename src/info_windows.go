package main

import (
	"fmt"
	"os"
)

func getInfo(infoType string) {
	switch infoType {
	case "packages":
		packageList, _ := getPackages()
		for _, p := range packageList {
			infoPrint(p)
		}
	case "users":
		userList, _ := getLocalUsers()
		for _, u := range userList {
			infoPrint(fmt.Sprint(u))
		}
	case "admins":
		adminList, _ := getLocalAdmins()
		for _, u := range adminList {
			infoPrint(fmt.Sprint(u))
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
