package gofindit

import (
	"fmt"
	"time"
)

type SearchQuery struct {
	Limit  uint               `json:"limit"`
	Skip   uint               `json:"skip"`
	Sort   string             `json:"sort"`    // asc or desc
	SortBy string             `json:"sort_by"` // Field to sort by
	Fields []SearchQueryField `json:"fields"`
}

func (sq *SearchQuery) Default() {
	// If limit is 0, set it to 10
	if (*sq).Limit == 0 {
		(*sq).Limit = 10
	}
}

func (sq *SearchQuery) Validate() error {
	// Check if the type is valid
	if sq.Sort != "" && sq.Sort != "asc" && sq.Sort != "desc" {
		return fmt.Errorf("invalid sort type %s", sq.Sort)
	}

	// Make sure sortBy is not empty if sort is set
	if sq.Sort != "" && sq.SortBy == "" {
		return fmt.Errorf("sort_by cannot be set without sort")
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

type SearchQueryField struct {
	Field string
	Type  string // "match", "partial", "range"
	Value any
}

func (dq *SearchQueryField) Validate() error {
	// Check if the name is not empty
	if dq.Field == "" {
		return fmt.Errorf("field name cannot be empty")
	}

	// Check if the type is empty
	if dq.Type == "" {
		dq.Type = "match"
	}

	// Check if the type is valid
	if dq.Type != "match" && dq.Type != "partial" && dq.Type != "range" {
		return fmt.Errorf("invalid type %s", dq.Type)
	}

	// Check for types like bool and time on partial, make invalid
	if dq.Type == "partial" {
		switch dq.Value.(type) {
		case bool:
			return fmt.Errorf("cannot use partial search on bool type")
		case time.Time:
			return fmt.Errorf("cannot use partial search on time type")
		}
	}

	// Check type for range and if bool or string, make invalid
	if dq.Type == "range" {
		switch dq.Value.(type) {
		case bool:
			return fmt.Errorf("cannot use range search on bool type")
		case string:
			return fmt.Errorf("cannot use range search on string type")
		}
	}

	return nil
}
