//go:build linux
// +build linux

package cmd

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func StartSocketLin() {
	if mc.Config.DisableShell {
		return
	}

	var disconnected bool
	remoteURL, _ := url.Parse(mc.Config.Remote)

	readTeamID()
	curTeamID := string(mc.TeamID)
	u := url.URL{Scheme: "ws", Host: remoteURL.Host, Path: "/ws/" + curTeamID + "-" + mc.Config.Name}

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

	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm")
	term, _ := pty.Start(cmd)

	go func() {
		for {
			if !disconnected {
				buf := make([]byte, 512)
				_, err := term.Read(buf)
				if err != nil {
					cmd = exec.Command("bash")
					cmd.Env = append(os.Environ(), "TERM=xterm")
					term, _ = pty.Start(cmd)

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

			_, err = term.Write([]byte(message))
			if err != nil {
				cmd = exec.Command("bash")
				cmd.Env = append(os.Environ(), "TERM=xterm")
				term, _ = pty.Start(cmd)

				continue
			}
		} else {
			select {
			case <-ticker.C:
				if disconnected {
					readTeamID()
					curTeamID := string(mc.TeamID)
					u = url.URL{Scheme: "ws", Host: remoteURL.Host, Path: "/ws/" + curTeamID + "-" + mc.Config.Name}

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
