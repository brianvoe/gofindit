package gofindit

import (
	"fmt"
	"reflect"
	"time"

	"github.com/brianvoe/gofindit/fields"
)

func getStructure(v any, parent string) (map[string]fields.Field, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("v is not a struct")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fieldsFinal := make(map[string]fields.Field)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name, _ := typeField.Tag.Lookup("find")
		if name == "-" {
			continue
		}

		if name == "" {
			name = typeField.Name
		}

		if parent != "" {
			name = parent + "." + name
		}

		fieldTag := typeField.Tag.Get("field")

		switch valueField.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool:
			// Simplified handling for basic types
			basicField, err := getBasicField(valueField, fieldTag)
			if err != nil {
				return nil, err
			}
			fieldsFinal[name] = basicField
		case reflect.Array, reflect.Slice:
			elemType := valueField.Type().Elem()
			if elemType.Kind() == reflect.Struct {
				// Handle slice of structs
				for j := 0; j < valueField.Len(); j++ {
					elemValue := valueField.Index(j)
					structFields, err := getStructure(elemValue.Interface(), fmt.Sprintf("%s[%d]", name, j))
					if err != nil {
						return nil, err
					}
					for k, v := range structFields {
						fieldsFinal[k] = v
					}
				}
			} else {
				// Simplified handling for slices of basic types
				basicField, err := getBasicField(valueField, fieldTag)
				if err != nil {
					return nil, err
				}
				fieldsFinal[name] = basicField
			}
		case reflect.Struct:
			// Recursive call for nested structs
			structFields, err := getStructure(valueField.Interface(), name)
			if err != nil {
				return nil, err
			}
			for k, v := range structFields {
				fieldsFinal[k] = v
			}
		}
	}

	return fieldsFinal, nil
}

// getBasicField handles the creation of fields based on basic types.
func getBasicField(valueField reflect.Value, fieldTag string) (fields.Field, error) {
	// Determine the field type and create the appropriate Field
	switch valueField.Kind() {
	// String
	case reflect.String:
		if fieldTag == "" {
			fieldTag = fields.TextType
		}

	// Number
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if fieldTag == "" {
			fieldTag = fields.NumberType
		}

	// Bool
	case reflect.Bool:
		if fieldTag == "" {
			fieldTag = fields.BooleanType
		}

	// Special handling for time.Time
	case reflect.Struct:
		if valueField.Type() == reflect.TypeOf(time.Time{}) {
			if fieldTag == "" {
				fieldTag = fields.DateType
			}
		} else {
			return nil, fmt.Errorf("struct type not supported directly, consider using nested struct handling")
		}

	default:
		return nil, fmt.Errorf("unsupported type: %v", valueField.Type())
	}

	return fields.GetField(fieldTag, nil)
}
