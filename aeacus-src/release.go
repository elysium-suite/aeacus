package main

import (
	"fmt"
)

func cleanUp() {
	infoPrint("Cleaning up the system...")
	// viminfo, scoring.conf, etc

	//    <Execute>rm -f /home/*/.local/share/recently-used.xbel</Execute>
	//    <Execute>echo Running installation commands</Execute>
	//    <Execute>rm -f /home/*/Desktop/*~</Execute>
	//    <Execute>rm -f /var/crash/*.crash</Execute>
	//    <Execute>rm -f /var/VMwareDnD/*</Execute>
}

func destroyImage() {
	// destroy the image if outside time range or time limit
	fmt.Println("destroying the system lol")
}
