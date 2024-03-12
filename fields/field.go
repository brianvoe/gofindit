package fields

import (
	"fmt"
	"sync"
)

// Field types
const (
	TextType    = "t"
	NumberType  = "n"
	BooleanType = "b"
	DateType    = "d"
)

// Field is an interface that all field types must implement
type Field interface {
	Type() string
	Value() []byte

	// Process will take in an any value and
	// use it to fill out the struct fields
	Process(val any) error

	// To use the struct values to calculate how search
	// bytes should be passed to the search function
	ToSearchBytes(val any) ([]byte, error)

	Search(val []byte) (bool, error)
	SearchRange(min, max []byte) (bool, error)
}

type storage struct {
	fields map[string]FieldFunc

	// lock
	mu sync.RWMutex
}

// store is the storage for all fields
var store = &storage{
	fields: make(map[string]FieldFunc),
}

// FieldFunc is a config passable function
// that returns a newly initialized Field
type FieldFunc func(config map[string]any) (Field, error)

// SetField sets a field in the store
// will overwrite if it already exists
func SetField(name string, field FieldFunc) {
	// Lock store
	store.mu.Lock()
	defer store.mu.Unlock()

	store.fields[name] = field
}

// DeleteField deletes a field from the store
func DeleteField(name string) {
	// Lock store
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.fields, name)
}

// GetField returns a field from the store
// and initializes it with the given config
func GetField(name string, config map[string]any) (Field, error) {
	fieldFunc, exists := store.fields[name]
	if !exists {
		return nil, fmt.Errorf("field type '%s' not found", name)
	}
	return fieldFunc(config)
}
