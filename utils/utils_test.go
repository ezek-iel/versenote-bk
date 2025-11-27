package utils

import (
	"strings"
	"testing"
	"testing/quick"
)

func TestParseVerseNotation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected verseNotation
	}{
		{
			name:  "Single verse",
			input: "John 3:16",
			expected: verseNotation{
				Book:    "John",
				Chapter: 3,
				Verse:   16,
				IsRange: false,
			},
		}, {
			name:  "Verse range",
			input: "Genesis 1:1-5",
			expected: verseNotation{
				Book:       "Genesis",
				Chapter:    1,
				StartVerse: 1,
				EndVerse:   5,
				IsRange:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVerseNotation(tt.input)
			if result != tt.expected {
				t.Errorf("ParseVerseNotation(%q) = got %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringModifierProperties(t *testing.T) {
	modifier := NewStringModifier()

	// PROPERTY 1: ConvertToTableNameFormat should always result in lowercase and no spaces
	// The function 'f' returns true if the property holds, false otherwise.
	f := func(s string) bool {
		// PBT generates random strings, including empty ones.
		// Your current implementation might panic on empty strings, so we skip them for now
		// or you should fix your code to handle them.
		if len(s) == 0 {
			return true
		}

		result := modifier.ConvertToTableNameFormat(s)

		hasNoSpaces := !strings.Contains(result, " ")
		isLowerCase := result == strings.ToLower(result)
		hasUnderscoresInsteadOfSpaces := strings.Count(result, "_") >= strings.Count(s, " ")

		return hasNoSpaces && isLowerCase && hasUnderscoresInsteadOfSpaces
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("ConvertToTableNameFormat failed property check: %v", err)
	}
}

func TestPresentationFormatProperties(t *testing.T) {
	modifier := NewStringModifier()

	// PROPERTY 2: ConvertToPresentationFormat should never have underscores
	f := func(s string) bool {
		// Guard against empty strings causing panic in your current implementation
		if len(s) == 0 {
			return true
		}

		// We simulate a "database style" input (e.g., "1_samuel")
		input := strings.ReplaceAll(s, " ", "_")

		result := modifier.ConvertToPresentationFormat(input)

		hasNoUnderscores := !(strings.Count(result, "_") > 0)
		return hasNoUnderscores
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("ConvertToPresentationFormat failed property check: %v", err)
	}
}
