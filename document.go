package gofindit

type Document struct {
	Original    any
	FieldValues map[string]any
	FieldTypes  map[string]string
}

func NewDocument(doc any) (*Document, error) {
	// Get structure of the document
	structure, err := getStructure(doc, "")
	if err != nil {
		return nil, err
	}

	// Get field value map
	fieldValueMap, err := getFieldValues(structure)
	if err != nil {
		return nil, err
	}

	// Get field types
	fieldTypes, err := getFieldTypes(structure)
	if err != nil {
		return nil, err
	}

	// Create a new document
	document := Document{
		Original:    doc,
		FieldValues: fieldValueMap,
		FieldTypes:  fieldTypes,
	}

	return &document, nil
}

func (d *Document) GetFieldValue(field string) (any, bool) {
	// Check if the field exists
	value, ok := (*d).FieldValues[field]
	if !ok {
		return nil, false
	}

	return value, true
}

func (d *Document) GetFieldType(field string) (string, bool) {
	// Check if the field exists
	value, ok := (*d).FieldTypes[field]
	if !ok {
		return "", false
	}

	return value, true
}
