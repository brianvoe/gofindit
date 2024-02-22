package gofindit

import "fmt"

type SearchQuery struct {
	Limit  uint            `json:"limit"`
	Skip   uint            `json:"skip"`
	Sort   string          `json:"sort"`
	Fields []DocumentQuery `json:"fields"`
}

func (sq *SearchQuery) Default() {
	// If limit is 0, set it to 10
	if (*sq).Limit == 0 {
		(*sq).Limit = 10
	}

	// If Sort is empty, set it to asc
	if (*sq).Sort == "" {
		(*sq).Sort = "asc"
	}
}

func (sq *SearchQuery) Validate() error {
	// Check if the type is valid
	switch (*sq).Sort {
	case "asc":
	case "desc":
	default:
		return fmt.Errorf("invalid sort type %s", (*sq).Sort)
	}

	// Check if the fields are valid
	for _, field := range (*sq).Fields {
		err := field.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
