package main

type RankingSave struct {
	SourceTimestamp string      `json:"sourceTimestamp"`
	Teams           []SavedTeam `json:"teams"`
	MatchLogs       []Match     `json:"matchLogs"`
}

type Match struct {
	Home       SavedTeam `json:"home"`
	Away       SavedTeam `json:"away"`
	Importance int       `json:"importance"`
}
