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

	file := File{os.Args[0], nil}

	if file.Verify() != nil {
		log.Fatal("File not of .ffs type!")
		return
	}

}

type File struct {
	name string
	err  error
}

func (f *File) Verify() error {
	format := strings.Split(string(f.name), ".")[len(f.name)-1]
	if format != "ffs" {
		return &File{}
	} else {
		return nil
	}
}

func (f *File) Error() string { return "An error occured!" }
