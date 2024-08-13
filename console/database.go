package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var CreateSaveTable string = `
CREATE TABLE IF NOT EXISTS Save (
	saveId VARCHAR(63) PRIMARY KEY,
	sourceDate DATE,
	name VARCHAR(255)
);`

var CreateSavedTeamTable string = `
CREATE TABLE IF NOT EXISTS SavedTeam (
	fifaCode VARCHAR(4) PRIMARY KEY,
	name VARCHAR(63),
	points INTEGER,
	flagImage TEXT
);`

var CreateMatchTable string = `
CREATE TABLE IF NOT EXISTS Match (
	saveId VARCHAR(63),
	matchId INTEGER,
	homeTeam VARCHAR(4),
	awayTeam VARCHAR(4),
	homeScore INTEGER,
	awayScore INTEGER,
	significance INTEGER,
	isKnockout BOOLEAN,
	homePenalties INTEGER,
	awayPenalties INTEGER,

	FOREIGN KEY (saveId) REFERENCES Save(saveId),
	FOREIGN KEY (homeTeam) REFERENCES SavedTeam(fifaCode) ON UPDATE CASCADE,
	FOREIGN KEY (awayTeam) REFERENCES SavedTeam(fifaCode) ON UPDATE CASCADE
);`

var CreateTeamTable string = `
CREATE TABLE IF NOT EXISTS Team (
	fifaCode VARCHAR(4) PRIMARY KEY,
	name VARCHAR(63)
);`

func createDatabase() {
	db, err := sql.Open("sqlite3", "./fifa.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(CreateSaveTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CreateSavedTeamTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CreateMatchTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CreateTeamTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables created successfully!")
}
