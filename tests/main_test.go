package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/ezek-iel/bible-notes-backend/converter"
	"github.com/ezek-iel/bible-notes-backend/query"
	"github.com/ezek-iel/bible-notes-backend/utils"
)

func compareVerses(got, expected []converter.Verse, t testing.TB) {
	giveError := func(field string, v1, v2 converter.Verse) {
		t.Errorf("%v is conflicting in \n%v and \n%v", field, v1, v2)
	}

	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("The verse list got and expected are different, %d, %d", len(got), len(expected))
	}

	for n, _ := range got {
		exp := expected[n]
		gt := got[n]

		if exp.Book != gt.Book {
			giveError("Book", exp, gt)
		}

		if exp.Chapter != gt.Chapter {
			giveError("Chapter", exp, gt)
		}

		if exp.Verse != gt.Verse {
			giveError("Verse", exp, gt)
		}

		if strings.TrimSpace(exp.Text) != strings.TrimSpace(gt.Text) {
			giveError("Text", exp, gt)
		}

		if exp.Comment != gt.Comment {
			giveError("Comment", exp, gt)
		}
	}
}

func TestQuerySingleBibleVerse(t *testing.T) {
	bibleDB := &query.BibleDatabase{DBPath: "test.db"}

	cases := []struct {
		verse         string
		errorExpected error
		result        []converter.Verse
	}{
		{"John 1:15",
			nil,
			[]converter.Verse{{
				Book:    0,
				Chapter: 1,
				Verse:   15,
				Text:    `John testified about him when he shouted to the crowds, “This is the one I was talking about when I said, ‘Someone is coming after me who is far greater than I am, for he existed long before me.’”`,
				Comment: ""}}},
		{
			"John 1:52",
			utils.ErrVerseNotFound,
			[]converter.Verse{{}},
		},
		{
			"Revelation 2:0",
			utils.ErrVerseNotFound,
			[]converter.Verse{{}},
		},
		{
			"Genesis 5:18-18",
			utils.ErrVerseNotFound,
			[]converter.Verse{{}},
		},
		{
			"Ezekiel 18:4-8",
			nil,
			[]converter.Verse{
				{Book: 0, Chapter: 18, Verse: 4, Text: `For all people are mine to judge — both parents and children alike. And this is my rule: The person who sins is the one who will die.`, Comment: ""},
				{Book: 0, Chapter: 18, Verse: 5, Text: `“Suppose a certain man is righteous and does what is just and right.`, Comment: ""},
				{Book: 0, Chapter: 18, Verse: 6, Text: `He does not feast in the mountains before Israel’s idols or worship them. He does not commit adultery or have intercourse with a woman during her menstrual period.`, Comment: `The Hebrew term (literally <i>round things</i>) probably alludes to dung; also in <a href='/NLT/26/18/12'>18:12</a>, <a href='/NLT/26/18/15'>15</a>.`},
				{Book: 0, Chapter: 18, Verse: 7, Text: `He is a merciful creditor, not keeping the items given as security by poor debtors. He does not rob the poor but instead gives food to the hungry and provides clothes for the needy.`},
				{Book: 0, Chapter: 18, Verse: 8, Text: `He grants loans without interest, stays away from injustice, is honest and fair when judging others,`},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.verse, func(t *testing.T) {
			result, err := bibleDB.Query(tc.verse)

			if !errors.Is(err, tc.errorExpected) {
				t.Fatalf("Expected error %v, got %v instead", tc.errorExpected, err)
			}

			compareVerses(result, tc.result, t)
		})
	}
}
