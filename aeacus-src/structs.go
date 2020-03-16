package main

type imageData struct {
	RunningTime int // change to time or smth idk
	Score       int
	ScoredVulns int
	Points      []scoreItem
	Contribs    int
	Penalties   []scoreItem
	Detracts    int
	TotalPoints int
}

type scoreItem struct {
	Message string
	Points  int
}

type scoringChecks struct {
	Name  string
	Title string
	User  string
    Remote string
    EndDate string
	Check []check
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
}
