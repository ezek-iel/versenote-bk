package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezek-iel/bible-notes-backend/converter"
)

func checkFile(filename string, ext string) bool {
	return strings.TrimSpace(filename) != "" && filepath.Ext(filename) == ext
}

func main() {
	source := flag.String("src", "", "the source file in JSON")
	destination := flag.String("dest", "", "the destination file in .db")
	flag.Parse()

	if !(checkFile(*source, ".json") || checkFile(*destination, ".db")) {
		fmt.Println("The src file is not in json or the destination in db")
		os.Exit(1)
	}

	convertErr := converter.ConvertBibleJsonToDB(*source, *destination)

	if convertErr != nil {
		fmt.Printf("an error occured: %v\n", convertErr)
		os.Exit(1)
	}
}
