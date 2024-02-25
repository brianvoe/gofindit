package gofindit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Take full raw query string and parse it into a SearchQuery struct
// Field types are a flat structure of the field and type of field
func StringToSearchQuery(input string, fieldTypes map[string]string) (*SearchQuery, error) {
	fields := make([]SearchQueryField, 0)
	var limit, skip uint64
	var sort string

	// If the input doesnt have an = sign then we know its not a valid query
	if !strings.Contains(input, "=") {
		return nil, errors.New("invalid input format")
	}

	// Parse the input string as a URL-encoded string of parameters
	params, err := url.ParseQuery(input)
	if err != nil {
		return nil, errors.New("invalid input format")
	}

	// Check if the "limit" parameter is present
	if limitValues, ok := params["limit"]; ok {
		if len(limitValues) > 0 {
			limit, err = strconv.ParseUint(limitValues[0], 10, 64)
			if err != nil {
				return nil, errors.New("invalid value for limit parameter")
			}
		}
	}

	// Check if the "skip" parameter is present
	if skipValues, ok := params["skip"]; ok {
		if len(skipValues) > 0 {
			skip, err = strconv.ParseUint(skipValues[0], 10, 64)
			if err != nil {
				return nil, errors.New("invalid value for skip parameter")
			}
		}
	}

	// Check if the "sort" parameter is present
	if sortValues, ok := params["sort"]; ok {
		if len(sortValues) > 0 {
			sort = sortValues[0]
		}
	}

	// Iterate over the map of parameters
	for fieldName, values := range params {
		if fieldName == "limit" || fieldName == "skip" || fieldName == "sort" {
			// Skip special parameters
			continue
		}
		if len(values) == 0 {
			continue
		}

		// Get the type from the fieldTypes map
		fieldType := ""
		if _, ok := fieldTypes[fieldName]; ok {
			fieldType = fieldTypes[fieldName]
		} else {
			return nil, fmt.Errorf("field name `%s` does not exist", fieldName)
		}

		searchType := ""
		fieldValue := values[0] // Get the first value for the field

		// Check if the field value has a type prefix (e.g. "match:", "partial:", "range:")
		value := fieldValue
		if strings.Contains(fieldValue, ":") {
			parts := strings.Split(fieldValue, ":")
			if len(parts) == 2 {
				searchType = parts[0]
				value = parts[1]
			}
		}

		// Map the value to the appropriate type
		valueAny, err := MapValueStringToQueryFieldValue(value, fieldType)
		if err != nil {
			return nil, err
		}

		// If searchType is empty, check if the value is a slice
		// If it is a slice, set the searchType to "range"
		if searchType == "" {
			if _, ok := valueAny.([]int); ok {
				searchType = "range"
			}
			if _, ok := valueAny.([]float64); ok {
				searchType = "range"
			}
		}

		// Otherwise, set the searchType to "match"
		if searchType == "" {
			searchType = "match"
		}

		// Create the SearchQueryField struct with the parsed field type and value
		searchQueryField := SearchQueryField{
			Field: fieldName,
			Type:  searchType,
			Value: valueAny,
		}

		// Append the SearchQueryField to the fields slice
		fields = append(fields, searchQueryField)
	}

	// If fields is empty, set to nil
	if len(fields) == 0 {
		fields = nil
	}

	// Create the SearchQuery struct with the map of SearchQueryFields and special fields
	searchQuery := &SearchQuery{
		Limit:  uint(limit),
		Skip:   uint(skip),
		Sort:   sort,
		Fields: fields,
	}

	return searchQuery, nil
}

func MapValueStringToQueryFieldValue(value string, fieldType string) (any, error) {
	// Switch statement on fieldType
	switch fieldType {
	case "string", "[]string":
		return value, nil
	case "bool":
		return strconv.ParseBool(value)
	case "int", "float":
		// Check for slice types
		if strings.Contains(value, ",") {
			values := strings.Split(value, ",")
			var floats []float64
			var isFloat = false

			// First pass: check if any value is a float
			for _, v := range values {
				if strings.Contains(v, ".") {
					if _, err := strconv.ParseFloat(v, 64); err == nil {
						isFloat = true
						break
					}
				}
			}

			// Convert to float slice if any value is a float
			if isFloat {
				for _, v := range values {
					floatVal, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, err // return error if conversion fails
					}
					floats = append(floats, floatVal)
				}
				return floats, nil
			}

			// If no floats, check for int and bool types
			var ints []int
			var isInt = true

			for _, v := range values {
				// Try int
				if _, err := strconv.Atoi(v); err != nil {
					isInt = false
				}
			}

			// Convert to the appropriate non-float slice type
			if isInt {
				for _, v := range values {
					intVal, _ := strconv.Atoi(v)
					ints = append(ints, intVal)
				}
				return ints, nil
			}

			// return string slice
			return values, nil
		}

		// If not a slice, convert to the appropriate type
		if fieldType == "float" {
			return strconv.ParseFloat(value, 64)
		}
		if fieldType == "int" {
			return strconv.Atoi(value)
		}
	default:
		return nil, fmt.Errorf("unknown type: %s", fieldType)
	}

	return value, nil
}

func JsontoSearchQueries(jsonBytes []byte) (*SearchQuery, error) {
	var searchQuery SearchQuery
	err := json.Unmarshal(jsonBytes, &searchQuery)
	if err != nil {
		return nil, err
	}
	return &searchQuery, nil
}
