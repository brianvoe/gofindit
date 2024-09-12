package tokenizers

import "testing"

// build simple tokenizer for tests
type SimpleTokenizer struct {
	// store the value
	val string
}

// Process will take in a string value and
// use it to fill out the struct fields
func (t *SimpleTokenizer) Process(val string) error {
	t.val = val
	return nil
}

// ToSearch will take in a string value and
// return a slice of strings that can be used
// to search
func (t *SimpleTokenizer) ToSearch(val string) ([]string, error) {
	return []string{val}, nil
}

// Search will take in a slice of strings and
// return true if the value is found
func (t *SimpleTokenizer) Search(val []string) (bool, error) {
	for _, v := range val {
		if v == t.val {
			return true, nil
		}
	}
	return false, nil
}

func TestTokenizerSetGetDelete(t *testing.T) {
	// create a new tokenizer
	tokenizer := &SimpleTokenizer{}

	// set the tokenizer
	SetTokenizer("simple", tokenizer)

	// get the tokenizer
	tok, err := GetTokenizer("simple", nil)
	if err != nil {
		t.Errorf("error getting tokenizer: %v", err)
	}

	// check if the tokenizer is the same
	if tok != tokenizer {
		t.Errorf("expected tokenizer to be the same")
	}

	// delete the tokenizer
	DeleteTokenizer("simple")

	// get the tokenizer
	_, err = GetTokenizer("simple", nil)
	if err == nil {
		t.Errorf("expected error getting tokenizer")
	}
}
