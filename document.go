package gofindit

type Document struct {
	Original any
	Flat     map[string]any
}

func (d *Document) Get(field string) (bool, any) {
	// Check if the field exists
	_, ok := (*d)[field]
	if !ok {
		return false, nil
	}

	return true, (*d)[field]
}

// Search returns a list of document IDs that match the given field and value
func (d *Document) Search(documentQuery []DocumentQuery) (bool, error) {
	// Loop through the document query, get the value from the field
	// If any one of them does not match, return false
	// If all of them match, return true
	for _, query := range documentQuery {
		// Find the field
		found, value := d.Get(query.Name)
		if !found {
			return false, nil
		}

		// Check if the value matches the query
		switch query.Type {
		case "match":
			// TODO: Implement match
			if value != query.Value {
				return false, nil
			}
		case "partial":
			// TODO: Implement partial
			if value != query.Value {
				return false, nil
			}
		case "range":
			// TODO: Implement range

		}
	}

	return true, nil
}

// DocumentQuery is a query to search fo
type DocumentQuery struct {
	Name  string
	Type  string // "match", "partial", "range"
	Value any
}
