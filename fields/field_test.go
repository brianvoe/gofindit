package fields

import (
	"bytes"
	"fmt"
	"testing"
)

// MockField for testing Field interface functionalities
type MockField struct {
	val []byte
}

func NewMockField(_ map[string]any) (Field, error) {
	return &MockField{}, nil
}

func (m *MockField) Type() string {
	return "mock"
}

func (m *MockField) Value() []byte {
	return m.val
}

func (m *MockField) ToSearchBytes(val any) ([]byte, error) {
	strVal, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("mock field requires a string value for search")
	}
	return []byte(strVal), nil
}

func (m *MockField) Process(val any) error {
	strVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("mock field requires a string value")
	}
	m.val = []byte(strVal)
	return nil
}

func (m *MockField) Search(val []byte) (bool, error) {
	return bytes.Equal(m.val, val), nil
}

func (m *MockField) SearchRange(min, max []byte) (bool, error) {
	return false, fmt.Errorf("range search not supported for mock field")
}

func TestFieldRegistrationAndRetrieval(t *testing.T) {
	SetField("mock", NewMockField)

	// Test retrieval of the newly set field
	mockField, err := GetField("mock", nil)
	if err != nil {
		t.Errorf("GetField() failed: %v", err)
	}

	if mockField.Type() != "mock" {
		t.Errorf("GetField() returned incorrect type: got %v, want %v", mockField.Type(), "mock")
	}

	// Test deletion
	DeleteField("mock")
	_, err = GetField("mock", nil)
	if err == nil {
		t.Error("GetField() should have failed after deletion, but it didn't")
	}
}

func TestMockField_ProcessAndSearch(t *testing.T) {
	// Register the mock field for this test
	SetField("mock", NewMockField)
	mockField, _ := GetField("mock", nil)

	// Process a value
	testValue := "test"
	mockField.Process(testValue)

	// Prepare search bytes
	searchBytes, _ := mockField.ToSearchBytes(testValue)

	// Perform search
	found, _ := mockField.Search(searchBytes)
	if !found {
		t.Errorf("Search() should find the processed value, but it didn't")
	}
}
