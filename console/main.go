package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	funcPtr := flag.String("f", "play", "determine the type of function for the console app")
	timestampPtr := flag.String("time", "2018-12-20", "determine the timestamp to extract FIFA rankings score")
	savePtr := flag.String("save", "-1", "determine the save to use for doing simulations")
	getRankTeamPtr := flag.String("t", "Vietnam", "get the current ranking of the given team")
	confederationPtr := flag.String("conf", "", "get a query but only limited to the given confederation")
	cutoffPtr := flag.String("cutoff", "", "Provide the cutoff date to be used for extracting timestamps")

	noSavePtr := flag.Bool("nosave", false, "will not write the current save when the program exits")

	flag.Parse()

	savePtrInt, _ := strconv.Atoi(*savePtr)
	formattedSave := fmt.Sprintf("%04d", savePtrInt)
	saveJson := "./saves/save-" + formattedSave + ".json"

	switch *funcPtr {
	case "create-db":
		createDatabase()

	case "insert-db":
		if *savePtr == "-1" {
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
		_, err := insertTeamsFromJson("../timestamps/2024-07-18.json")
		if err != nil {
			log.Fatal(err)
		}

	case "timestamp":
		if *cutoffPtr == "" {
			getTimestampData(*timestampPtr)
		} else {
			getAllTimestamps(*cutoffPtr)
		}

	case "create-save":
		createSave("../timestamps/" + *timestampPtr)

	case "play":
		if *savePtr != "-1" {
			fmt.Printf("Extracting save %s\n...", formattedSave)

			if _, err := os.Stat(saveJson); err != nil {
				log.Fatal(err)
			}
			simulate(saveJson, *noSavePtr)
		} else {
			fmt.Printf("No saves chosen. Creating new save from %s...\n", *timestampPtr)
			destFile := createSave("../timestamps/" + *timestampPtr)
			simulate(destFile, *noSavePtr)

			if *noSavePtr {
				os.Remove(destFile)
			}
		}

	case "get-rank":
		if _, err := os.Stat(saveJson); err != nil {
			log.Fatalf("Save file 'save-%s' does not exists!", formattedSave)
		}
		getRanking(*getRankTeamPtr, saveJson)

	case "get-list":
		saveJson := "save-" + *savePtr + ".json"
		files, err := os.ReadDir("./saves")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if saveJson == f.Name() {
				getSortedRankings(*confederationPtr, saveJson)
			}
		}
	}

}
