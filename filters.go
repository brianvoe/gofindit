package gofindit

import "strings"

type FilterFunc func([]string) ([]string, error)

var DefaultFilters = []FilterFunc{
	FilterLowercase,
}

// FilterLowercase converts all tokens to lowercase
func FilterLowercase(tokens []string) ([]string, error) {
	out := make([]string, len(tokens))
	for i, token := range tokens {
		out[i] = strings.ToLower(token)
	}
	return out, nil
}
