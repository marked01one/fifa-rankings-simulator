package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var SOCCER_COUNTRIES []string = []string{"United States", "Canada", "Australia"}
var CreateTimestampTable string = `
CREATE TABLE IF NOT EXISTS Timestamp (
	FifaCode VARCHAR(4),
	Name VARCHAR(64),
	Point INTEGER,
	Datum VARCHAR(15),
);`

func getTimestampData(datum string) []string {

	missing := make([]string, 0, 211)
	tableNum := 0
	CodeNameMap := make(map[string]string)
	MissingMap := make(map[string]int)
	timestamp := RankingTime{}
	teams := make([]SavedTeam, 0, 211)

	currentName := ""
	currentCode := ""

	log.Println("Scraping en.wikipedia.org for FIFA country codes...")
	iconsCollector := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"))

	iconsCollector.OnHTML(`table`, func(e *colly.HTMLElement) {
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
					currentName = strings.TrimSpace(ehd.ChildText("a"))
				} else if j == 1 {
					currentCode = strings.TrimSpace(ehd.Text)
				}
			})

			CodeNameMap[currentName] = currentCode
			MissingMap[currentName] = 1
		})

	})

	if err := iconsCollector.Visit("https://en.wikipedia.org/wiki/List_of_FIFA_country_codes"); err != nil {
		log.Fatal(err)
	}

	log.Println("Attempt to scrape www.transfermarkt.com")
	teamsCollector := colly.NewCollector(
		colly.AllowedDomains("www.transfermarkt.com"),
	)

	teamsCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
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
					rawName := strings.TrimSpace(els.ChildAttr("a", "title"))

					switch rawName {
					case "Türkiye":
						rawName = "Turkey"
					case "Turkiye":
						rawName = "Turkey"
					case "Bosnia-Herzegovina":
						rawName = "Bosnia and Herzegovina"
					case "Democratic Republic of the Congo":
						rawName = "DR Congo"
					case "Republic of the Congo":
						rawName = "Congo"
					case "The Gambia":
						rawName = "Gambia"
					case "Brunei Darussalam":
						rawName = "Brunei"
					case "Timor-Leste":
						rawName = "East Timor"
					case "United States Virgin Islands":
						rawName = "U.S. Virgin Islands"
					}

					team.Name = rawName

				case confederationId:
					team.Confederation = strings.TrimSpace(els.Text)
				case pointId:
					points, err := strconv.Atoi(els.Text)
					if err != nil {
						log.Fatal(err)
					} else {
						team.Points = points
					}
				}
			})

			MissingMap[team.Name] = 0

			team.FifaCode = CodeNameMap[team.Name]
			teams = append(teams, team)
		})

	})

	if datum == "" {
		datum = "2018-12-20"
	}

	for i := 1; i < 10; i++ {
		err := teamsCollector.Visit("https://www.transfermarkt.com/statistik/weltrangliste?datum=" + datum + "&page=" + fmt.Sprint(i))
		if err != nil {
			log.Fatal(err)
		}
	}

	for key, value := range MissingMap {
		if value == 1 {
			missing = append(missing, key)
		}
	}

	timestamp.Timestamp = datum
	timestamp.Teams = teams
	timestamp.Missing = missing

	file, err := json.MarshalIndent(timestamp, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("../timestamps/"+datum+".json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return missing
}

func getAllTimestamps(cutoff string) {
	layout := "Jan 2, 2006"

	cutoffDate, err := time.Parse(layout, cutoff)
	if err != nil {
		log.Fatalf("Unable to log cutoff '%s' with layout '%s'\n", cutoff, layout)
	}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.transfermarkt.com"),
	)

	collector.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s\n", r.URL)
	})

	stamps := []string{}

	collector.OnHTML(`div[class="inline-select"]`, func(e *colly.HTMLElement) {
		e.ForEach("option", func(i int, li *colly.HTMLElement) {

			timestamp, err := time.Parse("2006-01-02", li.Attr("value"))

			if err != nil {
				log.Fatalf("Unable to parse timestamp '%s'\n", li.Attr("value"))
			}

			if timestamp.Before(cutoffDate) {
				return
			}

			timeStr := timestamp.Format("2006-01-02")

			stamps = append(stamps, timeStr)

		})
	})
	log.Println("Attempting to scrape for all timestamps")
	err = collector.Visit("https://www.transfermarkt.com/statistik/weltrangliste")
	if err != nil {
		log.Fatal(err)
	}

	for _, stamp := range stamps {

		if _, err := os.Stat("../timestamps/" + stamp + ".json"); err == nil {
			log.Printf("Skipping timestamp %s...", stamp)
			continue
		}

		missing := getTimestampData(stamp)

		if len(missing) > 0 {
			out := "Missing countries for " + stamp + ":"

			for idx, m := range missing {
				if idx == len(missing)-1 {
					out += " " + m
				} else {
					out += " " + m + ","
				}
			}

			log.Println(out)
			log.Println()
		}
	}
}
