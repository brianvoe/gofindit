package gofindit

import (
	"errors"
	"math/rand/v2"
	"sync"
)

type Index struct {
	Documents map[string]*Document

	mu sync.RWMutex
}

func New() *Index {
	index := Index{
		Documents: make(map[string]*Document),
	}

	return &index
}

func (i *Index) Random() (string, any) {
	// Get array of document keys
	var keys []string
	for k := range i.Documents {
		keys = append(keys, k)
	}

	// Get random key
	key := keys[rand.IntN(len(keys))]
	return key, i.Documents[key].Original
}

func (i *Index) Index(id string, doc any) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Make sure id isnt taken
	if _, ok := i.Documents[id]; ok {
		return errors.New("id already taken")
	}

	docNew, err := NewDocument(doc)
	if err != nil {
		return err
	}

	i.Documents[id] = docNew

	return nil
}

// Get returns the document with the given ID
func (i *Index) Get(id string) (any, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	doc, ok := i.Documents[id]
	if !ok {
		return nil, errors.New("document not found")
	}

	return doc.Original, nil
}
