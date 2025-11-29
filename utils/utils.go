package utils

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"
	"errors"
)

var bibleBookMap map[int]string = map[int]string{
	1:  "Genesis",
	2:  "Exodus",
	3:  "Leviticus",
	4:  "Numbers",
	5:  "Deuteronomy",
	6:  "Joshua",
	7:  "Judges",
	8:  "Ruth",
	9:  "1 Samuel",
	10: "2 Samuel",
	11: "1 Kings",
	12: "2 Kings",
	13: "1 Chronicles",
	14: "2 Chronicles",
	15: "Ezra",
	16: "Nehemiah",
	17: "Esther",
	18: "Job",
	19: "Psalms",
	20: "Proverbs",
	21: "Ecclesiastes",
	22: "Song of Solomon",
	23: "Isaiah",
	24: "Jeremiah",
	25: "Lamentations",
	26: "Ezekiel",
	27: "Daniel",
	28: "Hosea",
	29: "Joel",
	30: "Amos",
	31: "Obadiah",
	32: "Jonah",
	33: "Micah",
	34: "Nahum",
	35: "Habakkuk",
	36: "Zephaniah",
	37: "Haggai",
	38: "Zechariah",
	39: "Malachi",
	40: "Matthew",
	41: "Mark",
	42: "Luke",
	43: "John",
	44: "Acts",
	45: "Romans",
	46: "1 Corinthians",
	47: "2 Corinthians",
	48: "Galatians",
	49: "Ephesians",
	50: "Philippians",
	51: "Colossians",
	52: "1 Thessalonians",
	53: "2 Thessalonians",
	54: "1 Timothy",
	55: "2 Timothy",
	56: "Titus",
	57: "Philemon",
	58: "Hebrews",
	59: "James",
	60: "1 Peter",
	61: "2 Peter",
	62: "1 John",
	63: "2 John",
	64: "3 John",
	65: "Jude",
	66: "Revelation",
}

const bibleVerseRegexPattern string = `(\w{3,})\s(\d{1,}):(\d{1,})-?(\d{1,})?`

func LogQuery(query string, args ...any) {
	log.Printf("Executing query: %s with args: %v", query, args)
}

var ErrVerseNotFound = errors.New("verse not found")

func CreateDatabaseConnection(filename string) (*sql.DB, error) {
	return sql.Open("sqlite", filename)
}

func GetBibleBookName(book int) string {
	return bibleBookMap[book]
}

type StringModifier struct{}

func NewStringModifier() *StringModifier {
	return &StringModifier{}
}

// "1 Samuel" "1_samuel"
func (s *StringModifier) ConvertToTableNameFormat(book string) string {
	return strings.ToLower(strings.ReplaceAll(book, " ", "_"))
}

func (s *StringModifier) ConvertToPresentationFormat(book string) string {
	bk := strings.ReplaceAll(book, "_", " ")

	capitalize := func(word string) string {
		return strings.ToUpper(string(word[0])) + string(word[1:])
	}

	result := strings.Builder{}
	for word := range strings.SplitSeq(bk, " ") {
		result.WriteString(capitalize(word))

	}

	return result.String()
}

type verseNotation struct {
	Book       string
	Chapter    int
	Verse      int
	IsRange    bool
	StartVerse int
	EndVerse   int
}

func ParseVerseNotation(verse string) verseNotation {
	bibleVerseRegex := regexp.MustCompile(bibleVerseRegexPattern)
	matches := bibleVerseRegex.FindStringSubmatch(verse)
	log.Printf("Regex matches: %v, match length: %d", matches, len(matches))
	return CreateVerseNotationFromStringArray(matches)
}

func CreateVerseNotationFromStringArray(verseParts []string) verseNotation {
	book := verseParts[1]
	chapter, _ := strconv.Atoi(verseParts[2])
	verseNum, _ := strconv.Atoi(verseParts[3])

	if len(verseParts) == 5 {

		if verseParts[4] == "" {
			return verseNotation{
				Book:    book,
				Chapter: chapter,
				Verse:   verseNum,
				IsRange: false,
			}
		}

		startVerse, _ := strconv.Atoi(verseParts[3])
		endVerse, _ := strconv.Atoi(verseParts[4])
		return verseNotation{
			Book:       book,
			Chapter:    chapter,
			IsRange:    true,
			StartVerse: startVerse,
			EndVerse:   endVerse,
		}
	}
	return verseNotation{}
}
