package fields

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

func init() {
	RegisterField("date", NewDate)
}

type DateField struct {
	value       []byte
	granularity string
}

// NewDate creates a new DateField with the given configuration
func NewDate(config map[string]any) (Field, error) {
	// Default granularity is "day"
	granularity := "day"
	if val, ok := config["granularity"]; ok {
		if gran, ok := val.(string); ok && isValidGranularity(gran) {
			granularity = gran
		} else {
			return nil, fmt.Errorf("invalid granularity value")
		}
	}

	return &DateField{granularity: granularity}, nil
}

// dateToBytes adjusts the provided time.Time value according to the specified granularity and converts it to bytes.
func dateToBytes(date any, granularity string) ([]byte, error) {
	dateVal, ok := date.(time.Time)
	if !ok {
		return nil, fmt.Errorf("DateField requires a time.Time value")
	}

	adjustedDate := adjustDateToGranularity(dateVal, granularity)
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, adjustedDate.Unix()); err != nil {
		return nil, fmt.Errorf("error converting date value to bytes: %v", err)
	}
	return buf.Bytes(), nil
}

func (d *DateField) Type() string {
	return "date"
}

func (d *DateField) Value() any {
	var timestamp int64
	buf := bytes.NewReader(d.value)
	binary.Read(buf, binary.BigEndian, &timestamp)
	return time.Unix(timestamp, 0)
}

func (d *DateField) ToByte(val any) ([]byte, error) {
	return dateToBytes(val, d.granularity)
}

func (d *DateField) Process(dateVal any) error {
	bytes, err := dateToBytes(dateVal, d.granularity)
	if err != nil {
		return err
	}
	d.value = bytes
	return nil
}

func (d *DateField) Search(searchValue []byte) (bool, error) {
	return bytes.Equal(d.value, searchValue), nil
}

func (d *DateField) SearchRange(min, max []byte) (bool, error) {
	return bytes.Compare(d.value, min) >= 0 && bytes.Compare(d.value, max) <= 0, nil
}

// adjustDateToGranularity is now a private method of DateField
func adjustDateToGranularity(t time.Time, granularity string) time.Time {
	// Map the granularity string to actual adjustments
	switch granularity {
	case "year":
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	case "month":
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "day":
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "hour":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case "minute":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case "second":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	default:
		return t // Default to no adjustment if granularity is unknown or invalid.
	}
}

// isValidGranularity checks if the provided granularity string is valid.
func isValidGranularity(granularity string) bool {
	switch granularity {
	case "year", "month", "day", "hour", "minute", "second":
		return true
	}
	return false
}
