package main

import "fmt"

func roundRobin(teams []string) [][][]string {
	n := len(teams)

	if n%2 == 1 {
		teams = append(teams, "")
		n++
	}

	var matchdays [][][]string

	for i := 0; i < n-1; i++ {
		var matchday [][]string

		for j := 0; j < n/2; j++ {
			if teams[j] != "" && teams[n-j-1] != "" {
				matchup := []string{teams[j], teams[n-j-1]}
				matchday = append(matchday, matchup)
			}
		}
		matchdays = append(matchdays, matchday)

		// Rotate teams
		teams = append([]string{teams[0]}, append(teams[n-1:], teams[1:n-1]...)...)
	}

	return matchdays
}

func main() {
	teams := []string{"Laos", "Myanmar", "Philippines", "Indonesia", "Vietnam"}

	schedule := roundRobin(teams)

	for i, day := range schedule {
		fmt.Printf("Day %d: %v\n", i+1, day)
	}
}
