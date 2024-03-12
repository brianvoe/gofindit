package fields

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func init() {
	SetField("num", NewNum)
}

// Num stores the numeric value directly as bytes
type Num struct {
	v     any // original value
	value []byte
}

// NewNum creates a new Num that will do an exact and range search
func NewNum(config map[string]any) (Field, error) {
	return &Num{}, nil
}

// NumToSearchBytes converts a numeric value to a byte slice
func numToSearchBytes(value any) ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)

	switch v := value.(type) {
	case int:
		err = binary.Write(buf, binary.BigEndian, int64(v))
	case int8, int16, int32, int64, float32, float64:
		err = binary.Write(buf, binary.BigEndian, v)
	case uint:
		err = binary.Write(buf, binary.BigEndian, uint64(v))
	case uint8, uint16, uint32, uint64:
		err = binary.Write(buf, binary.BigEndian, v)
	default:
		return nil, fmt.Errorf("unsupported type for Num: %T", v)
	}

	if err != nil {
		return nil, fmt.Errorf("error converting numeric value to bytes: %v", err)
	}

	return buf.Bytes(), nil
}

func (n *Num) Type() string {
	return NumberType
}

func (n *Num) Value() any {
	return n.v
}

// Process converts a numeric value to bytes
// and stores it in the Num struct using NumToSearchBytes
func (n *Num) Process(val any) error {
	// Set original value
	n.v = val

	bytes, err := numToSearchBytes(val)
	if err != nil {
		return fmt.Errorf("failed to process numeric value: %v", err)
	}
	n.value = bytes
	return nil
}

func (n *Num) ToSearchBytes(val any) ([]byte, error) {
	return numToSearchBytes(val)
}

// Search compares the given byte slice directly with the Num's stored byte slice
func (n *Num) Search(val []byte) (bool, error) {
	return bytes.Equal(n.value, val), nil
}

// SearchRange checks if the stored value is within the given range [min, max]
func (n *Num) SearchRange(min, max []byte) (bool, error) {
	return bytes.Compare(n.value, min) >= 0 && bytes.Compare(n.value, max) <= 0, nil
}
