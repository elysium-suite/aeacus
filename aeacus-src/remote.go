package main

import (
    "fmt"
    "net/http"
)

func getAuthToken() {
    fmt.Println("init connection")
}

func sendScore(score int) {
    // form http request
    // add auth token
    // add team id
    // add image name
    // add score
    fmt.Printf("sending %d", score)
}

func checkScoring(mc *metaConfig) bool {
    // hit endpoint with check
    return true
}

func checkServer(mc *metaConfig) ([]string, bool) {

    connStatus := make([]string, 6)

    // Internet check (requisite)
    if mc.Cli.Bool("v") {
        infoPrint("Checking for internet connection...")
    }
    _, err := http.Get("http://clients3.google.com/generate_204")
    if err != nil {
        connStatus[2] = "red"
        connStatus[3] = "FAIL"
    } else {
        connStatus[2] = "green"
        connStatus[3] = "OK"
    }

    // Scoring engine check (required)
    if mc.Cli.Bool("v") {
        infoPrint("Checking for scoring engine connection...")
    }
    _, err = http.Get("http://" + mc.Config.Remote)
    if err != nil {
        connStatus[4] = "red"
        connStatus[5] = "FAIL"
    } else {
        if checkScoring(mc) {
            connStatus[4] = "green"
            connStatus[5] = "OK"
        } else {
            connStatus[4] = "yellow"
            connStatus[5] = "ERROR"
        }
    }

    // Overall
    if connStatus[3] == "FAIL" && connStatus[5] == "OK" {
        connStatus[0] = "yellow"
        connStatus[1] = "Server connection good but no Internet. Assuming you're on an isolated LAN."
        return connStatus, true
    } else if connStatus[5] == "FAIL" {
        connStatus[0] = "red"
        connStatus[1] = "Failure! Can't access remote scoring server."
        return connStatus, false
    } else if connStatus[4] == "ERROR" {
        connStatus[0] = "red"
        connStatus[1] = "Score upload failure! Can't send scores to remote server."
        return connStatus, false
    } else {
        connStatus[0] = "green"
        connStatus[1] = "OK"
        return connStatus, true
    }

}
