package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var timestamp RankingTime

var NameMap map[string]string = map[string]string{}

func insertTeamsFromJson(saveJson string) (int, error) {
	tableNum := 0
	currentName := ""
	currentCode := ""
	// Initialize the FIFA code collector to get corresponding FIFA codes from Wikipedia
	fifaCodeCollector := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
	fifaCodeCollector.OnHTML(`table`, func(e *colly.HTMLElement) {
		tableNum++
		if tableNum > 4 {
			return
		}

		e.ForEach("tbody tr", func(i int, eh *colly.HTMLElement) {
			if i == 0 {
				return
			}

			eh.ForEach("td", func(j int, ehd *colly.HTMLElement) {
				if j == 0 {
					currentName = ehd.ChildText("a")
				} else if j == 1 {
					currentCode = strings.TrimSpace(ehd.Text)
				}
			})

			switch currentName {
			case "Turkey":
				NameMap["Türkiye"] = currentCode
			case "Bosnia and Herzegovina":
				NameMap["Bosnia-Herzegovina"] = currentCode
			case "DR Congo":
				NameMap["Democratic Republic of the Congo"] = currentCode
			case "Congo":
				NameMap["Republic of the Congo"] = currentCode
			case "Gambia":
				NameMap["The Gambia"] = currentCode
			case "Brunei":
				NameMap["Brunei Darussalam"] = currentCode
			case "East Timor":
				NameMap["Timor-Leste"] = currentCode
			case "U.S. Virgin Islands":
				NameMap["United States Virgin Islands"] = currentCode
			default:
				NameMap[currentName] = currentCode
			}

		})

	})

	err := fifaCodeCollector.Visit("https://en.wikipedia.org/wiki/List_of_FIFA_country_codes")

	if err != nil {
		log.Fatal(err)
	}

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

	TeamExec := "INSERT OR REPLACE INTO Team (fifaCode, name) VALUES \n"
	for i, team := range timestamp.Teams {
		TeamExec += fmt.Sprintf("('%s','%s')", NameMap[team.Name], team.Name)

		if i == len(timestamp.Teams)-1 {
			TeamExec += ";\n"
		} else {
			TeamExec += ",\n"
		}
	}

	fmt.Print()

	result, err := db.Exec(TeamExec)
	if err != nil {
		return -1, err
	}

	output, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(output), nil
}
