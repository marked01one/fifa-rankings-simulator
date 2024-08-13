package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var timestamp RankingTime

func insertTeamsFromJson(saveJson string) (int, error) {

	fifaCodeCollector := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org/wiki"))

	fifaCodeCollector.OnHTML("tr", func(e *colly.HTMLElement) {

	})

	db, err := sql.Open("sqlite3", "./fifa.db")
	if err != nil {
		return -1, err
	}

	data, err := os.ReadFile(saveJson)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(data, &timestamp)
	if err != nil {
		return -1, err
	}

	execStr := "INSERT INTO Team (fifaCode, name) VALUES \n"
	for i, team := range timestamp.Teams {

		underscoredTeam := strings.ReplaceAll(team.Name, " ", "_")
		footballString := "_national_football_team"

		if team.Name == soccerCountries[i] {
			footballString = "_men's_national_soccer_team"
			break
		}

		fifaCodeCollector.Visit("https://en.wikipedia.org/wiki/" + underscoredTeam + footballString)

		execStr += fmt.Sprintf("(%s,%s)", team.FifaCode, team.Name)
		if i == len(timestamp.Teams)-1 {
			execStr += ";\n"
		} else {
			execStr += ",\n"
		}
	}

	fmt.Print(execStr)

	result, err := db.Exec(execStr)
	if err != nil {
		return -1, err
	}

	output, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(output), nil
}
