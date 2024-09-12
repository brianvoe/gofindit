package tokenizers

import (
	"fmt"
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

type Tokenizer interface {
	// Process will take in a string value and
	// use it to fill out the struct fields
	Process(val string) error

	// ToSearch will take in a string value and
	// return a slice of strings that can be used
	// to search
	ToSearch(val string) ([]string, error)

	Search(val []string) (bool, error)
}

type storage struct {
	tokenizers map[string]Tokenizer

	// lock
	mu sync.RWMutex
}

// store is the storage for all tokenizers
var store = &storage{
	tokenizers: make(map[string]Tokenizer),
}

// GetTokenizer returns a tokenizer from the store
// and initializes it with the given config
func GetTokenizer(name string, config map[string]any) (Tokenizer, error) {
	tokenizer, exists := store.tokenizers[name]
	if !exists {
		return nil, fmt.Errorf("tokenizer type '%s' not found", name)
	}
	return tokenizer, nil
}

// SetTokenizer sets a tokenizer in the store
// will overwrite if it already exists
func SetTokenizer(name string, tokenizer Tokenizer) {
	// Lock store
	store.mu.Lock()
	defer store.mu.Unlock()

	store.tokenizers[name] = tokenizer
}

// DeleteTokenizer deletes a tokenizer from the store
func DeleteTokenizer(name string) {
	// Lock store
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.tokenizers, name)
}
