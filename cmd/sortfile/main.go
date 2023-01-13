package main

import (
	"log"
	"os"

	"github.com/KEINOS/go-sortfile/sortfile"
	"github.com/pkg/errors"
)

func main() {
	ExitOnError(Run())
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	if len(os.Args) < 2 {
		return errors.New("No arguments given")
	}

	inFile := os.Args[1]
	outFile := os.Args[2]

	err := sortfile.FromPath(inFile, outFile, false)

	return errors.Wrap(err, "Failed to sort file")
}
