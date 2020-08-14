package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	aeacusVersion = "1.4.0"
	scoringConf   = "scoring.conf"
	scoringData   = "scoring.dat"
	linuxDir      = "/opt/aeacus/"
	windowsDir    = "C:\\aeacus\\"
)

var (
	verboseEnabled = false
	debugEnabled   = false
	yesEnabled     = false
	mc             = &metaConfig{}
)

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func writeFile(fileName string, fileContent string) {
	err := ioutil.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		failPrint("Error writing file: " + err.Error())
	}
}

// grepString acts like grep, taking in a pattern to search for, and the
// fileText to search in. It returns the line which contains the string
// (if any).
func grepString(patternText, fileText string) string {
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + patternText + ".*$")
	return string(re.Find([]byte(fileText)))
}

func connectWs() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	wsPath := strings.Split(mc.Config.Remote, "://")[1]

	u1 := url.URL{Scheme: "ws", Host: wsPath, Path: "/shell/" + mc.TeamID + "/" + mc.Config.Name + "/clientOutput"}
	debugPrint("Connecting to " + u1.String())

	u2 := url.URL{Scheme: "ws", Host: wsPath, Path: "/shell/" + mc.TeamID + "/" + mc.Config.Name + "/clientInput"}
	debugPrint("Connecting to " + u2.String())

	stdout, _, err := websocket.DefaultDialer.Dial(u1.String(), nil)
	if err != nil {
		failPrint("dial: " + err.Error())
	}
	defer stdout.Close()

	stdin, _, err := websocket.DefaultDialer.Dial(u2.String(), nil)
	if err != nil {
		failPrint("dial: " + err.Error())
	}
	defer stdin.Close()

	done := make(chan struct{})
	stdout.WriteMessage(1, []byte("Connected"))

	go func() {
		defer close(done)
		for {
			_, message, err := stdin.ReadMessage()
			if err != nil {
				failPrint("read: " + err.Error())
				return
			}

			cmdInput := strings.TrimSpace(string(message))
			if cmdInput == "exit" {
				break
			}
			output, err := shellCommandOutput(cmdInput)
			if err != nil {
				err = stdout.WriteMessage(1, []byte("ERROR: "+err.Error()))
			} else {
				err = stdout.WriteMessage(1, []byte(output))
			}
			if err != nil {
				failPrint("write: " + err.Error())
				break
			}
		}
	}()

	for {
		select {
		case <-done:
			mc.ShellActive = false
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := stdin.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			mc.ShellActive = false
			return
		}
	}
}
