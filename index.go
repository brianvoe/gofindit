package gofindit

import (
	"errors"
	"sync"
)

type Index struct {
	Name      string
	Documents map[string]*Document

	Structure []StructField

	mu sync.RWMutex
}

func New(name string, structure any) (*Index, error) {
	index := Index{
		Name: name,
	}

	// Initialize the documents map
	index.Documents = make(map[string]*Document)

	// Get structure of struct
	s, err := getStructure(structure, "")
	if err != nil {
		return nil, err
	}
	index.Structure = s

	return &index, nil
}

func (i *Index) Index(id string, doc any) error {
	i.mu.Lock()
	defer i.mu.Unlock()

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

// Search returns a list of document IDs that match the given field and value
func (i *Index) Search(searchQuery SearchQuery) ([]any, error) {
	// Set default values if none set
	searchQuery.Default()

	// Validate the search query
	err := searchQuery.Validate()
	if err != nil {
		return nil, err
	}

	// Loop through docs and run search on each one and return the ones that match
	var results []any
	for _, doc := range i.Documents {
		match, err := doc.Search(searchQuery.Fields)
		if err != nil {
			return nil, err
		}
		if match {
			results = append(results, doc.Original)
		}
	}

	return results, nil
}
