package gofindit

import (
	"fmt"
	"reflect"
	"strings"
)

type Document struct {
	Original   any
	FieldValue map[string]any
}

func NewDocument(doc any) (*Document, error) {
	// Get field value map
	fieldValueMap, err := getFieldValueMap(doc)
	if err != nil {
		return nil, err
	}

	// Create a new document
	document := Document{
		Original:   doc,
		FieldValue: fieldValueMap,
	}

	return &document, nil
}

func (d *Document) GetFieldValue(field string) (bool, any) {
	// Check if the field exists
	_, ok := (*d).FieldValue[field]
	if !ok {
		return false, nil
	}

	return true, (*d).FieldValue[field]
}

// Search returns a list of document IDs that match the given field and value
func (d *Document) Search(documentQuery []DocumentQuery) (bool, error) {
	// Validate the document query
	for _, query := range documentQuery {
		err := query.Validate()
		if err != nil {
			return false, err
		}
	}

	// Loop through the document query, get the value from the field
	// If any one of them does not match, return false
	// If all of them match, return true
	for _, query := range documentQuery {
		// Find the field
		found, value := d.GetFieldValue(query.Name)
		if !found {
			return false, fmt.Errorf("field %s not found in struct", query.Name)
		}

		// Get document value reflect type
		fieldType := reflect.ValueOf(d).FieldByName(query.Name)
		fmt.Println(fieldType)

		// Check if the value matches the query
		switch query.Type {
		case "match":
			// Compare values
			if value != query.Value {
				return false, fmt.Errorf("field %s does not match value %s", query.Name, query.Value)
			}
		case "partial":
			// Convert both values to string
			// and then do strings.Contains
			switch value.(type) {
			case string:
				if !strings.Contains(value.(string), query.Value.(string)) {
					return false, fmt.Errorf("field %s does not contain value %s", query.Name, query.Value)
				}
			default:
				return false, fmt.Errorf("field %s is not a string", query.Name)
			}
		case "range":
			// TODO: Implement range

		}
	}

	// Nothing failed so return true
	return true, nil
}

// DocumentQuery is a query to search fo
type DocumentQuery struct {
	Name  string
	Type  string // "match", "partial", "range"
	Value any
}

func (dq *DocumentQuery) Validate() error {
	// Check if the name is not empty
	if dq.Name == "" {
		return fmt.Errorf("field name cannot be empty")
	}

	// Check if the type is empty
	if dq.Type == "" {
		dq.Type = "match"
	}

	// Check if the type is valid
	switch dq.Type {
	case "match":
	case "partial":
	case "range":
	default:
		return fmt.Errorf("invalid search type %s", (*dq).Type)
	}

	// Check if the value is valid

	return nil
}
