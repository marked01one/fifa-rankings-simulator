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
)

type ByPoint []Team

func (a ByPoint) Len() int           { return len(a) }
func (a ByPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPoint) Less(i, j int) bool { return a[i].Points > a[j].Points }

func simulate(saveFile string) {
	bytes, err := os.ReadFile("./saves/" + saveFile)
	if err != nil {
		log.Fatal(err)
	}
	var rankings RankingTime
	err = json.Unmarshal(bytes, &rankings)
	if err != nil {
		log.Fatal(err)
	}

	var homeName string
	var awayName string
	var home Team
	var away Team
	var homeId int
	var awayId int
	var importance int
	var result string
	var penalties string
	var isStopped string = ""

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("\nEnter home team: ")
		if scanner.Scan() {
			homeName = scanner.Text()
		}

		fmt.Print("Enter away team: ")
		if scanner.Scan() {
			awayName = scanner.Text()
		}

		for i, team := range rankings.Teams {
			switch team.Name {
			case homeName:
				home = team
				homeId = i
			case awayName:
				away = team
				awayId = i
			}
		}

		if home.Name != homeName {
			log.Println("Error: Country of '" + homeName + "' does not exist!")
			break
		}
		if away.Name != awayName {
			log.Println("Error: Country of '" + awayName + "' does not exist!")
			break
		}

		fmt.Print("Match importance: ")
		fmt.Scanln(&importance)

		fmt.Print("Match result [{home}-{away}]: ")
		fmt.Scanln(&result)

		if len(result) != 3 || result[1] != '-' {
			log.Println("Result format of '" + result + "' is incorrect!")
			break
		}

		fmt.Print("Penalties? [0 if home wins, 1 if away wins, skip of no penalties] ")
		fmt.Scanln(&penalties)

		isKnockout := "n"
		fmt.Print("Knockout? [y/N] ")
		fmt.Scanln(&isKnockout)

		homeResult, awayResult := getResultWeights(result[0], result[2], penalties)

		if (isKnockout) != "" {
			if homeResult == 0 {
				away.Points = int(calculateResult(away.Points, home.Points, importance, awayResult))
				log.Printf("%d\n", away.Points)
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

		fmt.Print("Continue? (type anything to stop) ")
		fmt.Scanln(&isStopped)
		if isStopped != "" {
			break
		}
	}
	fmt.Println("\nReceived stop signal. Saving results...")
	sort.Sort(ByPoint(rankings.Teams))
	file, err := json.MarshalIndent(rankings, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./saves/"+saveFile, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getRanking(team string, saveFile string) {
	bytes, err := os.ReadFile("./saves/" + saveFile)
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

func calculateResult(home, away, importance int, result float64) float64 {
	ratingDiff := (float64(home) - float64(away)) / 600 * -1
	expected := 1 / (math.Pow(10, ratingDiff) + 1)
	return math.Round(float64(home) + float64(importance)*(result-expected))
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
