package gofindit

import (
	"github.com/brianvoe/gofindit/fields"
)

type Document struct {
	Original any
	Values   map[string]fields.Field
}

func NewDoc(doc any) (*Document, error) {
	// Get structure of the document
	values, err := getStructure(doc, "")
	if err != nil {
		return nil, err
	}

	// Create a new document
	document := Document{
		Original: doc,
		Values:   values,
	}

	return &document, nil
}

func (d *Document) GetFieldValue(field string) (any, bool) {
	// Check if the field exists
	value, ok := (*d).Values[field]
	if !ok {
		return nil, false
	}

	return value.Value, true
}

func (d *Document) GetFieldType(field string) (string, bool) {
	// Check if the field exists
	value, ok := (*d).Values[field]
	if !ok {
		return "", false
	}

	return value.Type, true
}
