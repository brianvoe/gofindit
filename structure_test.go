package gofindit

import (
	"testing"
)

func TestGetStructurePerson(t *testing.T) {
	_, doc := generateDoc()

	// Get structure of the person
	structure, err := getStructure(doc, "")
	if err != nil {
		t.Errorf("Failed to get structure: %v", err)
	}

	// Define expected fields to be found
	expectedFields := map[string]bool{
		// type TestData struct {
		// 	Name      string    `find:"name"`
		// 	Age       int       `find:"age"`
		// 	Hobbies   []string  `find:"hobbies"`
		// 	Pets      []TestPet `find:"pets"`
		// 	Bio       string    `find:"bio"`
		// 	isStudent bool      `find:"isStudent"`
		// 	Birthday  time.Time `find:"birthday"`
		// }

		// type TestPet struct {
		// 	Name  string   `find:"name"`
		// 	Age   int      `find:"age"`
		// 	Type  string   `find:"type"`
		// 	Breed string   `find:"breed"`
		// 	Toys  []string `find:"toys"`
		// }

		"name":         true,
		"age":          true,
		"hobbies":      true,
		"pets":         true,
		"bio":          true,
		"isStudent":    true,
		"birthday":     true,
		"pets[0].name": true,
		"pets[0].age":  true,
	}

	// Check if all expected fields are found
	for field := range expectedFields {
		if _, found := structure[field]; !found {
			t.Errorf("Field %v not found in structure", field)
		}
	}

	// Optionally, verify the total number of fields matches expectations
	if len(structure) != len(expectedFields) {
		t.Errorf("Expected %d fields, found %d", len(expectedFields), len(structure))
	}
}
