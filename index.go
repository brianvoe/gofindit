package gofindit

import (
	"errors"
	"math/rand/v2"
	"sync"
)

type Index struct {
	Documents map[string]*Document

	Filters []FilterFunc

	// Cache
	Cache     bool
	CacheSize int

	mu sync.RWMutex
}

type Options struct {
	// Cache
	Cache     bool // Whether or not to cache search results
	CacheSize int  // The maximum number of search results to cache

	// Filters
	Filters []FilterFunc // The filters to apply to strings
}

func New() *Index {
	index := Index{
		Documents: make(map[string]*Document),
		Filters:   DefaultFilters,
		Cache:     true,
		CacheSize: 100,
	}

	return &index
}

// NewOptions returns a new index with the given options
func NewOptions(options Options) *Index {
	index := Index{
		Documents: make(map[string]*Document),
		Filters:   options.Filters,
		Cache:     options.Cache,
		CacheSize: options.CacheSize,
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

	docNew, err := NewDoc(doc)
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
