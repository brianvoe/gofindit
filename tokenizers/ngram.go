package tokenizers

import (
	"fmt"
	"strings"

	"golang.org/x/text/transform"
)

func init() {
	SetTokenizer("ngram", NewNGram(3, 10))
}

type NGram struct {
	min   int
	max   int
	index map[string]bool // Index of n-grams for efficient searching
}

// NewNGram returns a new NGram tokenizer
func NewNGram(min, max int) *NGram {
	// If min is greater than max, swap them
	if min > max {
		min, max = max, min
	}

	return &NGram{
		min:   min,
		max:   max,
		index: make(map[string]bool), // Initialize the index map
	}
}

// Process takes a string value, generates n-grams, and fills out the index
func (n *NGram) Process(val string) error {
	// Clean the input string
	val = cleanNGramStr(val)

	// After clean check if the string is empty
	if len(val) < n.min {
		return fmt.Errorf("input shorter than min n-gram length")
	}

	nGrams := generateNGrams(val, n.min, n.max)

	// Index the n-grams
	for _, nGram := range nGrams {
		n.index[nGram] = true
	}

	return nil
}

// Tokenize generates all possible n-grams from the input string
func (n *NGram) ToSearch(val string) ([]string, error) {
	// Clean the input string
	val = cleanNGramStr(val)

	// After clean check if the string is empty
	if len(val) < n.min {
		return nil, fmt.Errorf("input shorter than min n-gram length")
	}

	// Check if val is larger than the max n-gram length
	if len(val) > n.max {
		val = val[:n.max]
	}

	return []string{val}, nil
}

// Search checks if all specified n-grams exist in the index
func (n *NGram) Search(vals []string) (bool, error) {
	// Expected value is a single value of the string we are searching for
	if len(vals) != 1 {
		return false, fmt.Errorf("expected a single value, got %d", len(vals))
	}

	// Check if the value is in the index
	if _, exists := n.index[vals[0]]; !exists {
		return false, nil
	}

	return true, nil
}

func cleanNGramStr(val string) string {
	// Remove accents from the string
	val, _, _ = transform.String(normalizer, val)

	// Lowercase the string
	val = strings.ToLower(val)

	return val
}

func generateNGrams(val string, min, max int) []string {
	var nGrams []string
	for i := 0; i <= len(val)-min; i++ {
		for j := min; j <= max && i+j <= len(val); j++ {
			nGrams = append(nGrams, val[i:i+j])
		}
	}
	return nGrams
}
