package gofindit

type Document struct {
	Original any
	Values   map[string]Field
}

func NewDoc(doc any) (*Document, error) {
	return NewDocFilters(doc, DefaultFilters...)
}

func NewDocFilters(doc any, filters ...FilterFunc) (*Document, error) {
	// Get structure of the document
	values, err := getStructure(doc, "")
	if err != nil {
		return nil, err
	}
	// Loop through values and
	// if they are a string, tokenize them
	for k, v := range values {
		if v.Type() == "string" {
			// Tokenize the string
			tokens, error := Tokenize(v.Value.(string), filters)
			if error != nil {
				return nil, error
			}
			values[k] = FieldValue{
				Type:   values[k].Type,
				Value:  values[k].Value,
				Tokens: tokens,
			}
		}
	}

	// Create a new document
	document := Document{
		Original: doc,
		Values:   values,

		Filters: filters,
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
