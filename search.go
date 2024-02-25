package gofindit

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type SearchQuery struct {
	Limit  uint               `json:"limit"`
	Skip   uint               `json:"skip"`
	Sort   string             `json:"sort"`    // asc or desc
	SortBy string             `json:"sort_by"` // Field to sort by
	Fields []SearchQueryField `json:"fields"`
}

func (sq *SearchQuery) Sanatize() {
	// If limit is 0, set it to 10
	if (*sq).Limit == 0 {
		(*sq).Limit = 10
	}

	// Loop through fields and sanatize
	for i := range (*sq).Fields {
		(*sq).Fields[i].Sanatize()
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
	for _, field := range sq.Fields {
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

func (dq *SearchQueryField) Sanatize() {
	// If value is string, make it lower case
	switch dq.Value.(type) {
	case string:
		dq.Value = strings.ToLower(dq.Value.(string))
	}
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

// Search returns a list of document IDs that match the given field and value
func (i *Index) Search(searchQuery SearchQuery) ([]any, error) {
	// Set default values if none set
	searchQuery.Sanatize()

	// Validate the search query
	err := searchQuery.Validate()
	if err != nil {
		return nil, err
	}

	// If sortBy is not empty, make sure it is in the fields
	if searchQuery.SortBy != "" {
		found := false
		for _, field := range searchQuery.Fields {
			if field.Field == searchQuery.SortBy {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("sort_by field %s not found in fields", searchQuery.SortBy)
		}
	}

	// Loop through docs and run search on each one and return the ones that match
	var results []*Document
	for _, doc := range i.Documents {
		searchQueryField := searchQuery.Fields

		// If no fields, add document to results
		if len(searchQueryField) == 0 {
			results = append(results, doc)
			continue
		}

		// Establish matches
		matches := 0

		// Loop through search query fields
		for _, query := range searchQueryField {
			// Get Value and Type
			value, valueType, found := doc.GetFieldValue(query.Field)
			if !found {
				// Field not found
				continue
			}

			// Get Query Value and Type
			queryValue := query.Value
			queryType := query.Type

			// Check if the value matches the query
			switch queryType {
			case "match":
				// Convert both to string and do an exact match
				if fmt.Sprint(value) == fmt.Sprint(queryValue) {
					matches++
				}
			case "partial":
				// if bool or time then just continue
				if valueType == "bool" || valueType == "time" {
					continue
				}

				// Convert both values to string and then do strings.Contains
				switch valueType {
				// Any numbers convert to string and then do strings.Contains
				case "string", "int", "uint", "float":
					if strings.Contains(fmt.Sprint(value), fmt.Sprint(queryValue)) {
						matches++
					}
				case "[]string", "[]int", "[]uint", "[]float":
					// Loop through array and check if any of the values match
					for _, v := range value.([]any) {
						if strings.Contains(fmt.Sprint(v), fmt.Sprint(queryValue)) {
							matches++
							break // If one value matches, break the loop
						}
					}
				}
			case "range":
				// If valueType is not a number, then continue
				if valueType == "string" || valueType == "bool" {
					continue
				}

				// Initial declarations for min and max values for both numeric and time ranges
				var minNum, maxNum float64
				var minTime, maxTime time.Time
				var err error
				var isTimeRange, isNumRange bool

				// Determine the type of range query: Time or Numeric based on valueType
				switch valueType {
				case "time":
					isTimeRange = true
					// Handle time ranges
					if timeRange, ok := queryValue.([]time.Time); ok && len(timeRange) > 0 {
						minTime = timeRange[0]
						if len(timeRange) > 1 {
							maxTime = timeRange[1]
						}
					} else if singleTime, ok := queryValue.(time.Time); ok {
						minTime = singleTime
						// For single time, consider maxTime as zero value of time.Time, indicating no end range
					}
				case "int", "uint", "float":
					isNumRange = true
					// Handle numeric ranges
					if numRange, ok := queryValue.([]time.Time); ok && len(numRange) > 0 {
						minNum, err = toFloat64(numRange[0])
						if err != nil {
							continue // Skip if conversion fails
						}
						if len(numRange) > 1 {
							maxNum, err = toFloat64(numRange[1])
							if err != nil {
								continue // Skip if conversion fails
							}
						}
					} else if singleNum, ok := queryValue.(float64); ok {
						minNum = singleNum
						// For single number, consider maxNum as not set
					}
				}

				// Compare based on determined range type
				if isTimeRange {
					fieldValue := value.(time.Time)
					if !minTime.IsZero() && fieldValue.Before(minTime) {
						continue
					}
					if !maxTime.IsZero() && fieldValue.After(maxTime) {
						continue
					}
					matches++
				} else if isNumRange {
					fieldValue, err := toFloat64(value)
					if err != nil {
						continue // Skip if conversion fails
					}
					if fieldValue < minNum || (!isZero(maxNum) && fieldValue > maxNum) {
						continue
					}
					matches++
				}

			}

			// If matches is equal to the number of fields, then add the document to the results
			if matches == len(searchQueryField) {
				results = append(results, doc)
			}
		}
	}

	// Sort the results
	sortOrder := searchQuery.Sort
	sortBy := searchQuery.SortBy
	sort.SliceStable(results, func(i, j int) bool {
		// Sort by the sub field FieldValues
		if sortOrder == "desc" {
			return fmt.Sprint(results[i].FieldValues[sortBy]) > fmt.Sprint(results[j].FieldValues[sortBy])
		}

		return fmt.Sprint(results[i].FieldValues[sortBy]) < fmt.Sprint(results[j].FieldValues[sortBy])
	})

	// Handle skip
	if searchQuery.Skip > 0 {
		if int(searchQuery.Skip) > len(results) {
			results = results[:0]
		} else {
			results = results[searchQuery.Skip:]
		}
	}

	// Handle limit
	if searchQuery.Limit > 0 {
		if int(searchQuery.Limit) < len(results) {
			results = results[:searchQuery.Limit]
		}
	}

	// Loop through results and get the original document
	var originalResults []any
	for _, doc := range results {
		originalResults = append(originalResults, doc.Original)
	}

	return originalResults, nil
}
