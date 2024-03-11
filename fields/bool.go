package fields

import (
	"bytes"
	"fmt"
)

func init() {
	RegisterField("bool", NewBool)
}

// BoolField stores the boolean value as bytes.
type BoolField struct {
	value []byte
}

func NewBool(config map[string]any) (Field, error) {
	return &BoolField{}, nil
}

// Type returns the type of the Field.
func (b *BoolField) Type() string {
	return Boolean
}

// Value returns the stored value of the Field.
func (b *BoolField) Value() any {
	// Interpret the byte slice as a boolean.
	return b.value[0] == 1
}

// Process converts a boolean value to bytes and stores it in the BoolField struct.
func (b *BoolField) Process(val any) error {
	boolVal, ok := val.(bool)
	if !ok {
		return fmt.Errorf("BoolField requires a boolean value")
	}

	if boolVal {
		b.value = []byte{1}
	} else {
		b.value = []byte{0}
	}
	return nil
}

// ToSearchByte converts a boolean value to a byte slice.
func (b *BoolField) ToSearchByte(val any) ([]byte, error) {
	boolVal, ok := val.(bool)
	if !ok {
		return nil, fmt.Errorf("BoolField requires a boolean value")
	}

	if boolVal {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Search checks if the given byte slice represents the same boolean value.
func (b *BoolField) Search(searchValue []byte) (bool, error) {
	if len(searchValue) != 1 {
		return false, fmt.Errorf("invalid search value for BoolField")
	}
	return bytes.Equal(b.value, searchValue), nil
}

// SearchRange for BoolField is not applicable but implemented to satisfy the interface.
func (b *BoolField) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for BoolField")
}
