package gofindit

import (
	"fmt"
	"reflect"
	"time"
)

type StructField struct {
	Name      string
	Type      string
	Value     any
	Supported bool
}

// Using getStructure to get the structure of a struct
// and output map[string]any (name and value)
func getFieldValueMap(v any) (map[string]any, error) {
	// Get the structure of the document
	structure, err := getStructure(v, "")
	if err != nil {
		return nil, err
	}

	// Create a map of the structure
	structureMap := make(map[string]any)
	for _, field := range structure {
		structureMap[field.Name] = field.Value
	}

	return structureMap, nil
}

// getStructure returns an array of
func getStructure(v any, parent string) ([]StructField, error) {
	// Make sure v is a struct
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("v is not a struct")
	}

	// Get value if v is a pointer
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make([]StructField, 0)

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
			fields = append(fields, StructField{
				Name:      name,
				Type:      "string",
				Value:     valueField.String(),
				Supported: true,
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fields = append(fields, StructField{
				Name:      name,
				Type:      "int",
				Value:     valueField.Int(),
				Supported: true,
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fields = append(fields, StructField{
				Name:      name,
				Type:      "uint",
				Value:     valueField.Uint(),
				Supported: true,
			})
		case reflect.Float32, reflect.Float64:
			fields = append(fields, StructField{
				Name:      name,
				Type:      "float",
				Value:     valueField.Float(),
				Supported: true,
			})
		case reflect.Bool:
			fields = append(fields, StructField{
				Name:      name,
				Type:      "bool",
				Value:     valueField.Bool(),
				Supported: true,
			})
		case reflect.Array, reflect.Slice:
			// Get the type of the slice elements
			elemType := valueField.Type().Elem()

			switch elemType.Kind() {
			case reflect.String:
				fields = append(fields, StructField{
					Name:      name,
					Type:      "[]string",
					Value:     valueField.Interface().([]string),
					Supported: true,
				})
			case reflect.Int:
				fields = append(fields, StructField{
					Name:      name,
					Type:      "[]int",
					Value:     valueField.Interface().([]int),
					Supported: true,
				})
			default:
				fields = append(fields, StructField{
					Name:      name,
					Type:      elemType.Kind().String(),
					Value:     valueField.Interface(),
					Supported: false,
				})
			}
		case reflect.Struct:
			// Special handling for time.Time
			if valueField.Type() == reflect.TypeOf(time.Time{}) {
				fields = append(fields, StructField{
					Name:      name,
					Type:      "time.Time",
					Value:     valueField.Interface().(time.Time),
					Supported: true,
				})
			} else {
				// Recursive call for other structs
				structFields, err := getStructure(valueField.Interface(), name)
				if err != nil {
					return nil, err
				}

				// Add struct fields to fields
				fields = append(fields, structFields...)
			}
		default:
			// default to the type of the field
			fields = append(fields, StructField{
				Name:      name,
				Type:      valueField.Kind().String(),
				Value:     valueField.Interface(),
				Supported: false,
			})
		}
	}

	return fields, nil
}
