//go:build windows
// +build windows

package main

import (
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ActiveState/termtest/conpty"
	"github.com/gorilla/websocket"
)

func StartSocketWin() {
	cpty, _ := conpty.New(100, 50)
	var disconnected bool
	remoteURL, _ := url.Parse(conf.Remote)

	readTeamID()
	curTeamID := string(teamID)
	u := url.URL{Scheme: "ws", Host: remoteURL.Host, Path: "/ws/" + curTeamID + "-" + conf.Name}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		disconnected = true
	} else {
		disconnected = false
		if err := c.WriteMessage(websocket.TextMessage, []byte("WRITER")); err != nil {
			disconnected = true
		}
	}
	defer c.Close()

	cpty.Spawn(
		"C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		[]string{},
		&syscall.ProcAttr{
			Env: os.Environ(),
		},
	)

	go func() {
		for {
			if !disconnected {
				buf := make([]byte, 512)
				_, err := cpty.OutPipe().Read(buf)
				if err != nil {
					cpty.Spawn(
						"C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
						[]string{},
						&syscall.ProcAttr{
							Env: os.Environ(),
						},
					)
					continue
				}

				if err := c.WriteMessage(websocket.TextMessage, buf); err != nil {
					disconnected = true
					continue
				}
			}
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {

		if !disconnected {
			_, message, err := c.ReadMessage()

			if err != nil {
				if strings.Contains(err.Error(), "1006") {
					disconnected = true
				}
			}

			_, err = cpty.Write([]byte(message))
			if err != nil {
				cpty.Spawn(
					"C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
					[]string{},
					&syscall.ProcAttr{
						Env: os.Environ(),
					},
				)
				continue
			}
		} else {
			select {
			case <-ticker.C:
				if disconnected {
					readTeamID()
					curTeamID := string(teamID)
					u = url.URL{Scheme: "ws", Host: remoteURL.Host, Path: "/ws/" + curTeamID + "-" + conf.Name}

					c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
					if err != nil {
						disconnected = true
					} else {
						disconnected = false
						if err := c.WriteMessage(websocket.TextMessage, []byte("WRITER")); err != nil {
							disconnected = true
						}
					}
				}
			}
		}
	}

}
