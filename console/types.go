package main

type RankingTime struct {
	Timestamp string `json:"timestamp"`
	Teams     []Team `json:"teams"`
}

type Team struct {
	Name          string `json:"name"`
	Confederation string `json:"confederation"`
	Points        int    `json:"points"`
}

type RankingSave struct {
	SourceTimestamp string  `json:"sourceTimestamp"`
	Teams           []Team  `json:"teams"`
	MatchLogs       []Match `json:"matchLogs"`
}

type Match struct {
	Home       Team `json:"home"`
	Away       Team `json:"away"`
	Importance int  `json:"importance"`
}
