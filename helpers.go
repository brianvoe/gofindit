package gofindit

import (
	"fmt"
	"reflect"
)

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
