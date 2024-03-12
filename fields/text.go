package fields

import (
	"bytes"
	"fmt"
)

func init() {
	SetField("text", NewText)
}

// Text stores the text value as a byte slice
type Text struct {
	value []byte
}

// NewText creates a new Text that will do an exact match search
func NewText(config map[string]any) (Field, error) {
	return &Text{}, nil
}

// Type returns the type of the Field
func (t *Text) Type() string {
	return TextType
}

// Value returns the stored value of the Field as a string
func (t *Text) Value() []byte {
	return t.value
}

func (t *Text) ToString() string {
	return string(t.value)
}

// Process converts a text value to a byte slice and stores it in the Text struct
func (t *Text) Process(val any) error {
	strVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("Text requires a string value")
	}
	t.value = []byte(strVal)
	return nil
}

// ToSearchByte converts a text value to a byte slice
func (t *Text) ToSearchBytes(val any) ([]byte, error) {
	strVal, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("Text requires a string value")
	}
	return []byte(strVal), nil
}

// Search performs an efficient comparison between the stored byte slice and the search value
func (t *Text) Search(searchValue []byte) (bool, error) {
	return bytes.Equal(t.value, searchValue), nil
}

// SearchRange for Text is not applicable but implemented to satisfy the interface
func (t *Text) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for Text")
}
