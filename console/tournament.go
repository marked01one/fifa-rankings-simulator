package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TournamentStageJson struct {
	Name             string `json:"name"`
	Format           string `json:"format"`
	Slots            int    `json:"slots"`
	Importance       int    `json:"importance"`
	InheritFromIndex int    `json:"inheritFromIndex"`
	TwoLegged        bool   `json:"twoLegged"`

	Groups [][]string `json:"groups"`
}

type TournamentJson struct {
	Name   string                `json:"name"`
	Stages []TournamentStageJson `json:"stages"`
}

type GroupState struct {
	Teams map[string]GroupStateTeam
}

type GroupStateTeam struct {
	W      int
	D      int
	L      int
	GF     int
	GA     int
	GD     int
	Points int
}

func (tournament *TournamentJson) GetJson(path string) {

	bytes, err := os.ReadFile(path)

	if err != nil {
		log.Printf("Unable to read file from path '%s'", path)
		log.Fatal(err)
	}

	if json.Unmarshal(bytes, &tournament) != nil {
		log.Printf("Unable to parse TournamentJson!")
		log.Fatal(err)
	}
}

func (tournament *TournamentJson) PrintContent() {
	fmt.Printf("\n%s\n", tournament.Name)
	fmt.Print("-------------------------\n\n")

	for i, stage := range tournament.Stages {
		fmt.Printf("STAGE %02d: %s\n", i+1, stage.Name)
		fmt.Printf(" - FORMAT:\t\t%s\n", stage.Format)
		fmt.Printf(" - SLOTS:\t\t%02d\n", stage.Slots)
		fmt.Printf(" - IMPORTANCE:\t\t%02d\n", stage.Importance)

		if stage.InheritFromIndex != -1 {
			fmt.Printf(" - INHERIT FROM STAGE:\t%02d\n", stage.InheritFromIndex+1)
		}

		fmt.Print(" - TWO LEGGED:\t\t")
		if stage.TwoLegged {
			fmt.Print("true\n")
		} else {
			fmt.Print("false\n")
		}
		fmt.Println(" - GROUPS:")
		fmt.Println(stage.Groups)

		fmt.Print("\n\n")
	}
}

