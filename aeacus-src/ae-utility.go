package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os/exec"
	"runtime"
)

func sendNotification(userName string, notifyText string) {
	if runtime.GOOS == "linux" {
		commandText := "/sbin/runuser -l " + userName + " -c  '/usr/bin/notify-send -i /opt/aeacus/web/assets/logo.png \"Aeacus Scoring System\" \"" + notifyText + "\"'"
		cmd := exec.Command("sh", "-c", commandText)
		cmd.Run()
	} else {
		fmt.Println("not supported yet oopsies")
	}
}

func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	return string(fileContent), err
}

func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func passPrint(toPrint string) {
	printer(color.FgGreen, "PASS", toPrint)
}

func failPrint(toPrint string) {
	printer(color.FgRed, "FAIL", toPrint)
}

func warnPrint(toPrint string) {
	printer(color.FgYellow, "WARN", toPrint)
}

func infoPrint(toPrint string) {
	printer(color.FgBlue, "INFO", toPrint)
}

func printer(colorChosen color.Attribute, messageType string, toPrint string) {
	printer := color.New(colorChosen, color.Bold)
	fmt.Printf("[")
	printer.Printf(messageType)
	fmt.Printf("] %s", toPrint)
	if toPrint != "" {
		fmt.Printf("\n")
	}
}
