package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezek-iel/versenote-bk/converter"
	"github.com/ezek-iel/versenote-bk/query"
)

func checkFile(filename string, ext string) bool {
	return strings.TrimSpace(filename) != "" && filepath.Ext(filename) == ext
}

func main() {

	convertCommand := flag.NewFlagSet("convert", flag.ExitOnError)
	queryCommand := flag.NewFlagSet("query", flag.ExitOnError)

	source := convertCommand.String("source-file", "", "the source file in JSON")
	destination := convertCommand.String("destination-file", "", "the destination file in .db")

	bibleVerseNotation := queryCommand.String(
		"verse-notation",
		"",
		"the verse notation in the format Book Chapter:Verse or Book Chapter:StartVerse-EndVerse e.g John 1:14 or John 1:14-16")
	bibleDatabaseFileName := queryCommand.String("database-file", "", "the bible database file in .db")

	if len(os.Args) < 2 {
		fmt.Println("expected 'convert' or 'query' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "convert":
		convertCommand.Parse(os.Args[2:])
		if !(checkFile(*source, ".json") || checkFile(*destination, ".db")) {
			fmt.Println("The src file is not in json or the destination in db")
			convertCommand.PrintDefaults()
			os.Exit(1)
		}

		convertErr := converter.ConvertBibleJsonToDB(*source, *destination)

		if convertErr != nil {
			fmt.Printf("an error occured: %v\n", convertErr)
			os.Exit(1)
		}
	case "query":
		queryCommand.Parse(os.Args[2:])
		if !checkFile(*bibleDatabaseFileName, ".db") {
			fmt.Println("The database file is not in db format")
			queryCommand.PrintDefaults()
			os.Exit(1)
		}
		bibleDB := &query.BibleDatabase{
			DBPath: *bibleDatabaseFileName,
		}
		verses, queryErr := bibleDB.Query(*bibleVerseNotation)

		if queryErr != nil {
			fmt.Printf("an error occured: %v\n", queryErr)
			os.Exit(1)
		}

		for _, verse := range verses {
			fmt.Printf("%s %d:%d - %s\n", verse.GetBookName(), verse.Chapter, verse.Verse, verse.Text)
		}

	case "default":
		fmt.Println("expected 'convert' or 'query' subcommands")
		os.Exit(1)
	}

}
