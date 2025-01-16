package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ByPoint []SavedTeam

func (a ByPoint) Len() int           { return len(a) }
func (a ByPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPoint) Less(i, j int) bool { return a[i].Points > a[j].Points }

func simulate(saveFile string, noSave bool) {
	if noSave {
		fmt.Println("WARNING: Simulator is running in no-save mode!")
	}

	bytes, err := os.ReadFile(saveFile)
	if err != nil {
		log.Fatal(err)
	}
	var rankings RankingTime
	err = json.Unmarshal(bytes, &rankings)
	if err != nil {
		log.Fatal(err)
	}

	var importance int = 0
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var homeName string
		var awayName string
		var home SavedTeam
		var away SavedTeam
		var homeId int
		var awayId int
		var result string
		var penalties string
		var isStopped string = ""
		var isKnockout string = ""

		fmt.Print("\nTeams? (Home-Away) ")
		if scanner.Scan() {
			names := strings.Split(scanner.Text(), "-")
			if len(names) != 2 {
				fmt.Println("Error: there are not exact two names!")
				break
			}
			homeName = strings.ToLower(strings.TrimSpace(names[0]))
			awayName = strings.ToLower(strings.TrimSpace(names[1]))
		}

		for i, team := range rankings.Teams {
			if len(homeName) == 3 && len(awayName) == 3 {
				switch strings.ToLower(team.FifaCode) {
				case homeName:
					home = team
					homeId = i
				case awayName:
					away = team
					awayId = i
				}
			} else {
				switch strings.ToLower(team.Name) {
				case homeName:
					home = team
					homeId = i
				case awayName:
					away = team
					awayId = i
				}
			}
		}

		if strings.ToLower(home.Name) != homeName {
			log.Println("Error: Country of '" + homeName + "' does not exist!")
			break
		}
		if strings.ToLower(away.Name) != awayName {
			log.Println("Error: Country of '" + awayName + "' does not exist!")
			break
		}

		fmt.Print("Results? (Home-Away) ")
		fmt.Scanln(&result)

		if len(result) != 3 || result[1] != '-' {
			log.Println("Result format of '" + result + "' is incorrect!")
			break
		}

		if importance != 0 {
			fmt.Printf("Importance? (current at %d) ", importance)
		} else {
			fmt.Print("Importance? ")
		}

		fmt.Scanln(&importance)

		if importance >= 35 {
			fmt.Print("Knockout? [y/N] ")
			fmt.Scanln(&isKnockout)
		}

		if result[0] == result[2] && strings.ToLower(isKnockout) != "n" && isKnockout != "" {
			fmt.Print("Penalties? [0 if home wins, 1 if away wins, skip of no penalties] ")
			fmt.Scanln(&penalties)
		}

		homeResult, awayResult := getResultWeights(result[0], result[2], penalties)

		if isKnockout != "" {
			if homeResult == 0 {
				away.Points = int(calculateResult(away.Points, home.Points, importance, awayResult))
			} else if awayResult == 0 {
				home.Points = int(calculateResult(home.Points, away.Points, importance, homeResult))
			} else {
				home.Points = int(calculateResult(home.Points, away.Points, importance, homeResult))
				away.Points = int(calculateResult(away.Points, home.Points, importance, awayResult))
			}
		} else {
			home.Points = int(calculateResult(home.Points, away.Points, importance, homeResult))
			away.Points = int(calculateResult(away.Points, home.Points, importance, awayResult))
		}
		fmt.Println("New total for " + home.Name + ": " + fmt.Sprint(home.Points) + " (" + fmt.Sprint(home.Points-rankings.Teams[homeId].Points) + ")")
		fmt.Println("New total for " + away.Name + ": " + fmt.Sprint(away.Points) + " (" + fmt.Sprint(away.Points-rankings.Teams[awayId].Points) + ")")

		rankings.Teams[homeId] = home
		rankings.Teams[awayId] = away

		fmt.Print("\nContinue? (type anything to stop) ")
		fmt.Scanln(&isStopped)
		if isStopped != "" {
			break
		}
	}

	fmt.Println("Received stop signal")

	if noSave {
		fmt.Println("Exiting...")
		return
	}

	fmt.Println("Saving results...")
	sort.Sort(ByPoint(rankings.Teams))
	file, err := json.MarshalIndent(rankings, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(saveFile, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getRanking(team string, saveJson string) {
	bytes, err := os.ReadFile(saveJson)
	if err != nil {
		log.Fatal(err)
	}
	var rankings RankingTime
	err = json.Unmarshal(bytes, &rankings)
	if err != nil {
		log.Fatal(err)
	}

	for i, t := range rankings.Teams {
		if t.Name == team {
			fmt.Printf("Current ranking for %s: %d\n", t.Name, i+1)
			return
		}
	}
	fmt.Printf("Could not find team with name '%s'\n", team)
}

func getSortedRankings(confederation string, saveJson string) {
	bytes, err := os.ReadFile(saveJson)
	if err != nil {
		log.Fatal(err)
	}
	var rankings RankingTime
	err = json.Unmarshal(bytes, &rankings)
	if err != nil {
		log.Fatal(err)
	}

	if confederation == "" {
		fmt.Print("World\tPoints\tName\n")
		fmt.Println("------------------------------------------------")
		for i, t := range rankings.Teams {
			fmt.Printf("%d\t%d\t%s\n", i+1, t.Points, t.Name)
		}
		return
	}

	counter := 1
	fmt.Printf("\n%s\tWorld\tPoints\tName\n", confederation)
	fmt.Println("------------------------------------------------")
	for i, t := range rankings.Teams {
		if t.Confederation == confederation {
			fmt.Printf("%d\t%d\t%d\t%s\n", counter, i+1, t.Points, t.Name)
			counter++
		}
	}

}

func calculateResult(home, away, importance int, result float64) float64 {
	ratingDiff := (float64(home) - float64(away)) / 600 * -1
	expected := 1 / (math.Pow(10, ratingDiff) + 1)
	finalResult := float64(home) + float64(importance)*(result-expected)
	return math.Round(finalResult)
}

func getResultWeights(home, away byte, penalties string) (float64, float64) {
	homeScore, _ := strconv.Atoi(string(home))
	awayScore, _ := strconv.Atoi(string(away))

	if homeScore > awayScore {
		return 1, 0
	}

	if awayScore > homeScore {
		return 0, 1
	}

	switch penalties {
	case "0":
		return 0.75, 0.5
	case "1":
		return 0.5, 0.75
	}

	return 0.5, 0.5
}
