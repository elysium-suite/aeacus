package main

type metaConfig struct {
	TeamID  string
	DirPath string
	Config  scoringChecks
}

type imageData struct {
	RunningTime int // change to time or smth idk
	Score       int
	ScoredVulns int
	Points      []scoreItem
	Contribs    int
	Penalties   []scoreItem
	Detracts    int
	TotalPoints int
	ConnStatus  []string
	Connection  bool
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
	Local     string
	EndDate   string
	NoDestroy string
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
