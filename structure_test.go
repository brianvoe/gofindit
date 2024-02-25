package gofindit

import (
	"reflect"
	"testing"
)

func TestGetFieldTypeMap(t *testing.T) {
	type Test struct {
		Name string `find:"name"`
		Age  int    `find:"age"`
	}

	structure, err := getFieldTypes(Test{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := map[string]string{
		"name": "string",
		"age":  "int",
	}

	if !reflect.DeepEqual(structure, expected) {
		t.Errorf("Expected %v, got %v", expected, structure)
	}
}

func TestGetBaseType(t *testing.T) {
	typeStr := getBaseType("string")
	if typeStr != "string" {
		t.Errorf("Expected string, got %v", typeStr)
	}

	// Test int types
	typeStr = getBaseType(11)
	if typeStr != "int" {
		t.Errorf("Expected int, got %v", typeStr)
	}
	typeStr = getBaseType(int8(11))
	if typeStr != "int" {
		t.Errorf("Expected int, got %v", typeStr)
	}
	typeStr = getBaseType(int16(11))
	if typeStr != "int" {
		t.Errorf("Expected int, got %v", typeStr)
	}
	typeStr = getBaseType(int32(11))
	if typeStr != "int" {
		t.Errorf("Expected int, got %v", typeStr)
	}
	typeStr = getBaseType(int64(11))
	if typeStr != "int" {
		t.Errorf("Expected int, got %v", typeStr)
	}

	// Test uint types
	typeStr = getBaseType(uint(11))
	if typeStr != "uint" {
		t.Errorf("Expected uint, got %v", typeStr)
	}
	typeStr = getBaseType(uint8(11))
	if typeStr != "uint" {
		t.Errorf("Expected uint, got %v", typeStr)
	}
	typeStr = getBaseType(uint16(11))
	if typeStr != "uint" {
		t.Errorf("Expected uint, got %v", typeStr)
	}
	typeStr = getBaseType(uint32(11))
	if typeStr != "uint" {
		t.Errorf("Expected uint, got %v", typeStr)
	}
	typeStr = getBaseType(uint64(11))
	if typeStr != "uint" {
		t.Errorf("Expected uint, got %v", typeStr)
	}

	// Test float types
	typeStr = getBaseType(11.0)
	if typeStr != "float" {
		t.Errorf("Expected float, got %v", typeStr)
	}
	typeStr = getBaseType(float32(11.0))
	if typeStr != "float" {
		t.Errorf("Expected float, got %v", typeStr)
	}
	typeStr = getBaseType(float64(11.0))
	if typeStr != "float" {
		t.Errorf("Expected float, got %v", typeStr)
	}

	// Test bool type
	typeStr = getBaseType(true)
	if typeStr != "bool" {
		t.Errorf("Expected bool, got %v", typeStr)
	}

	// Test slice types
	typeStr = getBaseType([]string{"test"})
	if typeStr != "[]string" {
		t.Errorf("Expected []string, got %v", typeStr)
	}
	typeStr = getBaseType([]int{1})
	if typeStr != "[]int" {
		t.Errorf("Expected []int, got %v", typeStr)
	}
	typeStr = getBaseType([]uint{1})
	if typeStr != "[]uint" {
		t.Errorf("Expected []uint, got %v", typeStr)
	}
	typeStr = getBaseType([]float64{1.0})
	if typeStr != "[]float" {
		t.Errorf("Expected []float, got %v", typeStr)
	}
	typeStr = getBaseType([]bool{true})
	if typeStr != "[]bool" {
		t.Errorf("Expected []bool, got %v", typeStr)
	}
}
