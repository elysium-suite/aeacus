package main

import (
	"time"
)

type metaConfig struct {
	TeamID  string
	DirPath string
	Config  scoringChecks
	Image   imageData
}

type imageData struct {
	RunningTime time.Time
	Score       int
	ScoredVulns int
	Points      []scoreItem
	Contribs    int
	Penalties   []scoreItem
	Detracts    int
	TotalPoints int
	Conn        connData
	Connection  bool
}

type connData struct {
	ServerColor   string
	ServerStatus  string
	NetColor      string
	NetStatus     string
	OverallColor  string
	OverallStatus string
}

type scoreItem struct {
	Message string
	Points  int
}

type scoringChecks struct {
	Name      string
	Title     string
	User      string
	OS        string
	Remote    string
	Password  string
	Local     bool
	EndDate   string
	NoDestroy bool
	Check     []check
}

type check struct {
	Message string
	Points  int
	Pass    []condition
	Fail    []condition
}

type condition struct {
	Type string
	Arg1 string
	Arg2 string
	Arg3 string
	Arg4 string
}
