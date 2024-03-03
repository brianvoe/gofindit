package gofindit

import (
	"strings"
	"unicode"
)

// Tokenize takes a string and returns an array of tokens
func Tokenize(str string, filters []FilterFunc) ([]string, error) {
	// Process the initial string
	tokens := strings.FieldsFunc(str, func(r rune) bool {
		// Split on any character that is not a letter or a number
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Loop through filters and apply them to the tokens
	var err error
	for _, filter := range filters {
		tokens, err = filter(tokens)
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}
