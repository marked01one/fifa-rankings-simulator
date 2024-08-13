package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	funcPtr := flag.String("f", "play", "determine the type of function for the console app")
	timestampPtr := flag.String("time", "2018-12-20", "determine the timestamp to extract FIFA rankings score")
	savePtr := flag.String("save", "0001", "determine the save to use for doing simulations")
	getRankTeamPtr := flag.String("t", "Vietnam", "get the current ranking of the given team")
	confederationPtr := flag.String("conf", "", "get a query but only limited to the given confederation")

	flag.Parse()

	switch *funcPtr {
	case "create-db":
		createDatabase()
	case "timestamp":
		getTimestamp(*timestampPtr)
	case "create-save":
		createSave("../timestamps/" + *timestampPtr)
	case "play":
		saveJson := "save-" + *savePtr + ".json"
		files, err := os.ReadDir("./saves")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if saveJson == f.Name() {
				simulate(saveJson)
			}
		}
	case "get-rank":
		saveJson := "save-" + *savePtr + ".json"
		files, err := os.ReadDir("./saves")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if saveJson == f.Name() {
				getRanking(*getRankTeamPtr, saveJson)
			}
		}
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
