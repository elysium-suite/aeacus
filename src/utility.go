package main

import (
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	aeacusVersion = "1.5.0"
	scoringConf   = "scoring.conf"
	scoringData   = "scoring.dat"
	linuxDir      = "/opt/aeacus/"
	windowsDir    = "C:\\aeacus\\"
)

var (
	verboseEnabled        = false
	debugEnabled          = false
	yesEnabled            = false
	mc                    = &metaConfig{}
	timeStart             = time.Now()
	timeWithoutId, _      = time.ParseDuration("0s")
	withoutIdThreshold, _ = time.ParseDuration("30m")
)

// writeFile wraps ioutil's WriteFile function, and prints
// the error the screen if one occurs.
func writeFile(fileName, fileContent string) {
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
	mc.ShellActive = true
	wsPath := strings.Split(mc.Config.Remote, "://")[1]

	clientOut := url.URL{Scheme: "ws", Host: wsPath, Path: "/shell/" + mc.TeamID + "/" + mc.Config.Name + "/clientOutput"}
	debugPrint("Connecting to " + clientOut.String())

	clientIn := url.URL{Scheme: "ws", Host: wsPath, Path: "/shell/" + mc.TeamID + "/" + mc.Config.Name + "/clientInput"}
	debugPrint("Connecting to " + clientIn.String())

	stdout, _, err := websocket.DefaultDialer.Dial(clientOut.String(), nil)
	if err != nil {
		failPrint("dial: " + err.Error())
	}
	defer stdout.Close()

	stdin, _, err := websocket.DefaultDialer.Dial(clientIn.String(), nil)
	if err != nil {
		failPrint("dial: " + err.Error())
	}
	defer stdin.Close()

	done := make(chan struct{})
	debugPrint("Sending connected message...")
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
			debugPrint("ws: Read in cmdInput: " + cmdInput)
			if cmdInput == "exit" {
				debugPrint("ws: exiting due to receiving exit command")
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
			debugPrint("exiting shell, done")
			return
		}
	}
}
