package main

import (
	"fmt"
	"os"
)

// getInfo is a helper function to retrieve information about the
// system.
func getInfo(infoType string) {
	switch infoType {
	case "programs":
		programList, _ := getPrograms()
		for _, p := range programList {
			info(p)
		}
	case "users":
		userList, _ := getLocalUsers()
		for _, u := range userList {
			info(fmt.Sprint(u))
		}
	case "admins":
		adminList, _ := getLocalAdmins()
		for _, u := range adminList {
			info(fmt.Sprint(u))
		}
	default:
		if infoType == "" {
			fail("No info type provided. See the README for supported types.")
		} else {
			fail("No info for \"" + infoType + "\" found.")
		}
		os.Exit(1)
	}
}
