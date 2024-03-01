package gofindit

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// stringToValue attempts to convert a string to any of the basic types
// and array of basic types.
func stringToAny(value string) (any, error) {
	// If the string contains commas, treat it as a CSV
	if strings.Contains(value, ",") {
		elements := strings.Split(value, ",")
		var boolArray []bool
		var intArray []int
		var floatArray []float64
		var stringArray []string
		var anyArray []any

		// Flags to determine the type of the array
		typeFlags := []string{}

		for _, element := range elements {
			trimmedElement := strings.TrimSpace(element)

			// Check if it can be a bool
			if trimmedElement == "true" || trimmedElement == "false" {
				typeFlags = append(typeFlags, "bool")
				continue
			}

			// Check if it can be an int
			if _, err := strconv.Atoi(trimmedElement); err == nil {
				typeFlags = append(typeFlags, "int")
				continue
			}

			// Check if it can be a float
			if _, err := strconv.ParseFloat(trimmedElement, 64); err == nil {
				typeFlags = append(typeFlags, "float")
				continue
			}

			// If it's not a bool, int or float, it's a string
			typeFlags = append(typeFlags, "string")
		}

		// Loop through the type flags and if they are not all the same, then it's an array of any
		shouldBeAny := false
		for _, flag := range typeFlags {
			if flag != typeFlags[0] {
				shouldBeAny = true
				break
			}
		}

		// Loop through the elements and convert them to the appropriate type
		for _, element := range elements {
			trimmedElement := strings.TrimSpace(element)

			if shouldBeAny {
				anyArray = append(anyArray, trimmedElement)
				continue
			}
			if typeFlags[0] == "bool" {
				boolVal, _ := strconv.ParseBool(trimmedElement)
				boolArray = append(boolArray, boolVal)
			}
			if typeFlags[0] == "int" {
				intVal, _ := strconv.Atoi(trimmedElement)
				intArray = append(intArray, intVal)
			}
			if typeFlags[0] == "float" {
				floatVal, _ := strconv.ParseFloat(trimmedElement, 64)
				floatArray = append(floatArray, floatVal)
			}
			if typeFlags[0] == "string" {
				stringArray = append(stringArray, trimmedElement)
			}
		}

		// If it's an array of any, return the any array
		if shouldBeAny {
			return anyArray, nil
		}

		// If it's an array of bool, return the bool array
		if typeFlags[0] == "bool" {
			return boolArray, nil
		}

		// If it's an array of int, return the int array
		if typeFlags[0] == "int" {
			return intArray, nil
		}

		// If it's an array of float, return the float array
		if typeFlags[0] == "float" {
			return floatArray, nil
		}

		// If it's an array of string, return the string array
		if typeFlags[0] == "string" {
			return stringArray, nil
		}
	}

	// For single values, the original logic is preserved
	if value == "true" || value == "false" {
		return strconv.ParseBool(value)
	}
	if val, err := strconv.Atoi(value); err == nil {
		return val, nil
	}
	if val, err := strconv.ParseFloat(value, 64); err == nil {
		return val, nil
	}
	return value, nil
}

// toFloat64 attempts to convert an any to a float64.
// Returns the converted float64 value and nil error if successful, 0 and an error otherwise.
func toFloat64(value any) (float64, error) {
	var result float64
	var err error

	// Use reflection to determine the type and convert accordingly
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		result = v.Float()
	default:
		err = fmt.Errorf("toFloat64: cannot convert %v to float64", reflect.TypeOf(value))
	}

	return result, err
}

func isZero(value any) bool {
	// Use reflect to get the value
	v := reflect.ValueOf(value)

	// Check if the value is valid (not nil)
	if !v.IsValid() {
		return true // nil is considered as "zero"
	}

	// Different handling for different kinds of types
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan, reflect.Interface:
		// For these types, consider them "zero" if they are nil
		return v.IsNil()
	default:
		// Use the IsZero method to check for zero value of other types
		return v.IsZero()
	}
}
