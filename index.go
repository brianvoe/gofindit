package gofindit

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"
	"time"
)

type Index struct {
	Documents map[string]*Document

	mu sync.RWMutex
}

func New() *Index {
	index := Index{
		Documents: make(map[string]*Document),
	}

	return &index
}

func (i *Index) Random() (string, any) {
	// Get array of document keys
	var keys []string
	for k := range i.Documents {
		keys = append(keys, k)
	}

	// Get random key
	key := keys[rand.IntN(len(keys))]
	return key, i.Documents[key].Original
}

func (i *Index) Index(id string, doc any) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Make sure id isnt taken
	if _, ok := i.Documents[id]; ok {
		return errors.New("id already taken")
	}

	docNew, err := NewDocument(doc)
	if err != nil {
		return err
	}

	i.Documents[id] = docNew

	return nil
}

// Get returns the document with the given ID
func (i *Index) Get(id string) (any, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	doc, ok := i.Documents[id]
	if !ok {
		return nil, errors.New("document not found")
	}

	return doc.Original, nil
}

// Search returns a list of document IDs that match the given field and value
func (i *Index) Search(searchQuery SearchQuery) ([]any, error) {
	// Set default values if none set
	searchQuery.Default()

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
