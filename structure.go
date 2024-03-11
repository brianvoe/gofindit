package gofindit

import (
	"fmt"
	"reflect"
	"time"

	"github.com/brianvoe/gofindit/fields"
)

// getStructure returns an array of
func getStructure(v any, parent string) (map[string]fields.Field, error) {
	// Make sure v is a struct
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("v is not a struct")
	}

	// Get value if v is a pointer
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make(map[string]fields.Field)

	// Loop through fields and add name and type to fields
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name, _ := typeField.Tag.Lookup("find")
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

		//

		// Handle different types
		switch valueField.Kind() {
		case reflect.String:
			fields[name] = FieldValue{
				Type:  "string",
				Value: valueField.String(),
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fields[name] = FieldValue{
				Type:  "int",
				Value: valueField.Int(),
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fields[name] = FieldValue{
				Type:  "uint",
				Value: valueField.Uint(),
			}
		case reflect.Float32, reflect.Float64:
			fields[name] = FieldValue{
				Type:  "float",
				Value: valueField.Float(),
			}
		case reflect.Bool:
			fields[name] = FieldValue{
				Type:  "bool",
				Value: valueField.Bool(),
			}
		case reflect.Array, reflect.Slice:
			// Get the type of the slice elements
			elemType := valueField.Type().Elem()

			switch elemType.Kind() {
			case reflect.String:
				fields[name] = FieldValue{
					Type:  "[]string",
					Value: valueField.Interface().([]string),
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fields[name] = FieldValue{
					Type:  "[]int",
					Value: valueField.Interface().([]int),
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fields[name] = FieldValue{
					Type:  "[]uint",
					Value: valueField.Interface().([]uint),
				}
			case reflect.Float32, reflect.Float64:
				fields[name] = FieldValue{
					Type:  "[]float",
					Value: valueField.Interface().([]float64),
				}
			case reflect.Bool:
				fields[name] = FieldValue{
					Type:  "[]bool",
					Value: valueField.Interface().([]bool),
				}
			default:
				fields[name] = FieldValue{
					Type:  elemType.Kind().String(),
					Value: valueField.Interface(),
				}
			}
		case reflect.Struct:
			// Special handling for time.Time
			if valueField.Type() == reflect.TypeOf(time.Time{}) {
				fields[name] = FieldValue{
					Type:  "time.Time",
					Value: valueField.Interface().(time.Time),
				}
			} else {
				// Recursive call for other structs
				structFields, err := getStructure(valueField.Interface(), name)
				if err != nil {
					return nil, err
				}

				// Add struct fields to fields
				for k, v := range structFields {
					fields[k] = v
				}
			}
		default:
			// default to the type of the field
			fields[name] = FieldValue{
				Type:  valueField.Type().String(),
				Value: valueField.Interface(),
			}
		}
	}

	return fields, nil
}

// Base types
// var baseTypes = []string{"string", "int", "uint", "float", "bool", "[]string", "[]int", "[]uint", "[]float", "[]bool", "time"}

func getBaseType(value any) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	case reflect.Array, reflect.Slice:
		// Get the type of the slice by calling getBaseType on the first element
		elemType := reflect.TypeOf(value).Elem()
		return "[]" + getBaseType(reflect.New(elemType).Elem().Interface())
	case reflect.Struct:
		// Special handling for time.Time
		if reflect.TypeOf(value) == reflect.TypeOf(time.Time{}) {
			return "time"
		}

		// Recursive call for other structs
		return "struct"
	}

	return reflect.TypeOf(value).Kind().String()
}
