package gofindit

import (
	"fmt"
	"reflect"
	"time"
)

// Supported types
var supportedTypes = []string{"string", "int", "uint", "float", "bool", "time.Time", "[]string", "[]int"}

func isSupportedType(t string) bool {
	for _, supportedType := range supportedTypes {
		if t == supportedType {
			return true
		}
	}

	return false
}

type StructKeyType struct {
	Name      string
	Type      string
	Supported bool
}

// getStructure returns a map of the structure of a struct
// in a flat map[string]string
func getStructure(v any, parent string) ([]StructKeyType, error) {
	// Make sure v is a struct
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("v is not a struct")
	}

	// Get value if v is a pointer
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make([]StructKeyType, 0)

	// Loop through fields and add name and type to fields
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Tag.Get("find")
		if name == "-" {
			continue
		}

		// If name is empty, use the field name
		if name == "" {
			name = typeField.Name
		}

		// If parent is not empty, add a dot
		if parent != "" {
			name = parent + "." + name
		}

		// Handle different types
		switch valueField.Kind() {
		case reflect.String:
			fields = append(fields, StructKeyType{Name: name, Type: "string", Supported: true})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fields = append(fields, StructKeyType{Name: name, Type: "int", Supported: true})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fields = append(fields, StructKeyType{Name: name, Type: "uint", Supported: true})
		case reflect.Float32, reflect.Float64:
			fields = append(fields, StructKeyType{Name: name, Type: "float", Supported: true})
		case reflect.Bool:
			fields = append(fields, StructKeyType{Name: name, Type: "bool", Supported: true})
		case reflect.Array, reflect.Slice:
			// Get the type of the slice elements
			elemType := valueField.Type().Elem()

			switch elemType.Kind() {
			case reflect.String:
				fields = append(fields, StructKeyType{Name: name, Type: "[]string", Supported: true})
			case reflect.Int:
				fields = append(fields, StructKeyType{Name: name, Type: "[]int", Supported: true})
			default:
				fields = append(fields, StructKeyType{Name: name, Type: elemType.Kind().String(), Supported: false})
			}
		case reflect.Struct:
			// Special handling for time.Time
			if valueField.Type() == reflect.TypeOf(time.Time{}) {
				fields = append(fields, StructKeyType{Name: name, Type: "time.Time"})
			} else {
				// Recursive call for other structs
				structFields, err := getStructure(valueField.Interface(), name)
				if err != nil {
					return nil, err
				}

				// Add struct fields to fields
				for _, structField := range structFields {
					fields = append(fields, structField)
				}
			}
		default:
			// default to the type of the field
			fields = append(fields, StructKeyType{Name: name, Type: valueField.Kind().String()})
		}
	}

	return fields, nil
}
