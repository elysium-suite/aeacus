package cmd

import (
	"fmt"
	"os"
)

// GetInfo is a helper function to retrieve
// generic information about the system
func GetInfo(infoType string) {
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
