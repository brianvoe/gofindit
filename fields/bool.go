package fields

import (
	"bytes"
	"fmt"
)

func init() {
	SetField("bool", NewBool)
}

// Bool stores the boolean value as bytes
type Bool struct {
	value []byte
}

func NewBool(config map[string]any) (Field, error) {
	return &Bool{}, nil
}

// Type returns the type of the Field.
func (b *Bool) Type() string {
	return BooleanType
}

// Value returns the stored value of the Field
func (b *Bool) Value() []byte {
	// Interpret the byte slice as a boolean
	return b.value
}

func (b *Bool) ToBool() bool {
	return b.value[0] == 1
}

// Process converts a boolean value to bytes and stores it in the Bool struct
func (b *Bool) Process(val any) error {
	boolVal, ok := val.(bool)
	if !ok {
		return fmt.Errorf("Bool requires a boolean value")
	}

	if boolVal {
		b.value = []byte{1}
	} else {
		b.value = []byte{0}
	}
	return nil
}

// ToSearchByte converts a boolean value to a byte slice
func (b *Bool) ToSearchBytes(val any) ([]byte, error) {
	boolVal, ok := val.(bool)
	if !ok {
		return nil, fmt.Errorf("Bool requires a boolean value")
	}

	if boolVal {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Search checks if the given byte slice represents the same boolean value
func (b *Bool) Search(searchValue []byte) (bool, error) {
	if len(searchValue) != 1 {
		return false, fmt.Errorf("invalid search value for Bool")
	}
	return bytes.Equal(b.value, searchValue), nil
}

// SearchRange for Bool is not applicable but implemented to satisfy the interface
func (b *Bool) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for Bool")
}
