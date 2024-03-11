package fields

import (
	"bytes"
	"fmt"
)

func init() {
	RegisterField("text", NewText)
}

// TextField stores the text value as a byte slice.
type TextField struct {
	value []byte
}

// NewText
func NewText(config map[string]any) (Field, error) {
	return &TextField{}, nil
}

// Type returns the type of the Field.
func (t *TextField) Type() string {
	return Text
}

// Value returns the stored value of the Field as a string.
func (t *TextField) Value() any {
	return string(t.value)
}

// ToByte converts a text value to a byte slice.
func (t *TextField) ToByte(val any) ([]byte, error) {
	strVal, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("TextField requires a string value")
	}
	return []byte(strVal), nil
}

// Process converts a text value to a byte slice and stores it in the TextField struct.
func (t *TextField) Process(val any) error {
	strVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("TextField requires a string value")
	}
	t.value = []byte(strVal)
	return nil
}

// Search performs an efficient comparison between the stored byte slice and the search value.
func (t *TextField) Search(searchValue []byte) (bool, error) {
	return bytes.Equal(t.value, searchValue), nil
}

// SearchRange for TextField is not applicable but implemented to satisfy the interface.
func (t *TextField) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for TextField")
}
