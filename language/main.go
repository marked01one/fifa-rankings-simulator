package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 1 {
		log.Fatal("No arguments given!")
		return
	}

	file := FileName(os.Args[0])
	if !file.IsFFS() {
		log.Fatal("File not of .ffs type!")
		return
	}

}

type FileName string

func (f FileName) IsFFS() bool {
	format := strings.Split(string(f), ".")[len(f)-1]
	if format != "ffs" {
		return false
	} else {
		return true
	}
}
