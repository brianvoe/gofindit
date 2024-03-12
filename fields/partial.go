package fields

import (
	"bytes"
	"fmt"
	"strings"
)

func init() {
	SetField("partial", NewPartial)
}

// Partial stores the text value as a byte slice for partial text search
type Partial struct {
	value []byte
}

// NewPartial creates a new Partial field
func NewPartial(config map[string]any) (Field, error) {
	return &Partial{}, nil
}

func stringToBytes(val any) ([]byte, error) {
	strVal, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("value requires a string")
	}

	// Lowercase the string and convert it to a byte slice
	strVal = strings.ToLower(strVal)

	return []byte(strVal), nil
}

// Type returns the type of the Field
func (p *Partial) Type() string {
	return TextType
}

// Value returns the stored value of the Field as a byte slice
func (p *Partial) Value() []byte {
	return p.value
}

// Process converts a text value to a byte slice and stores it in the Partial struct
func (p *Partial) Process(val any) error {
	bytes, err := stringToBytes(val)
	if err != nil {
		return err
	}
	p.value = bytes
	return nil
}

func (p *Partial) ToSearchBytes(val any) ([]byte, error) {
	return stringToBytes(val)
}

// Search checks if the stored value contains the search value as a substring
func (p *Partial) Search(searchValue []byte) (bool, error) {
	// Use bytes.Contains to check if p.value contains searchValue
	return bytes.Contains(p.value, searchValue), nil
}

// SearchRange for Partial is not applicable but implemented to satisfy the interface
func (p *Partial) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for Partial")
}