func (tournament *TournamentJson) SimulateWithSave(saveJson string) {
	rankings := getSave(saveJson)

	fmt.Printf("\n%s\n", tournament.Name)
	fmt.Print("-------------------------\n\n")

	var groupStates []GroupState

	scanner := bufio.NewScanner(os.Stdin)

	for k, stage := range tournament.Stages {
		fmt.Printf("STAGE %02d: %s\n\n", k+1, stage.Name)

		switch stage.Format {
		case "group":
			fmt.Printf(
				" - SLOTS: %d (%d per group, %d wildcards)\n\n",
				stage.Slots,
				stage.Slots/len(stage.Groups),
				stage.Slots%len(stage.Groups),
			)

			// Re-initialize the group state manager with the new groups
			// In case the current stage does not inherit outputs from any previous stage
			if stage.InheritFromIndex == -1 && len(stage.Groups) > 0 {
				groupStates = []GroupState{}

				for _, group := range stage.Groups {
					gState := GroupState{}
					gState.Teams = map[string]GroupStateTeam{}
					for _, team := range group {
						gState.Teams[team] = GroupStateTeam{W: 0, D: 0, L: 0, GF: 0, GA: 0, GD: 0}
					}
					groupStates = append(groupStates, gState)
				}
			}

			// Generate round robin schedule from the given groups
			scheduleList := [][][][]string{}
			for _, group := range stage.Groups {
				schedule := roundRobin(group)
				scheduleList = append(scheduleList, schedule)
			}

			for i := 0; i < len(scheduleList[0]); i++ {
				for gIdx, group := range scheduleList {
					fmt.Printf("Matchday %d -- Group %d\n", i+1, gIdx+1)
					for _, matchup := range group[i] {
						fmt.Printf("|%s vs. %s| ==> ", matchup[0], matchup[1])

						var homeTeam string = matchup[0]
						var awayTeam string = matchup[1]

						// Get user inputted result
						if scanner.Scan() {
							for {
								scores := strings.Split(scanner.Text(), "-")
								scoreCount := len(scores)
								if scoreCount == 2 {
									homeResult, awayResult := readResult(scores)
									// Extract temporary copies of states
									tempHome := groupStates[gIdx].Teams[homeTeam]
									tempAway := groupStates[gIdx].Teams[awayTeam]

									savedTeamHome, _ := rankings.getTeam(matchup[0])
									savedTeamAway, _ := rankings.getTeam(matchup[1])

									if homeResult < awayResult {
										// Case 1: Home team lost
										tempHome.L++
										tempAway.W++
										tempAway.Points += 3
									} else if homeResult > awayResult {
										// Case 2: Home team wins
										tempAway.L++
										tempHome.W++
										tempHome.Points += 3
									} else {
										// Case 3: Draw
										tempAway.D++
										tempHome.D++
										tempAway.Points++
										tempHome.Points++
									}

									tempHome.GF += homeResult
									tempHome.GA += awayResult
									tempHome.GD = (tempHome.GF - tempHome.GA)

									tempAway.GF += awayResult
									tempAway.GA += homeResult
									tempAway.GD = (tempAway.GF - tempAway.GA)

									// Re-assigns team state copies
									groupStates[gIdx].Teams[homeTeam] = tempHome
									groupStates[gIdx].Teams[awayTeam] = tempAway

									homeWeight, awayWeight := getResultWeights(homeResult, awayResult, "")
									fmt.Printf("%d\n", stage.Importance)

									homePoints := int(calculateResult(savedTeamHome.Points, savedTeamAway.Points, stage.Importance, homeWeight))
									awayPoints := int(calculateResult(savedTeamAway.Points, savedTeamHome.Points, stage.Importance, awayWeight))

									fmt.Printf("\n%s: %d (%d)\n",
										savedTeamHome.FifaCode,
										homePoints,
										homePoints-savedTeamHome.Points,
									)
									fmt.Printf("%s: %d (%d)\n\n",
										savedTeamAway.FifaCode,
										awayPoints,
										awayPoints-savedTeamAway.Points,
									)

									savedTeamHome.Points = homePoints
									savedTeamAway.Points = awayPoints

									rankings.updateTeam(savedTeamHome)
									rankings.updateTeam(savedTeamAway)

									break
								}
							}
						}
					}
					fmt.Println()
				}
				fmt.Println()

				fmt.Printf("Matchday %d Results ----\n\n", i+1)

				for gIdx, group := range groupStates {
					keysSorted := getSortedKeysFromMap(group.Teams)

					fmt.Printf("---- GROUP %d ----\n", gIdx+1)

					fmt.Println("Team\tW\tL\tD\tGF\tGA\tGD\tPoints")

					for _, team := range keysSorted {
						fmt.Printf("%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
							team,
							group.Teams[team].W,
							group.Teams[team].D,
							group.Teams[team].L,
							group.Teams[team].GF,
							group.Teams[team].GA,
							group.Teams[team].GD,
							group.Teams[team].Points,
						)
					}
				}
				fmt.Println()
			}

		case "knockout":
			return
		}

	}
}

func (tournament *TournamentJson) Simulate(save RankingTime) RankingTime {
	return save
}

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

func getSortedKeysFromMap(input map[string]GroupStateTeam) []string {
	output := []string{}

	for key := range input {
		output = append(output, key)
	}

	sort.SliceStable(output, func(i, j int) bool {
		if input[output[i]].Points != input[output[j]].Points {
			return input[output[i]].Points > input[output[j]].Points
		}

		if input[output[i]].GD != input[output[j]].GD {
			return input[output[i]].GD > input[output[j]].GD
		}

		return input[output[i]].GF > input[output[j]].GF

	})

	return output
}

func readResult(scores []string) (int, int) {

	homeResult, err := strconv.Atoi(strings.TrimSpace(scores[0]))
	if err != nil {
		log.Println("'homeScore' cannot be retrieved!")
		log.Fatal(err)
	}

	awayResult, err := strconv.Atoi(strings.TrimSpace(scores[1]))
	if err != nil {
		log.Println("'awayScore' cannot be retrieved!")
		log.Fatal(err)
	}
	return homeResult, awayResult
}
