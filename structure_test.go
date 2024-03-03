package gofindit

import (
	"testing"
)

func TestGetStructure(t *testing.T) {
	type Test struct {
		Name string `find:"Name"`
		Age  int    `find:"Age"`
	}

	doc := Test{
		Name: "Test",
		Age:  10,
	}

	// Get structure of the document
	structure, err := getStructure(doc, "")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Check if the structure is correct
	if len(structure) != 2 {
		t.Errorf("Expected 2, got %v", len(structure))
	}

	// Check if the structure is correct
	if structure["Name"].Type != "string" {
		t.Errorf("Expected string, got %v", structure["Name"].Type)
	}
	if structure["Age"].Type != "int" {
		t.Errorf("Expected int, got %v", structure["Age"].Type)
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
