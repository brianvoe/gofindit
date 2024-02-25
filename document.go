package gofindit

type Document struct {
	Original    any
	FieldValues map[string]any
	FieldTypes  map[string]string
}

func NewDocument(doc any) (*Document, error) {
	// Get field value map
	fieldValueMap, err := getFieldValues(doc)
	if err != nil {
		return nil, err
	}

	// Get field types
	fieldTypes, err := getFieldTypes(doc)
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

func (d *Document) GetFieldValue(field string) (any, string, bool) {
	// Check if the field exists
	_, ok := (*d).FieldValues[field]
	if !ok {
		return nil, "", false
	}

	// Get value type using reflect
	value, ok := (*d).FieldValues[field]
	if !ok {
		return nil, "", false
	}

	// Get value type using reflect
	valueType := getBaseType(value)

	return (*d).FieldValues[field], valueType, true
}
