package main

import (
	"fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
	"io/ioutil"
	"github.com/fatih/color"
)

func scoreImage(mc *metaConfig, id *imageData) {
    // Check connection and configuration
    var connStatus []string
    if mc.Config.Remote != "" {
        readTeamID(mc)
        connStatus, connection := checkServer(mc)
        if ! connection {
            failPrint("Can't access remote scoring server!")
            genReport(mc, id, connStatus)
            os.Exit(1)
        }
    }

    // Score image
    if runtime.GOOS == "linux" {
        scoreL(mc, id)
    } else {
        scoreW(mc, id)
    }

    // Check if points increased/decreased
    prevPoints, err := readFile(mc.DirPath + "web/assets/previous.txt")
    if err == nil {
        prevScore, _ := strconv.Atoi(prevPoints)
        if prevScore < id.Score {
            sendNotification(mc.Config.User, "You gained points!")
            // TODO play gain noise
        } else if prevScore > id.Score {
            sendNotification(mc.Config.User, "You lost points!")
            // TODO play loss noise
        }
    } else {
        fmt.Println(err)
    }
    writeFile(mc.DirPath + "web/assets/previous.txt", strconv.Itoa(id.Score))

    genReport(mc, id, connStatus)
}

func destroyImage() {
	// destroy the image if outside time range
    if runtime.GOOS == "linux" {
    	fmt.Println("oops dfestroying thel inux system")
    } else {
        fmt.Println("cant do that yet. not supported on windows. enjoy ur undestryoed imaeg")
    }
}

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
