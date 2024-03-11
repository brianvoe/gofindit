package fields

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func init() {
	RegisterField("num", NewNum)
}

// NumField stores the numeric value directly as bytes.
type NumField struct {
	val []byte
}

func NewNum(config map[string]any) (Field, error) {
	return &NumField{}, nil
}

// NumToBytes converts a numeric value to a byte slice.
func numToBytes(value any) ([]byte, error) {
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
		return nil, fmt.Errorf("unsupported type for NumField: %T", v)
	}

	if err != nil {
		return nil, fmt.Errorf("error converting numeric value to bytes: %v", err)
	}

	return buf.Bytes(), nil
}

func (n *NumField) Type() string {
	return Number
}

func (n *NumField) Value() any {
	return n.Value
}

func (n *NumField) ToByte(val any) ([]byte, error) {
	return numToBytes(val)
}

// Process converts a numeric value to bytes and stores it in the NumField struct using NumToBytes.
func (n *NumField) Process(val any) error {
	bytes, err := numToBytes(val)
	if err != nil {
		return fmt.Errorf("failed to process numeric value: %v", err)
	}
	n.val = bytes
	return nil
}

// Search compares the given byte slice directly with the NumField's stored byte slice.
func (n *NumField) Search(val []byte) (bool, error) {
	return bytes.Equal(n.val, val), nil
}

// SearchRange checks if the stored value is within the given range [min, max].
func (n *NumField) SearchRange(min, max []byte) (bool, error) {
	return bytes.Compare(n.val, min) >= 0 && bytes.Compare(n.val, max) <= 0, nil
}
