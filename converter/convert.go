package converter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ezek-iel/bible-notes-backend/utils"
	_ "modernc.org/sqlite"
)

type Verse struct {
	Pk          int    `json:"pk"`
	Translation string `json:"translation"`
	Book        int    `json:"book"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	Text        string `json:"text"`
	Comment     string `json:"comment"`
}

var stringModifier = utils.NewStringModifier()

func (v *Verse) GetBookName() string {
	return utils.GetBibleBookName(v.Book)
}



func getJsonDataFromFile(filename string) ([]Verse, error) {
	verseList := []Verse{}

	fileContents, readFileErr := os.ReadFile(filename)

	if readFileErr != nil {
		return verseList, readFileErr
	}

	unmarshalErr := json.Unmarshal(fileContents, &verseList)

	if unmarshalErr != nil {
		return verseList, unmarshalErr
	}

	return verseList, nil
}

func ConvertBibleJsonToDB(filename string, destination string) error {
	db, createDatabaseConnectionError := utils.CreateDatabaseConnection(destination)

	if createDatabaseConnectionError != nil {
		return createDatabaseConnectionError
	}

	verseList, convertFromJsonErr := getJsonDataFromFile(filename)

	if convertFromJsonErr != nil {
		return convertFromJsonErr
	}

	for _, verse := range verseList {
		verseBookName := verse.GetBookName()

		bookExists, bookExistsErr := bookTableExists(db, verseBookName)

		if bookExistsErr != nil {
			return bookExistsErr
		}

		if !bookExists {
			createBook(db, verseBookName)
		}

		insertVerseErr := insertVerse(db, verse)

		if insertVerseErr != nil {
			return insertVerseErr
		}

	}

	return nil
}

func insertVerse(db *sql.DB, v Verse) error {
	verseBookName := stringModifier.ConvertToTableNameFormat(v.GetBookName())
	query := fmt.Sprintf("INSERT INTO \"%s\" (pk, chapter, verse, text, comment) VALUES (?, ?, ?, ?, ?)", verseBookName)
	log.Printf("Inserting into %s: Chapter %d, Verse %d", verseBookName, v.Chapter, v.Verse)
	_, err := db.Exec(query, v.Pk, v.Chapter, v.Verse, v.Text, v.Comment)
	return err
}

func createBook(db *sql.DB, bookname string) error {
	bn := stringModifier.ConvertToTableNameFormat(bookname)
	query := fmt.Sprintf("CREATE TABLE \"%s\" (pk int, chapter int, verse int, text text, comment text)", bn)
	log.Printf("Creating bookname %s", bn)
	_, err := db.Exec(query)
	return err
}

func bookTableExists(db *sql.DB, bookName string) (bool, error) {
	query := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?"
	var count int

	bk := stringModifier.ConvertToTableNameFormat(bookName)

	err := db.QueryRow(query, bk).Scan(&count)
	if err != nil {
		return false, err
	}

	bookExists := count > 0

	if bookExists {
		log.Printf("Book %v exists", bk)
	} else {
		log.Printf("Book %v does not exist", bk)
	}
	return bookExists, nil
}
