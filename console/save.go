package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type RankingTime struct {
	Timestamp string      `json:"timestamp"`
	Missing   []string    `json:"missing"`
	Teams     []SavedTeam `json:"teams"`
}

type SavedTeam struct {
	Name          string `json:"name"`
	FifaCode      string `json:"fifaCode"`
	Confederation string `json:"confederation"`
	Points        int    `json:"points"`
}

func (ranking *RankingTime) getTeam(fifaCode string) (SavedTeam, error) {
	for _, t := range ranking.Teams {
		if fifaCode == t.FifaCode {
			return t, nil
		}
	}

	return SavedTeam{}, fmt.Errorf("no team found associated with FIFA code '%s'", fifaCode)
}

func (ranking *RankingTime) updateTeam(team SavedTeam) error {
	for idx, t := range ranking.Teams {
		if team.FifaCode == t.FifaCode {
			ranking.Teams[idx] = team
			return nil
		}
	}

	return fmt.Errorf("no team found associated with FIFA code '%s'", team.FifaCode)
}

func createSave(saveTimestamp string) string {
	saves, err := os.ReadDir("./saves")
	if err != nil {
		log.Fatal(err)
	}

	length := fmt.Sprint(len(saves) + 1)

	destName := "save-"

	if len(length) < 4 {
		for i := 0; i < 4-len(length); i++ {
			destName += "0"
		}
	}
	destName += length
	destFile := "./saves/" + destName + ".json"
	sourceFile := saveTimestamp + ".json"

	copyTimestamp(sourceFile, destFile)

	return destFile
}

func copyTimestamp(source, destination string) {
	// Open the source file
	srcFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure all data is written to save
	err = destFile.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func getSave(saveJson string) RankingTime {
	if _, err := os.Stat(saveJson); err != nil {
		log.Fatal(err)
	}

	rankings := RankingTime{}

	bytes, err := os.ReadFile(saveJson)

	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(bytes, &rankings); err != nil {
		log.Fatal(err)
	}

	return rankings
}
