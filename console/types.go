package main

type RankingTime struct {
	Timestamp string      `json:"timestamp"`
	Teams     []SavedTeam `json:"teams"`
}

type SavedTeam struct {
	Name          string `json:"name"`
	FifaCode      string `json:"fifaCode"`
	Confederation string `json:"confederation"`
	Points        int    `json:"points"`
}

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
