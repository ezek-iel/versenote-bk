package query

import (
	"database/sql"
	"errors"

	"github.com/ezek-iel/versenote-bk/converter"
	"github.com/ezek-iel/versenote-bk/utils"
)

type BibleDatabase struct {
	DBPath string
	DB     *sql.DB
}

func (b *BibleDatabase) Query(verseNotation string) ([]converter.Verse, error) {
	err := b.connect()
	if err != nil {
		return []converter.Verse{}, err
	}
	defer b.Close()

	verseInfo := utils.ParseVerseNotation(verseNotation)

	if verseInfo.IsRange {
		return b.queryVerseRange(verseInfo.Book, verseInfo.Chapter, verseInfo.StartVerse, verseInfo.EndVerse)
	}

	verse, queryErr := b.queryVerse(verseInfo.Book, verseInfo.Chapter, verseInfo.Verse)

	return []converter.Verse{verse}, queryErr
}

func (b *BibleDatabase) Close() error {
	return b.DB.Close()
}

func (b *BibleDatabase) connect() error {
	db, err := utils.CreateDatabaseConnection(b.DBPath)
	if err != nil {
		return err
	}
	b.DB = db
	return nil
}

func (b *BibleDatabase) queryVerse(book string, chapter int, verse int) (converter.Verse, error) {
	bibleVerseQuery := `SELECT pk, verse, text, comment FROM ` + utils.NewStringModifier().ConvertToTableNameFormat(book) + ` WHERE chapter = ? AND verse = ?;`

	utils.LogQuery(bibleVerseQuery, chapter, verse)
	row := b.DB.QueryRow(bibleVerseQuery, chapter, verse)
	var verseResult converter.Verse
	readRowError := row.Scan(&verseResult.Pk, &verseResult.Verse, &verseResult.Text, &verseResult.Comment)

	verseResult.Chapter = chapter
	if readRowError != nil {

		if errors.Is(readRowError, sql.ErrNoRows) {
			return converter.Verse{}, utils.ErrVerseNotFound
		}
		return converter.Verse{}, readRowError
	}
	return verseResult, nil
}

func (b *BibleDatabase) queryVerseRange(book string, chapter int, startVerse int, endVerse int) ([]converter.Verse, error) {

	if startVerse >= endVerse {
		return []converter.Verse{{}}, utils.ErrVerseNotFound
	}

	bibleVerseQuery := `SELECT pk, verse, text, comment FROM ` + utils.NewStringModifier().ConvertToTableNameFormat(book) + ` WHERE chapter = ? AND verse BETWEEN ? AND ?;`
	utils.LogQuery(bibleVerseQuery, chapter, startVerse, endVerse)
	rows, queryError := b.DB.Query(bibleVerseQuery, chapter, startVerse, endVerse)
	if queryError != nil {
		return []converter.Verse{}, queryError
	}
	defer rows.Close()

	var verses []converter.Verse

	rowCount := 0

	for rows.Next() {
		var verseResult converter.Verse
		readRowError := rows.Scan(&verseResult.Pk, &verseResult.Verse, &verseResult.Text, &verseResult.Comment)
		if readRowError != nil {
			return []converter.Verse{}, readRowError
		}
		verseResult.Chapter = chapter
		verses = append(verses, verseResult)
		rowCount++
	}

	if rowCount == 0 {
		return verses, utils.ErrVerseNotFound
	}
	return verses, nil
}
