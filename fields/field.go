package fields

import "fmt"

// Field types
const (
	Text    = "t"
	Number  = "n"
	Boolean = "b"
	Date    = "d"
)

type Field interface {
	Type() string
	Value() any

	// To use the struct values to calculate how search bytes
	// should be passed to the search function
	ToSearchByte(val any) ([]byte, error)

	// Process will take in an any value
	// and use it to fill out the struct fields
	Process(val any) error

	Search(val []byte) (bool, error)
	SearchRange(min, max []byte) (bool, error)
}

type FieldFunc func(config map[string]any) (Field, error)

var Store = make(map[string]FieldFunc, 0)

func RegisterField(name string, field FieldFunc) {
	Store[name] = field
}

func GetField(name string, config map[string]any) (Field, error) {
	fieldFunc, exists := Store[name]
	if !exists {
		return nil, fmt.Errorf("field type '%s' not found", name)
	}
	return fieldFunc(config)
}