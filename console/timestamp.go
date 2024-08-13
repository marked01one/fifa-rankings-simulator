package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var soccerCountries []string = []string{"United States", "Canada", "Australia"}

func getTimestamp(datum string) {
	log.Println("Attemt to scrape www.transfermarkt.com")
	teamsCollector := colly.NewCollector(colly.AllowedDomains("www.transfermarkt.com"))
	fifaCodeCollector := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org/wiki"))

	timestamp := RankingTime{}
	teams := make([]SavedTeam, 0, 211)

	teamsCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	teamsCollector.OnHTML(`table[class="items"]`, func(e *colly.HTMLElement) {
		var confederationId int
		var pointId int

		e.ForEach("thead th", func(i int, eh *colly.HTMLElement) {
			switch eh.ChildText("a") {
			case "Confederation":
				confederationId = i
			case "Points":
				pointId = i
			}
		})

		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
			team := SavedTeam{}
			el.ForEach("td", func(i int, els *colly.HTMLElement) {
				switch i {
				case 1:
					team.Name = els.ChildAttr("a", "title")
				case confederationId:
					team.Confederation = els.Text
				case pointId:
					points, err := strconv.Atoi(els.Text)
					if err != nil {
						log.Fatal(err)
					} else {
						team.Points = points
					}
				}
			})

			teams = append(teams, team)
		})
	})

	if datum == "" {
		datum = "2018-12-20"
	}

	for i := 1; i < 10; i++ {
		teamsCollector.Visit("https://www.transfermarkt.com/statistik/weltrangliste?datum=" + datum + "&page=" + fmt.Sprint(i))
	}
	timestamp.Timestamp = datum
	timestamp.Teams = teams

	for _, team := range timestamp.Teams {

		underscoredTeam := strings.ReplaceAll(team.Name, " ", "_")
		footballString := "_national_football_team"

		for i := 0; i < len(soccerCountries); i++ {
			if team.Name == soccerCountries[i] {
				footballString = "_men's_national_soccer_team"
				break
			}
		}

		fifaCodeCollector.Visit("https://en.wikipedia.org/wiki/" + underscoredTeam + footballString)
	}

	file, err := json.MarshalIndent(timestamp, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("../timestamps/"+datum+".json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
