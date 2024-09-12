package tokenizers

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
)

func init() {
	SetTokenizer("words", NewWords())
}

type Words struct {
	words []string
}

// NewWords will remove accents, lowercase and
// make it searchable via words
func NewWords() *Words {
	return &Words{}
}

// Process will take in an any value and
// use it to fill out the struct fields
func (w *Words) Process(str string) error {
	// Set the words
	var err error
	w.words, err = w.ToSearch(str)
	if err != nil {
		return err
	}

	return nil
}

// ToSearchBytes will return the bytes to search
func (w *Words) ToSearch(str string) ([]string, error) {
	if str == "" {
		return nil, errors.New("empty string")
	}

	var err error

	// Remove accents from the string
	str, _, err = transform.String(normalizer, str)
	if err != nil {
		return nil, err
	}

	// Lowercase the string
	str = strings.ToLower(str)

	// Remove apostrophes from the string
	str = strings.ReplaceAll(str, "'", "")

	// Tokenize the string
	tokens := strings.FieldsFunc(str, func(r rune) bool {
		// Split on any character that is not a letter, a number, or an apostrophe.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return tokens, nil
}

// Search checks if the provided words appear in order within the Words struct's words slice.
func (w *Words) Search(val []string) (bool, error) {
	// Start index for search
	searchIndex := 0

	for _, searchWord := range val {
		found := false
		for i := searchIndex; i < len(w.words); i++ {
			if w.words[i] == searchWord {
				// If the word is found, update searchIndex to start from the next word
				searchIndex = i + 1
				found = true
				break // Break the inner loop and continue with the next searchWord
			}
		}
		// If any of the words is not found, return false
		if !found {
			return false, nil
		}
	}

	// If all words were found in order, return true
	return true, nil
}
