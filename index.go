package gofindit

import (
	"sync"
)

type Index struct {
	Name      string
	Documents map[string]Document

	Structure []StructKeyType

	mu sync.RWMutex
}

func New(name string, structure any) (*Index, error) {
	index := Index{
		Name: name,
	}

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

	i.Documents[id] = doc

	return nil
}

// Search returns a list of document IDs that match the given field and value
func (i *Index) Search(searchQuery SearchQuery) ([]string, error) {

}
