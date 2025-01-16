package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Arguments struct {
	Function          *string
	Timestamp         *string
	Save              *int
	ListTeam          *string
	ListConfederation *string
	TimestampCutoff   *string
	NoSave            *bool
}

func main() {

	args := Arguments{
		Function:          flag.String("f", "play", "determine the type of function for the console app"),
		Timestamp:         flag.String("time", "2018-12-20", "determine the timestamp to extract FIFA rankings score"),
		Save:              flag.Int("save", -1, "determine the save to use for doing simulations"),
		ListTeam:          flag.String("t", "Vietnam", "get the current ranking of the given team"),
		ListConfederation: flag.String("conf", "", "get a query but only limited to the given confederation"),
		TimestampCutoff:   flag.String("cutoff", "", "Provide the cutoff date to be used for extracting timestamps"),
		NoSave:            flag.Bool("nosave", false, "Will not write the current save when the program exits"),
	}

	flag.Parse()

	formattedSave := fmt.Sprintf("%04d", *args.Save)
	saveJson := "./saves/save-" + formattedSave + ".json"

	switch *args.Function {
	case "create-db":
		createDatabase()

	case "insert-db":
		if *args.Save == -1 {
			log.Fatal("No save chosen! Choose a save insert into the DB!")
		}

		if _, err := os.Stat(saveJson); err != nil {
			log.Fatalf("Save file 'save-%s' does not exists!", formattedSave)
		}

		_, err := insertTeamsFromJson(saveJson)
		if err != nil {
			log.Fatal(err)
		}

	case "insert-teams":
		if _, err := insertTeamsFromJson("../timestamps/2024-07-18.json"); err != nil {
			log.Fatal(err)
		}

	case "timestamp":
		if *args.TimestampCutoff == "" {
			getTimestampData(*args.Timestamp)
		} else {
			getAllTimestamps(*args.TimestampCutoff)
		}

	case "create-save":
		createSave("../timestamps/" + *args.Timestamp)

	case "play":
		if *args.Save != -1 {
			fmt.Printf("Extracting save %s\n...", formattedSave)

			if _, err := os.Stat(saveJson); err != nil {
				log.Fatal(err)
			}
			simulate(saveJson, *args.NoSave)
		} else {
			fmt.Printf("No saves chosen. Creating new save from %s...\n", *args.Timestamp)
			destFile := createSave("../timestamps/" + *args.Timestamp)
			simulate(destFile, *args.NoSave)

			if *args.NoSave {
				os.Remove(destFile)
			}
		}

	case "get-rank":
		if _, err := os.Stat(saveJson); err != nil {
			log.Fatalf("Save file 'save-%s' does not exists!", formattedSave)
		}
		getRanking(*args.ListTeam, saveJson)

	case "get-list":
		if _, err := os.Stat(saveJson); err != nil {
			log.Fatalf("Save file 'save-%s' does not exists!", formattedSave)
		}
		getSortedRankings(*args.ListConfederation, saveJson)

	case "parse-tournament":
		testPath := "../concepts/tournament_templates/aff_cup.json"
		tournament := TournamentJson{}
		tournament.GetJson(testPath)
		tournament.SimulateWithSave(saveJson)

	}
}
