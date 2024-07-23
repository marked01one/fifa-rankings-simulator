package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func createSave(saveTimestamp string) {
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
