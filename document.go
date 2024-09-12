package gofindit

import (
	"github.com/brianvoe/gofindit/fields"
)

type Document struct {
	Original any
	Fields   map[string]fields.Field
}

func NewDoc(doc any) (*Document, error) {
	// Get structure of the document
	fields, err := getStructure(doc, "")
	if err != nil {
		return nil, err
	}

	// Create a new document
	document := Document{
		Original: doc,
		Fields:   fields,
	}

	return &document, nil
}

func (d *Document) GetField(field string) (fields.Field, bool) {
	val, ok := d.Fields[field]
	return val, ok
}
